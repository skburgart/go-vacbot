package vacbot

const (
	COMMAND_MOVE_FORWARD = `<ctl td="Move"><move action="forward" /></ctl>`
	COMMAND_SPIN_LEFT    = `<ctl td="Move"><move action="SpinLeft" /></ctl>`
	COMMAND_SPIN_RIGHT   = `<ctl td="Move"><move action="SpinRight" /></ctl>`
	COMMAND_TURN_AROUND  = `<ctl td="Move"><move action="TurnAround" /></ctl>`
	COMMAND_STOP_MOVING  = `<ctl td="Move"><move action="stop" /></ctl>`

	COMMAND_CLEAN_AUTO        = `<ctl td="Clean"><clean type="auto" speed="standard" /></ctl>`
	COMMAND_CLEAN_AUTO_STRONG = `<ctl td="Clean"><clean type="auto" speed="strong" /></ctl>`

	COMMAND_CLEAN_BORDER        = `<ctl td="Clean"><clean type="border" speed="standard" /></ctl>`
	COMMAND_CLEAN_BORDER_STRONG = `<ctl td="Clean"><clean type="border" speed="strong" /></ctl>`

	COMMAND_CLEAN_SPOT        = `<ctl td="Clean"><clean type="spot" speed="standard" /></ctl>`
	COMMAND_CLEAN_SPOT_STRONG = `<ctl td="Clean"><clean type="spot" speed="strong" /></ctl>`

	COMMAND_CLEAN_SINGLEROOM        = `<ctl td="Clean"><clean type="singleroom" speed="standard" /></ctl>`
	COMMAND_CLEAN_SINGLEROOM_STRONG = `<ctl td="Clean"><clean type="singleroom" speed="strong" /></ctl>`

	COMMAND_CLEAN_STOP = `<ctl td="Clean"><clean type="stop" speed="standard" /></ctl>`

	COMMAND_CHARGE = `<ctl td="Charge"><charge type="go"/></ctl>`

	COMMAND_GET_BATTERY_INFO = `<ctl td="GetBatteryInfo" />`
	COMMAND_GET_CLEAN_STATE  = `<ctl td="GetCleanState" />`
)
