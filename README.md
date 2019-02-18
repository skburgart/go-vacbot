# go-vacbot

A golang library for communicating with Ecovacs Deebot robot vacuums.

Inspired by [sucks](https://github.com/wpietri/sucks) from [William Pietri](https://github.com/wpietri).

## Quick Start
First create a config file similar to `vacbot.example.json`.

```golang
package main

import (
	"flag"
	"time"

	vacbot "github.com/skburgart/go-vacbot"
)

func main() {
	vacbotConfigFile := flag.String("vacbotconfig", "vacbot.json", "json file containing vacbot configuration")
	flag.Parse()

	v := vacbot.NewFromConfigFile(*vacbotConfigFile)
	v.SpinLeft()
	time.Sleep(2 * time.Second)
	v.StopMoving()

	v.SpinRight()
	time.Sleep(2 * time.Second)
	v.StopMoving()

	v.Forward()
	time.Sleep(2 * time.Second)
	v.StopMoving()

	v.CleanAuto()
	time.Sleep(5 * time.Second)
	v.CleanStop()
}
```
