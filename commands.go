package vacbot

import "encoding/xml"

const (
	COMMAND_MOVE_FORWARD = `<ctl td="Move"><move action="forward" /></ctl>`
	COMMAND_SPIN_LEFT    = `<ctl td="Move"><move action="SpinLeft" /></ctl>`
	COMMAND_SPIN_RIGHT   = `<ctl td="Move"><move action="SpinRight" /></ctl>`
	COMMAND_TURN_AROUND  = `<ctl td="Move"><move action="TurnAround" /></ctl>`
	COMMAND_STOP_MOVING  = `<ctl td="Move"><move action="stop" /></ctl>`

	COMMAND_GET_BATTERY_INFO = `<ctl td="GetBatteryInfo" />`
)

type BatteryResponse struct {
	XMLName xml.Name    `xml:"query"`
	Xmlns   string      `xml:"xmlns,attr"`
	Ctl     *BatteryCtl `xml:"ctl"`
}

type BatteryCtl struct {
	Id      string   `xml:"id,attr"`
	Ret     string   `xml:"ret,attr"`
	Errno   string   `xml:"errno,attr"`
	Battery *Battery `xml:"battery"`
}

type Battery struct {
	Power string `xml:"power,attr"`
}
