package main

import (
	"flag"
	"log"

	"gitlab.delat.net/skburgart/go-vacbot/vacbot"
)

func main() {
	configFile := flag.String("configfile", "vacbot.json", "json file containing vacbot configuration")
	flag.Parse()

	log.Println("starting go-deebot service")
	v := vacbot.New(*configFile)
	v.TurnLeft(90)
	v.TurnRight(45)
	v.TurnRight(45)
	v.TurnRight(180)
	v.TurnLeft(90)
	v.TurnLeft(90)
}
