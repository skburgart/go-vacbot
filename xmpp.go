package vacbot

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"strconv"

	xmpp "github.com/mattn/go-xmpp"
)

const (
	XMPP_URL = "msg-%s.ecouser.net:5223"
)

type VacbotXMPP struct {
	client       *xmpp.Client
	from         string
	to           string
	batteryLevel int
}

func NewVacbotXMPP(userId, userAccessToken, deviceJID string) *VacbotXMPP {
	xmppPassword := fmt.Sprintf("0/%s/%s", config.Resource, userAccessToken)
	xmppOpts := xmpp.Options{
		Host:     getXmppUrl(),
		User:     fmt.Sprintf("%s@%s", userId, config.Realm),
		Password: xmppPassword,
		NoTLS:    true,
		Debug:    false,
		Session:  true,
	}
	xmppClient, err := xmppOpts.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	vx := &VacbotXMPP{
		client: xmppClient,
		from:   xmppClient.JID(),
		to:     deviceJID,
	}

	go vx.Recv()

	return vx
}

func (vx *VacbotXMPP) Recv() {
	log.Println("starting recv")
	for {
		stanza, err := vx.client.Recv()
		if err != nil {
			log.Fatal(err)
		}
		batteryLevel, err := parseBatteryLevel(stanza)
		if err != nil {
			continue
		}

		vx.batteryLevel = batteryLevel
	}
}

func parseBatteryLevel(stanza interface{}) (int, error) {
	switch t := stanza.(type) {
	case xmpp.IQ:
		if t.Query == nil {
			return 0, errors.New("nil query")
		}
		parsedResponse := &BatteryResponse{}
		err := xml.Unmarshal(t.Query, parsedResponse)
		if err != nil {
			log.Printf("failed xml unmarshal: %v\n", err)
		}
		batteryLevel, err := strconv.Atoi(parsedResponse.Ctl.Battery.Power)
		if err != nil {
			log.Printf("failed converting battery level to int unmarshal: %v\n", err)
		}

		return batteryLevel, nil
	default:
		return 0, errors.New("message not IQ")
	}
}

func (vx *VacbotXMPP) issueCommand(command string) {
	_, err := vx.client.RawInformationQuery(vx.from, vx.to, "1", xmpp.IQTypeSet, "com:ctl", command)
	if err != nil {
		log.Fatal(err)
	}
}

func getXmppUrl() string {
	return fmt.Sprintf(XMPP_URL, config.Continent)
}
