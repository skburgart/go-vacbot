# go-vacbot

A golang library for communicating with Ecovacs Deebot robot vacuums.

Based off of [sucks](https://github.com/wpietri/sucks) by [William Pietri](https://github.com/wpietri).

## Quick Start
First create a config file similar to `vacbot.example.json`.

```
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
	v.SpinRight()
	v.Forward()
	time.Sleep(2 * time.Second)
	v.StopMoving()
}
```
