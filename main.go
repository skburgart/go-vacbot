package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	log.Println("starting go-deebot service")
	init_api()
}

func init_api() {
	log.Println("initializing deebot api")
	uid, access_token := login(email, password_hash)
	authCode := get_auth_code(uid, access_token)
	userId, userAccessToken := get_user_access_token(uid, authCode)
	log.Printf("userId '%s' userAccessToken '%s'\n", userId, userAccessToken)
}

var (
	MAIN_URL = "https://eco-%s-api.ecovacs.com/v1/private/%s/%s/%s/%s/%s/%s/%s"
	USER_URL = "https://users-%s.ecouser.net:8000/user.do"
)

func login(email, passwordHash string) (string, string) {
	loginMap := map[string]string{
		"account":  encrypt(email),
		"password": encrypt(password_hash),
	}
	responseJson := call_main_api("user/login", loginMap)

	code := responseJson["code"].(string)
	if code != "0000" {
		log.Fatal("login error")
	}
	log.Println("login successful")

	data := responseJson["data"].(map[string]interface{})
	uid := data["uid"].(string)
	accessToken := data["accessToken"].(string)

	return uid, accessToken
}

func get_auth_code(uid, accessToken string) string {
	authMap := map[string]string{
		"uid":         uid,
		"accessToken": accessToken,
	}
	responseJson := call_main_api("user/getAuthCode", authMap)

	code := responseJson["code"].(string)
	if code != "0000" {
		log.Fatal("get auth code error")
	}
	log.Println("get auth code successful")

	data := responseJson["data"].(map[string]interface{})
	authCode := data["authCode"].(string)

	return authCode
}

func call_main_api(endpoint string, args map[string]string) map[string]interface{} {
	args["requestId"] = md5hash(time.Now().String())
	sign(args)

	client := &http.Client{}

	url := fmt.Sprintf("%s/%s", get_main_url(), endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	for _, k := range sortedKeys(args) {
		q.Add(k, args[k])
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(resp_body, &result)
	return result
}

func call_user_api(function string, args map[string]string) map[string]interface{} {
	args["todo"] = function

	jsonArgs, err := json.Marshal(args)
	if err != nil {
		log.Fatal("error marshalling user api json")
	}

	resp, err := http.Post(get_user_url(), "application/json", bytes.NewBuffer(jsonArgs))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(resp_body, &result)
	return result
}

func call_login_by_it_token(uid, auth_code string) map[string]interface{} {
	args := map[string]string{
		"country":  strings.ToUpper(country),
		"resource": resource,
		"realm":    realm,
		"userId":   uid,
		"token":    auth_code,
	}

	return call_user_api("loginByItToken", args)
}

func get_user_access_token(uid, authCode string) (string, string) {
	responseJson := call_login_by_it_token(uid, authCode)

	result := responseJson["result"].(string)
	if result != "ok" {
		log.Fatal("get user access token error")
	}
	log.Println("get user access token successful")

	userId := responseJson["userId"].(string)
	userAccessToken := responseJson["token"].(string)

	return userId, userAccessToken
}

func get_main_url() string {
	return fmt.Sprintf(MAIN_URL, country, country, lang, device_id, app_code, app_version, channel, device_type)
}

func get_user_url() string {
	return fmt.Sprintf(USER_URL, continent)
}
