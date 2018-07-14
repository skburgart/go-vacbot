package main

type Config struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	DeviceId     string `json:"device_id"`
	Resource     string `json:"resource"`
	Country      string `json:"country"`
	Continent    string `json:"continent"`
	Lang         string `json:"lang"`
	AppCode      string `json:"app_code"`
	AppVersion   string `json:"app_version"`
	Channel      string `json:"channel"`
	DeviceType   string `json:"device_type"`
	Timezone     string `json:"timezone"`
	Realm        string `json:"realm"`
}
