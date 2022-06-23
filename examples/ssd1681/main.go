// This example shows how to use ssd1681 200x200 display over SPI
// Tested on Feather-RP2040 https://www.adafruit.com/product/4884
//
// An example of the ssd1681 is https://www.adafruit.com/product/4196
// (the model from June 7th 2021)
//

package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ssd1681"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// Pins available on feather-rp2040
const (
	enPin    = machine.GPIO2
	busyPin  = machine.GPIO9
	resetPin = machine.GPIO10
	dcPin    = machine.GPIO12
	csPin    = machine.GPIO13
)

func main() {
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 20000000, // 20MHz
	})
	display := ssd1681.New(
		machine.SPI0,
		resetPin,
		dcPin,
		csPin,
		enPin,
		busyPin)

	display.Configure(ssd1681.Config{
		Width:  200,
		Height: 200,
	})

	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	// Blank display
	display.FillRectangle(0, 0, 200, 200, white)
	display.Display()

	// Say hello
	tinyfont.WriteLine(&display, &freemono.Regular9pt7b, 5, 15, "Hello from TinyGo!", black)
	display.Display()

	display.DeepSleep()
}
