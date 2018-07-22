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
)

const (
	MAIN_URL = "https://eco-%s-api.ecovacs.com/v1/private/%s/%s/%s/%s/%s/%s/%s"
	USER_URL = "https://users-%s.ecouser.net:8000/user.do"
)

var config Config

type Client struct {
	vx *VacbotXMPP
}

func New(c Config) *Client {
	config = c
	uid, accesstoken := login(config.Email, config.PasswordHash)
	authCode := getAuthCode(uid, accesstoken)
	userId, userAccessToken := getUserAccessToken(uid, authCode)
	deviceJID := getFirstDeviceAddress(userId, userAccessToken)
	vx := NewVacbotXMPP(userId, userAccessToken, deviceJID)

	return &Client{
		vx: vx,
	}
}

func (c *Client) RecvHandler(handlerFunc func(interface{}, error)) {
	go func() {
		for {
			handlerFunc(c.vx.client.Recv())
		}
	}()
}

func NewFromConfigFile(configFile string) *Client {
	config = LoadConfiguration(configFile)
	return New(config)
}

func (c *Client) FetchBatteryLevel() {
	c.vx.issueCommand(COMMAND_GET_BATTERY_INFO)
}

func (c *Client) Forward() {
	c.vx.issueCommand(COMMAND_MOVE_FORWARD)
}

func (c *Client) SpinLeft() {
	c.vx.issueCommand(COMMAND_SPIN_LEFT)
}

func (c *Client) SpinRight() {
	c.vx.issueCommand(COMMAND_SPIN_RIGHT)
}

func (c *Client) TurnAround() {
	c.vx.issueCommand(COMMAND_TURN_AROUND)
}

func (c *Client) StopMoving() {
	c.vx.issueCommand(COMMAND_STOP_MOVING)
}

func (c *Client) FetchCleanState() {
	c.vx.issueCommand(COMMAND_GET_CLEAN_STATE)
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

func login(email, passwordHash string) (string, string) {
	loginMap := map[string]string{
		"account":  encrypt(config.Email),
		"password": encrypt(config.PasswordHash),
	}
	responseJson := callMainApi("user/login", loginMap)

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

func getAuthCode(uid, accessToken string) string {
	authMap := map[string]string{
		"uid":         uid,
		"accessToken": accessToken,
	}
	responseJson := callMainApi("user/getAuthCode", authMap)

	code := responseJson["code"].(string)
	if code != "0000" {
		log.Fatal("get auth code error")
	}
	log.Println("get auth code successful")

	data := responseJson["data"].(map[string]interface{})
	authCode := data["authCode"].(string)

	return authCode
}

func callMainApi(endpoint string, args map[string]string) map[string]interface{} {
	args["requestId"] = md5hash(time.Now().String())
	sign(args)

	client := &http.Client{}

	url := fmt.Sprintf("%s/%s", getMainUrl(), endpoint)
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
	respBody, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	return result
}

func callUserApi(function string, args map[string]interface{}) map[string]interface{} {
	args["todo"] = function

	jsonArgs, err := json.Marshal(args)
	if err != nil {
		log.Fatal("error marshalling user api json")
	}

	resp, err := http.Post(getUserUrl(), "application/json", bytes.NewBuffer(jsonArgs))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	return result
}

func getDevices(userId, userAccessToken string) map[string]interface{} {
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
	return callUserApi("GetDeviceList", args)
}

func getFirstDeviceAddress(userId, userAccessToken string) string {
	deviceJson := getDevices(userId, userAccessToken)
	deviceList := deviceJson["devices"].([]interface{})
	firstDevice := deviceList[0]
	return getDeviceAddress(firstDevice.(map[string]interface{}))
}

func getDeviceAddress(deviceJson map[string]interface{}) string {
	deviceId := deviceJson["did"].(string)
	deviceClass := deviceJson["class"].(string)

	return fmt.Sprintf("%s@%s.ecorobot.net/atom", deviceId, deviceClass)
}

func callLoginByItToken(uid, authCode string) map[string]interface{} {
	args := map[string]interface{}{
		"country":  strings.ToUpper(config.Country),
		"resource": config.Resource,
		"realm":    config.Realm,
		"userId":   uid,
		"token":    authCode,
	}

	return callUserApi("loginByItToken", args)
}

func getUserAccessToken(uid, authCode string) (string, string) {
	responseJson := callLoginByItToken(uid, authCode)

	result := responseJson["result"].(string)
	if result != "ok" {
		log.Fatal("get user access token error")
	}
	log.Println("get user access token successful")

	userId := responseJson["userId"].(string)
	userAccessToken := responseJson["token"].(string)

	return userId, userAccessToken
}

func getMainUrl() string {
	return fmt.Sprintf(MAIN_URL, config.Country, config.Country, config.Lang, config.DeviceId, config.AppCode, config.AppVersion, config.Channel, config.DeviceType)
}

func getUserUrl() string {
	return fmt.Sprintf(USER_URL, config.Continent)
}
