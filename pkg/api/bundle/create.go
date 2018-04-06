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
	"github.com/george-e-shaw-iv/nixtools"
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

	/* Check if account port is in use by system */
	check, err := net.Listen("tcp", ":"+strconv.Itoa(createBundleRequestData.AccPort))
	if err != nil {
		logger.Println(req.URL.Path + "::" + "a service is already listening on port " + strconv.Itoa(createBundleRequestData.AccPort) + "::" + err.Error())
		http.Error(res, "A service is already listening on port "+strconv.Itoa(createBundleRequestData.AccPort), http.StatusInternalServerError)
		return false
	}
	check.Close()

	/* Check if public port is in use by system */
	check, err = net.Listen("tcp", ":"+strconv.Itoa(createBundleRequestData.PubPort))
	if err != nil {
		logger.Println(req.URL.Path + "::" + "a service is already listening on port " + strconv.Itoa(createBundleRequestData.PubPort) + "::" + err.Error())
		http.Error(res, "A service is already listening on port "+strconv.Itoa(createBundleRequestData.PubPort), http.StatusInternalServerError)
		return false
	}
	check.Close()

	/* Check if public/account ports are in use by another bundle */
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

	/* Check if admin settings are set and valid */
	ds, err := database.Open("server/" + database.DB_SETTINGS)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	err = ds.CheckAdminSettings()
	if err != nil {
		ds.Close()
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}
	ds.Close()

	/* Check if bundle username can be used within the system */
	if userExists := nixtools.UserExists(createBundleRequestData.Name); userExists {
		logger.Println(req.URL.Path + "::username already exists under that bundle name within the system")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	/* IF ALL CHECKS PASSED, start the bundle creation process */

	/* Create the bundle directory within the gPanel directory */
	newBundle := "bundles/" + createBundleRequestData.Name
	err = os.Mkdir(newBundle, 0777)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Create the log folder within the new bundle directory */
	err = os.Mkdir(newBundle+"/logs", 0777)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Create and open the newly created bundle's main datastore */
	ds, err = database.Open(newBundle + "/" + database.DB_MAIN)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Put the new bundle's ports within the datastore */
	var databaseBundlePorts struct {
		Account int `json:"account"`
		Public  int `json:"public"`
	}
	databaseBundlePorts.Account = createBundleRequestData.AccPort
	databaseBundlePorts.Public = createBundleRequestData.PubPort

	err = ds.Put(database.BUCKET_PORTS, []byte("bundle_ports"), databaseBundlePorts)
	if err != nil {
		ds.Close()
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Create the default bundle username and password */
	var defaultBundleUser database.Struct_Users
	var tempPass = encryption.RandomString(16)

	defaultBundleUser.Pass, err = encryption.HashPassword(tempPass)
	if err != nil {
		ds.Close()
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	defaultBundleUser.Secret = ""

	err = ds.Put(database.BUCKET_USERS, []byte("root"), defaultBundleUser)
	if err != nil {
		ds.Close()
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	ds.Close()

	/* Create a system user for the new bundle */
	sysUser, err := nixtools.GetUser(createBundleRequestData.Name, true)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Initialize SSH access for the new user and give the root account access */
	err = sysUser.InitSSH(true)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Create document_root within new system user's home directory */
	err = sysUser.CreateDirectory("document_root", 0777)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Write the default index within the document_root */
	err = sysUser.WriteFile("document_root/index.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777, []byte(DEFAULT_INDEX))
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Add newly created bundle to list of current bundles (run-time list) */
	bundles[createBundleRequestData.Name], err = gpaccount.New(newBundle+"/", createBundleRequestData.Name, databaseBundlePorts.Account, databaseBundlePorts.Public)

	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Start the account server and public server of the newly created bundle */
	err = bundles[createBundleRequestData.Name].Start()
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	err = bundles[createBundleRequestData.Name].Public.Start()
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	/* Get SMTP and Admin general settings for use in the email to be sent to new bundle designated email */
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

	/* Send email with information about system and bundle account to designated bundle email AND gPanel admin */
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
