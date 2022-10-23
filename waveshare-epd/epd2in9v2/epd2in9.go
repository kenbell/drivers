// Package epd2in9v2 implements a driver for Waveshare 2.9in V2 black and white e-paper device.
//
// Note: this is for the V2 device, the V1 device uses a different chipset.
//
// Datasheets:
//
//	https://www.waveshare.com/w/upload/7/79/2.9inch-e-paper-v2-specification.pdf
package epd2in9v2

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
)

type Config struct {
	Width        int16 // Width is the display resolution
	Height       int16
	LogicalWidth int16    // LogicalWidth must be a multiple of 8 and same size or bigger than Width
	Rotation     Rotation // Rotation is clock-wise
}

type Device struct {
	bus          drivers.SPI
	cs           machine.Pin
	dc           machine.Pin
	rst          machine.Pin
	busy         machine.Pin
	logicalWidth int16
	width        int16
	height       int16
	buffer       []uint8
	bufferLength uint32
	rotation     Rotation
}

type Rotation uint8

// Look up table for full updates
var lutFullUpdate = [159]uint8{
	0x80, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x0, 0x0, 0x0,
	0x10, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0,
	0x80, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x0, 0x0, 0x0,
	0x10, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x14, 0x8, 0x0, 0x0, 0x0, 0x0, 0x1,
	0xA, 0xA, 0x0, 0xA, 0xA, 0x0, 0x1,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x14, 0x8, 0x0, 0x1, 0x0, 0x0, 0x1,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x0, 0x0, 0x0,
	0x22, 0x17, 0x41, 0x0, 0x32, 0x36,
}

// Look up table for partial updates, faster but there will be some ghosting
var lutPartialUpdate = [159]uint8{
	0x0, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x80, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x40, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0A, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2,
	0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x0, 0x0, 0x0,
	0x22, 0x17, 0x41, 0xB0, 0x32, 0x36,
}

// New returns a new epd2in9 driver. Pass in a fully configured SPI bus.
func New(bus drivers.SPI, csPin, dcPin, rstPin, busyPin machine.Pin) Device {
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	busyPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	return Device{
		bus:  bus,
		cs:   csPin,
		dc:   dcPin,
		rst:  rstPin,
		busy: busyPin,
	}
}

// Configure sets up the device.
func (d *Device) Configure(cfg Config) {
	if cfg.LogicalWidth != 0 {
		d.logicalWidth = cfg.LogicalWidth
	} else {
		d.logicalWidth = 128
	}
	if cfg.Width != 0 {
		d.width = cfg.Width
	} else {
		d.width = 128
	}
	if cfg.Height != 0 {
		d.height = cfg.Height
	} else {
		d.height = 296
	}
	d.rotation = cfg.Rotation
	d.bufferLength = (uint32(d.logicalWidth) * uint32(d.height)) / 8
	d.buffer = make([]uint8, d.bufferLength)
	for i := uint32(0); i < d.bufferLength; i++ {
		d.buffer[i] = 0xFF
	}

	d.cs.Low()
	d.dc.Low()
	d.rst.Low()

	d.Reset()
	time.Sleep(100 * time.Millisecond)

	d.SendCommand(SW_RESET)
	d.WaitUntilIdle()

	d.SendCommand(DRIVER_OUTPUT_CONTROL)
	d.SendData(uint8((d.height - 1) & 0xFF))
	d.SendData(uint8(((d.height - 1) >> 8) & 0xFF))
	d.SendData(0x00) // GD = 0; SM = 0; TB = 0;

	d.SendCommand(DATA_ENTRY_MODE_SETTING)
	d.SendData(0x03) // X increment; Y increment

	d.SendCommand(DISPLAY_UPDATE_CONTROL_1)
	d.SendData(0x00)
	d.SendData(0x80)

	d.WaitUntilIdle()
	d.SetLUT(true)
}

// Reset resets the device
func (d *Device) Reset() {
	d.rst.High()
	time.Sleep(10 * time.Millisecond)
	d.rst.Low()
	time.Sleep(2 * time.Millisecond)
	d.rst.High()
	time.Sleep(10 * time.Millisecond)
}

// DeepSleep puts the display into deepsleep
func (d *Device) DeepSleep() {
	d.SendCommand(DEEP_SLEEP_MODE)
	d.WaitUntilIdle()
}

// SendCommand sends a command to the display
func (d *Device) SendCommand(command uint8) {
	d.sendDataCommand(true, command)
}

// SendData sends a data byte to the display
func (d *Device) SendData(data uint8) {
	d.sendDataCommand(false, data)
}

// sendDataCommand sends image data or a command to the screen
func (d *Device) sendDataCommand(isCommand bool, data uint8) {
	if isCommand {
		d.dc.Low()
	} else {
		d.dc.High()
	}
	d.cs.Low()
	d.bus.Transfer(data)
	d.cs.High()
}

