package vacbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	xmpp "github.com/mattn/go-xmpp"
)

var config Config

type VacBot struct {
}

func New(configFile string) *VacBot {
	v := &VacBot{}

	config = LoadConfiguration(configFile)
	uid, access_token := login(config.Email, config.PasswordHash)
	authCode := get_auth_code(uid, access_token)
	userId, userAccessToken := get_user_access_token(uid, authCode)
	xmppPassword := fmt.Sprintf("0/%s/%s", config.Resource, userAccessToken)
	_, err := xmpp.NewClientNoTLS(get_xmpp_url(), fmt.Sprintf("%s@%s", userId, config.Realm), xmppPassword, false)
	if err != nil {
		log.Fatal(err)
	}

	return v
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	config.Resource = config.DeviceId[:8]
	return config
}

var (
	MAIN_URL = "https://eco-%s-api.ecovacs.com/v1/private/%s/%s/%s/%s/%s/%s/%s"
	USER_URL = "https://users-%s.ecouser.net:8000/user.do"
	XMPP_URL = "msg-%s.ecouser.net:5223"
)

func login(email, passwordHash string) (string, string) {
	loginMap := map[string]string{
		"account":  encrypt(config.Email),
		"password": encrypt(config.PasswordHash),
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

func call_user_api(function string, args map[string]interface{}) map[string]interface{} {
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

func get_devices(userId, userAccessToken string) map[string]interface{} {
	args := map[string]interface{}{
		"userid": userId,
		"auth": map[string]string{
			"with":     "users",
			"userid":   userId,
			"realm":    config.Realm,
			"token":    userAccessToken,
			"resource": config.Resource,
		},
	}
	return call_user_api("GetDeviceList", args)
}

func get_first_device_address(userId, userAccessToken string) string {
	deviceJson := get_devices(userId, userAccessToken)
	deviceList := deviceJson["devices"].([]interface{})
	firstDevice := deviceList[0]
	return get_device_address(firstDevice.(map[string]interface{}))
}

func get_device_address(deviceJson map[string]interface{}) string {
	deviceId := deviceJson["did"].(string)
	deviceClass := deviceJson["class"].(string)

	return fmt.Sprintf("%s@%s.ecorobot.net/atom", deviceId, deviceClass)
}

func call_login_by_it_token(uid, auth_code string) map[string]interface{} {
	args := map[string]interface{}{
		"country":  strings.ToUpper(config.Country),
		"resource": config.Resource,
		"realm":    config.Realm,
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
	return fmt.Sprintf(MAIN_URL, config.Country, config.Country, config.Lang, config.DeviceId, config.AppCode, config.AppVersion, config.Channel, config.DeviceType)
}

func get_user_url() string {
	return fmt.Sprintf(USER_URL, config.Continent)
}
func get_xmpp_url() string {
	return fmt.Sprintf(XMPP_URL, config.Continent)
}
