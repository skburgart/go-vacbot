package vacbot

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"sort"
	"time"
)

var (
	CLIENT_KEY         = "eJUWrzRv34qFSaYk"
	SECRET             = "Cyu5jcR4zyK6QEPn1hdIGXB5QIDAQABMA0GC"
	ECOVACS_PUBLIC_KEY = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDb8V0OYUGP3Fs63E1gJzJh+7iq
eymjFUKJUqSD60nhWReZ+Fg3tZvKKqgNcgl7EGXp1yNifJKUNC/SedFG1IJRh5hB
eDMGq0m0RQYDpf9l0umqYURpJ5fmfvH/gjfHe3Eg/NTLm7QEa0a0Il2t3Cyu5jcR
4zyK6QEPn1hdIGXB5QIDAQAB
-----END PUBLIC KEY-----`
)

func encrypt(message string) string {
	block, _ := pem.Decode([]byte(ECOVACS_PUBLIC_KEY))

	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(message))
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(cipherText)
}

func meta() map[string]string {
	return map[string]string{
		"country":    config.Country,
		"lang":       config.Lang,
		"deviceId":   config.DeviceId,
		"appCode":    config.AppCode,
		"appVersion": config.AppVersion,
		"channel":    config.Channel,
		"deviceType": config.DeviceType,
	}
}

func sign(args map[string]string) map[string]string {
	args["authTimespan"] = fmt.Sprintf("%v", time.Now().UnixNano()/1000000)
	args["authTimeZone"] = config.Timezone

	textToSign := CLIENT_KEY

	// prepare signature map
	signMap := meta()
	for k, v := range args {
		signMap[k] = v
	}

	// construct message to sign
	for _, k := range sortedKeys(signMap) {
		textToSign += fmt.Sprintf("%s=%s", k, signMap[k])
	}
	textToSign += SECRET

	args["authAppkey"] = CLIENT_KEY
	args["authSign"] = md5hash(textToSign)

	return args
}

func sortedKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func md5hash(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}
