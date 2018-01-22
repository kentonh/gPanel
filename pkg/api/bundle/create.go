package bundle

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/emailer"
	"github.com/Ennovar/gPanel/pkg/encryption"
	"github.com/Ennovar/gPanel/pkg/gpaccount"
	"github.com/Ennovar/gPanel/pkg/system"
)

func Create(res http.ResponseWriter, req *http.Request, logger *log.Logger, bundles map[string]*gpaccount.Controller) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var createBundleRequestData struct {
		Name    string `json:"name"`
		AccPort int    `json:"account_port"`
		PubPort int    `json:"public_port"`
		Email   string `json:"email"`
	}

	err := json.NewDecoder(req.Body).Decode(&createBundleRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	check, err := net.Listen("tcp", ":"+strconv.Itoa(createBundleRequestData.AccPort))
	if err != nil {
		logger.Println(req.URL.Path + "::" + "a service is already listening on port " + strconv.Itoa(createBundleRequestData.AccPort) + "::" + err.Error())
		http.Error(res, "A service is already listening on port "+strconv.Itoa(createBundleRequestData.AccPort), http.StatusInternalServerError)
		return false
	}
	check.Close()

	check, err = net.Listen("tcp", ":"+strconv.Itoa(createBundleRequestData.PubPort))
	if err != nil {
		logger.Println(req.URL.Path + "::" + "a service is already listening on port " + strconv.Itoa(createBundleRequestData.PubPort) + "::" + err.Error())
		http.Error(res, "A service is already listening on port "+strconv.Itoa(createBundleRequestData.PubPort), http.StatusInternalServerError)
		return false
	}
	check.Close()

	err = nil
	for k, v := range bundles {
		if k == createBundleRequestData.Name {
			err = errors.New("Bundle \"" + k + "\" already exists")
			break
		}

		if v.Port == createBundleRequestData.AccPort ||
			v.Port == createBundleRequestData.PubPort ||
			v.Public.Port == createBundleRequestData.AccPort ||
			v.Public.Port == createBundleRequestData.PubPort {
			err = errors.New("An existing bundle is using the port \"" + strconv.Itoa(v.Port) + "\" already")
			break
		}
	}

	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	newBundle := "bundles/" + createBundleRequestData.Name
	err = os.Mkdir(newBundle, 0777)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	err = os.Mkdir(newBundle+"/logs", 0777)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	ds, err := database.Open(newBundle + "/" + database.DB_MAIN)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	var databaseBundlePorts struct {
		Account int `json:"account"`
		Public  int `json:"public"`
	}
	databaseBundlePorts.Account = createBundleRequestData.AccPort
	databaseBundlePorts.Public = createBundleRequestData.PubPort

	err = ds.Put(database.BUCKET_PORTS, []byte("bundle_ports"), databaseBundlePorts)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	var defaultBundleUser database.Struct_Users
	var tempPass = encryption.RandomString(16)

	defaultBundleUser.Pass, err = encryption.HashPassword(tempPass)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	defaultBundleUser.Secret = ""

	err = ds.Put(database.BUCKET_USERS, []byte("root"), defaultBundleUser)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	ds.Close()

	err, err2 := system.CreateBundleUser(createBundleRequestData.Name)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error() + " AND " + err2.Error())
		http.Error(res, err2.Error(), http.StatusInternalServerError)
		return false
	}

	bundles[createBundleRequestData.Name] = gpaccount.New(newBundle+"/", createBundleRequestData.Name, databaseBundlePorts.Account, databaseBundlePorts.Public)
	_ = bundles[createBundleRequestData.Name].Start()
	_ = bundles[createBundleRequestData.Name].Public.Start()

	ds, err = database.Open("server/" + database.DB_SETTINGS)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	var smtpSettings database.Struct_SMTP
	var adminSettings database.Struct_Admin

	err = ds.Get(database.BUCKET_GENERAL, []byte("smtp"), &smtpSettings)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	err = ds.Get(database.BUCKET_GENERAL, []byte("admin"), &adminSettings)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	mail, err := emailer.New(smtpSettings.Type, emailer.Credentials{
		Username: smtpSettings.Username,
		Password: smtpSettings.Password,
		Server:   smtpSettings.Server,
		Port:     smtpSettings.Port,
	})
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	nameservers, err := ds.ListNameservers()
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	var msg string
	if len(nameservers) > 0 {
		var nameserversString string
		for _, v := range nameservers {
			nameserversString += "- " + v + "\r\n"
		}

		msg = string("Your new gPanel Bundle has been successfully registered.\r\n\n" +
			"Account Port: " + strconv.Itoa(createBundleRequestData.AccPort) + "\r\n" +
			"Public Port: " + strconv.Itoa(createBundleRequestData.PubPort) + "\r\n\n" +
			"Default account username: root\r\n" +
			"Default account password: " + tempPass + "\r\n\n" +
			"Any questions, comments, or concerns can be directed toward your server administrator " + adminSettings.Name +
			" at " + adminSettings.Email + "\r\n\n" + "This server impliments the following nameservers:\r\n" + nameserversString)
	} else {
		msg = string("Your new gPanel Bundle has been successfully registered.\r\n\n" +
			"Account Port: " + strconv.Itoa(createBundleRequestData.AccPort) + "\r\n" +
			"Public Port: " + strconv.Itoa(createBundleRequestData.PubPort) + "\r\n\n" +
			"Default account username: root\r\n" +
			"Default account password: " + tempPass + "\r\n\n" +
			"Any questions, comments, or concerns can be directed toward your server administrator " + adminSettings.Name +
			" at " + adminSettings.Email)
	}

	err = mail.SendSimple(createBundleRequestData.Email, "New gPanel Bundle - "+createBundleRequestData.Name, msg)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	err = mail.SendSimple(adminSettings.Email, "New gPanel Bundle - "+createBundleRequestData.Name, msg)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(createBundleRequestData.Name))

	return true
}
