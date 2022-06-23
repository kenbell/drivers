// Package ssd1681 implements a driver for the SSD1681 e-paper displays.
//
// This package is based on the Adafruit Circuit Python implementation
// for SSD1681: https://github.com/adafruit/Adafruit_CircuitPython_EPD/blob/main/adafruit_epd/ssd1681.py
//
// The implementation uses a simplistic buffer in MCU RAM (5kB for 200x200
// display).  The Adafruit SSD1681 board also includes SPI RAM, but that is
// not used in this implementation to keep it generic.
//
package ssd1681

import (
	"errors"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
)

var (
	errDrawingOutOfBounds = errors.New("rectangle coordinates outside display area")
)

type Device struct {
	bus      drivers.SPI
	dcPin    machine.Pin
	resetPin machine.Pin
	csPin    machine.Pin
	enPin    machine.Pin
	busyPin  machine.Pin
	width    int16
	height   int16
	buffer   []byte
}

// Config is the configuration for the display
type Config struct {
	Width  int16
	Height int16
}

// New creates a new SSD1351 connection. The SPI wire must already be configured.
func New(bus drivers.SPI, resetPin, dcPin, csPin, enPin, busyPin machine.Pin) Device {
	return Device{
		bus:      bus,
		dcPin:    dcPin,
		resetPin: resetPin,
		csPin:    csPin,
		enPin:    enPin,
		busyPin:  busyPin,
	}
}

// Configure initializes the display with default configuration
func (d *Device) Configure(cfg Config) {
	if cfg.Width == 0 {
		cfg.Width = 200
	}

	if cfg.Height == 0 {
		cfg.Height = 200
	}

	d.width = cfg.Width
	d.height = cfg.Height

	d.buffer = make([]byte, ((d.width+7)/8)*d.height)

	// configure GPIO pins
	d.dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.resetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.enPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.busyPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	d.enPin.High()

	// reset the device
	d.resetPin.High()
	time.Sleep(100 * time.Millisecond)
	d.resetPin.Low()
	time.Sleep(100 * time.Millisecond)
	d.resetPin.High()
	time.Sleep(200 * time.Millisecond)
	d.Wait()
	d.Command(SW_RESET)
	d.Wait()

	// Initialization

	// Gate settings: width gates, sequence G0 G1 ..., non-interlaced, from G0-G199
	d.Command(DRIVER_CONTROL, byte(d.width-1), byte((d.width-1)>>8), 0)

	// RAM address: increment in X & Y, update in Y
	d.Command(DATA_MODE, 0x3)

	// RAM X start/end: start 0
	d.Command(SET_RAMXPOS, 0, byte(ceil(d.width, 8)))

	// RAM Y start/end: start 0, end from height
	d.Command(SET_RAMYPOS, 0, 0, byte(d.height-1), byte((d.height-1)>>8))

	// Set border waveform: GS transition, follow LUT 1
	d.Command(WRITE_BORDER, 0x05)

	// Temp sensor: use internal temp sensor
	d.Command(TEMP_CONTROL, 0x80)

	d.Wait()
}

// SetPixel sets a pixel in the buffer
func (d *Device) SetPixel(x int16, y int16, c color.RGBA) {
	if x < 0 || y < 0 || x >= d.width || y >= d.height {
		return
	}
	d.FillRectangle(x, y, 1, 1, c)
}

// FillRectangle fills a portion of the display with a color
func (d *Device) FillRectangle(x, y, width, height int16, c color.RGBA) error {
	if x < 0 || y < 0 || width <= 0 || height <= 0 ||
		x >= d.width || (x+width) > d.width || y >= d.height || (y+height) > d.height {
		return errDrawingOutOfBounds
	}

	// Convert to B/W via grayscale luminance
	lum := color.GrayModel.Convert(c).(color.Gray).Y
	colBit := byte(255) // 'white'
	if lum < 128 {
		colBit = 0 // 'black'
	}

	for ypos := y; ypos < y+height; ypos++ {
		// This code could probably be a lot more time-efficient for byte-aligned
		// rectangles, but for now this is simple and works :)
		for xpos := x; xpos < x+width; xpos++ {
			bufPos := (ypos * ceil(d.width, 8)) + xpos/8
			mask := byte(1 << (7 - (xpos % 8)))

			d.buffer[bufPos] = (d.buffer[bufPos] & ^mask) | (colBit & mask)
		}
	}

	return nil
}

// Display populates the hardware RAM from the MCU framebuffer, then
// tells the hardware to do a full refresh.
func (d *Device) Display() error {
	d.Command(SET_RAMXCOUNT, 0)
	d.Command(SET_RAMYCOUNT, 0, 0)

	for row := 0; row < int(d.height); row++ {
		d.Command(WRITE_BWRAM, d.buffer[row*int(ceil(d.width, 8)):(row+1)*int(ceil(d.width, 8))]...)
	}

	d.Command(DISP_CTRL2, 0xF7)
	d.Command(MASTER_ACTIVATE)
	d.Wait()

	return nil
}

// Size returns the current size of the display
func (d *Device) Size() (w, h int16) {
	return d.width, d.height
}

// Sleep puts the display to sleep, to re-awaken call
// Configure to do a hardware reset.
func (d *Device) DeepSleep() {
	d.Command(DEEP_SLEEP, 0x01)
}

// Command sends a command byte to the display
func (d *Device) Command(command byte, data ...byte) {
	d.dcPin.Set(false)
	d.csPin.Low()
	d.bus.Transfer(command)
	d.dcPin.Set(true)
	d.bus.Tx(data, nil)
	d.csPin.High()
}

// Waits until device signals it is ready
func (d *Device) Wait() {
	for d.busyPin.Get() {
		time.Sleep(10 * time.Millisecond)
	}
}

func ceil(n int16, d int16) int16 {
	return (n + d - 1) / d
}
