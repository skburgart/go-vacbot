package main

import (
	"log"

	"gitlab.delat.net/skburgart/go-deebot/vacbot"
)

func main() {
	log.Println("starting go-deebot service")
	vacbot.New()
}
