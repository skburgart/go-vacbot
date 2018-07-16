package vacbot

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
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
		Host:     get_xmpp_url(),
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
	//go Recv(xmppClient)

	return &VacbotXMPP{
		client: xmppClient,
		from:   xmppClient.JID(),
		to:     deviceJID,
	}
}

func Recv(xmppClient *xmpp.Client) {
	log.Println("starting recv")
	for {
		stanza, err := xmppClient.Recv()
		if err != nil {
			log.Fatal(err)
		}
		spew.Dump(stanza)
	}
}

func (vx *VacbotXMPP) issueCommand(command string) {
	_, err := vx.client.RawInformationQuery(vx.from, vx.to, "1", xmpp.IQTypeSet, "com:ctl", command)
	if err != nil {
		log.Fatal(err)
	}
}

func get_xmpp_url() string {
	return fmt.Sprintf(XMPP_URL, config.Continent)
}
