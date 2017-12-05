package database

type Struct_SMTP struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
}
