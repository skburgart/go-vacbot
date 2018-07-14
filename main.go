package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	client := &http.Client{}

	url := fmt.Sprintf("%s/%s", get_main_url(), endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	for _, k := range sortedKeys(args) {
		q.Add(k, args[k])
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	log.Printf(string(resp_body))
}

func get_main_url() string {
	return fmt.Sprintf(MAIN_URL, country, country, lang, device_id, app_code, app_version, channel, device_type)
}
