package vacbot

import (
	"fmt"
	"log"

	xmpp "github.com/mattn/go-xmpp"
)

const (
	XMPP_URL = "msg-%s.ecouser.net:5223"
)

type VacbotXMPP struct {
	client *xmpp.Client
	from   string
	to     string
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

	return vx
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
