package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println("starting go-deebot service")
	init_api()
}

func init_api() {
	log.Println("initializing deebot api")
	loginMap := map[string]string{
		"account":  encrypt(email),
		"password": encrypt(password_hash),
	}
	call_main_api("user/login", loginMap)
}

var (
	MAIN_URL = "https://eco-%s-api.ecovacs.com/v1/private/%s/%s/%s/%s/%s/%s/%s"
)

func call_main_api(endpoint string, args map[string]string) {
	args["requestId"] = md5hash(time.Now().String())
	sign(args)
}

func get_main_url() string {
	return fmt.Sprintf(MAIN_URL, country, country, lang, device_id, app_code, app_version, channel, device_type)
}