// SetLUT sets the look up tables for full or partial updates
func (d *Device) SetLUT(fullUpdate bool) {
	lut := lutFullUpdate
	if !fullUpdate {
		lut = lutPartialUpdate
	}

	d.SendCommand(WRITE_LUT_REGISTER)
	for i := 0; i < 153; i++ {
		d.SendData(lut[i])
	}
	d.WaitUntilIdle()

	d.SendCommand(0x3f)
	d.SendData(lut[153])
	d.SendCommand(0x03) // gate voltage
	d.SendData(lut[154])
	d.SendCommand(0x04)  // source voltage
	d.SendData(lut[155]) // VSH
	d.SendData(lut[156]) // VSH2
	d.SendData(lut[157]) // VSL
	d.SendCommand(0x2c)  // VCOM
	d.SendData(lut[158])
}

// SetPixel modifies the internal buffer in a single pixel.
// The display have 2 colors: black and white
// We use RGBA(0,0,0, 255) as white (transparent)
// Anything else as black
func (d *Device) SetPixel(x int16, y int16, c color.RGBA) {
	x, y = d.xy(x, y)
	if x < 0 || x >= d.logicalWidth || y < 0 || y >= d.height {
		return
	}
	byteIndex := (int32(x) + int32(y)*int32(d.logicalWidth)) / 8
	if c.R == 0 && c.G == 0 && c.B == 0 { // TRANSPARENT / WHITE
		d.buffer[byteIndex] |= 0x80 >> uint8(x%8)
	} else { // WHITE / EMPTY
		d.buffer[byteIndex] &^= 0x80 >> uint8(x%8)
	}
}

// Display sends the buffer to the screen.
func (d *Device) Display() error {
	d.setMemoryArea(0, 0, d.logicalWidth-1, d.height-1)
	for j := int16(0); j < d.height; j++ {
		d.setMemoryPointer(0, j)
		d.SendCommand(WRITE_RAM)
		for i := int16(0); i < d.logicalWidth/8; i++ {
			d.SendData(d.buffer[i+j*(d.logicalWidth/8)])
		}
	}

	d.SendCommand(DISPLAY_UPDATE_CONTROL_2)
	d.SendData(0xC7)
	d.SendCommand(MASTER_ACTIVATION)
	d.WaitUntilIdle()
	return nil
}

// ClearDisplay erases the device SRAM
func (d *Device) ClearDisplay() {
	d.setMemoryArea(0, 0, d.logicalWidth-1, d.height-1)
	d.setMemoryPointer(0, 0)
	d.SendCommand(WRITE_RAM)
	for i := uint32(0); i < d.bufferLength; i++ {
		d.SendData(0xFF)
	}
	d.Display()
}

// setMemoryArea sets the area of the display that will be updated
func (d *Device) setMemoryArea(x0 int16, y0 int16, x1 int16, y1 int16) {
	d.SendCommand(SET_RAM_X_ADDRESS_START_END_POSITION)
	d.SendData(uint8((x0 >> 3) & 0xFF))
	d.SendData(uint8((x1 >> 3) & 0xFF))
	d.SendCommand(SET_RAM_Y_ADDRESS_START_END_POSITION)
	d.SendData(uint8(y0 & 0xFF))
	d.SendData(uint8((y0 >> 8) & 0xFF))
	d.SendData(uint8(y1 & 0xFF))
	d.SendData(uint8((y1 >> 8) & 0xFF))
}

// setMemoryPointer moves the internal pointer to the speficied coordinates
func (d *Device) setMemoryPointer(x int16, y int16) {
	d.SendCommand(SET_RAM_X_ADDRESS_COUNTER)
	d.SendData(uint8((x >> 3) & 0xFF))
	d.SendCommand(SET_RAM_Y_ADDRESS_COUNTER)
	d.SendData(uint8(y & 0xFF))
	d.SendData(uint8((y >> 8) & 0xFF))
	d.WaitUntilIdle()
}

// WaitUntilIdle waits until the display is ready
func (d *Device) WaitUntilIdle() {
	for d.busy.Get() {
		time.Sleep(100 * time.Millisecond)
	}
}

// IsBusy returns the busy status of the display
func (d *Device) IsBusy() bool {
	return d.busy.Get()
}

// ClearBuffer sets the buffer to 0xFF (white)
func (d *Device) ClearBuffer() {
	for i := uint32(0); i < d.bufferLength; i++ {
		d.buffer[i] = 0xFF
	}
}

// Size returns the current size of the display.
func (d *Device) Size() (w, h int16) {
	if d.rotation == ROTATION_90 || d.rotation == ROTATION_270 {
		return d.height, d.logicalWidth
	}
	return d.logicalWidth, d.height
}

// SetRotation changes the rotation (clock-wise) of the device
func (d *Device) SetRotation(rotation Rotation) {
	d.rotation = rotation
}

// xy chages the coordinates according to the rotation
func (d *Device) xy(x, y int16) (int16, int16) {
	switch d.rotation {
	case NO_ROTATION:
		return x, y
	case ROTATION_90:
		return d.width - y - 1, x
	case ROTATION_180:
		return d.width - x - 1, d.height - y - 1
	case ROTATION_270:
		return y, d.height - x - 1
	}
	return x, y
}
