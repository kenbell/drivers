package epd2in9v2

// Registers
const (
	DRIVER_OUTPUT_CONTROL                = 0x01
	GATE_DRIVING_VOLTAGE_CONTROL         = 0x03
	SOURCE_DRIVING_VOLTAGE_CONTROL       = 0x04
	DEEP_SLEEP_MODE                      = 0x10
	DATA_ENTRY_MODE_SETTING              = 0x11
	SW_RESET                             = 0x12
	MASTER_ACTIVATION                    = 0x20
	DISPLAY_UPDATE_CONTROL_1             = 0x21
	DISPLAY_UPDATE_CONTROL_2             = 0x22
	WRITE_RAM                            = 0x24
	VCOM_SENSE                           = 0x28
	VCOM_SENSE_DURATION                  = 0x29
	PROGRAM_VCOM_OTP                     = 0x2A
	WRITE_REGISTER_FOR_VCOM_CONTROL      = 0x2B
	WRITE_VCOM_REGISTER                  = 0x2C
	OTP_REGISTER_READ_FOR_DISPLAY_OPTION = 0x2D
	USER_ID_READ                         = 0x2E
	PROGRAM_WS_OTP                       = 0x30
	LOAD_WS_OTP                          = 0x31
	WRITE_LUT_REGISTER                   = 0x32
	PROGRAM_OTP_SELECTION                = 0x36
	WRITE_REGISTER_FOR_USER_ID           = 0x38
	OTP_PROGRAM_MODE                     = 0x39
	SET_RAM_X_ADDRESS_START_END_POSITION = 0x44
	SET_RAM_Y_ADDRESS_START_END_POSITION = 0x45
	SET_RAM_X_ADDRESS_COUNTER            = 0x4E
	SET_RAM_Y_ADDRESS_COUNTER            = 0x4F

	NO_ROTATION  Rotation = 0
	ROTATION_90  Rotation = 1 // 90 degrees clock-wise rotation
	ROTATION_180 Rotation = 2
	ROTATION_270 Rotation = 3
)
