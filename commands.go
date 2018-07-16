package vacbot

var (
	COMMAND_MOVE_FORWARD = `<ctl td="Move"><move action="forward" /></ctl>`
	COMMAND_SPIN_LEFT    = `<ctl td="Move"><move action="SpinLeft" /></ctl>`
	COMMAND_SPIN_RIGHT   = `<ctl td="Move"><move action="SpinRight" /></ctl>`
	COMMAND_TURN_AROUND  = `<ctl td="Move"><move action="TurnAround" /></ctl>`
	COMMAND_STOP_MOVING  = `<ctl td="Move"><move action="stop" /></ctl>`

	COMMAND_GET_BATTERY_INFO = `<ctl td="GetBatteryInfo" />`
)
