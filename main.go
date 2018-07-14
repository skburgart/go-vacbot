package main

import (
	"flag"
	"log"

	"gitlab.delat.net/skburgart/go-deebot/vacbot"
)

func main() {
	configFile := flag.String("configfile", "vacbot.json", "json file containing vacbot configuration")
	flag.Parse()

	log.Println("starting go-deebot service")
	vacbot.New(*configFile)
}
