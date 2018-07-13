package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
)

var (
	CLIENT_KEY = "eJUWrzRv34qFSaYk"
	SECRET     = "Cyu5jcR4zyK6QEPn1hdIGXB5QIDAQABMA0GC"
)

var ECOVACS_PUBLIC_KEY = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDb8V0OYUGP3Fs63E1gJzJh+7iq
eymjFUKJUqSD60nhWReZ+Fg3tZvKKqgNcgl7EGXp1yNifJKUNC/SedFG1IJRh5hB
eDMGq0m0RQYDpf9l0umqYURpJ5fmfvH/gjfHe3Eg/NTLm7QEa0a0Il2t3Cyu5jcR
4zyK6QEPn1hdIGXB5QIDAQAB
-----END PUBLIC KEY-----`

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
