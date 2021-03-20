// Example using i2c BME280 environment sensor with 128x32 SSD1306 display.
//
package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/bme280"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/st7735"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

const (
	spiScreenResetPin = machine.A0
	spiScreenDcPin    = machine.A1
	spiScreenCsPin    = machine.A2
	spiScreenBlPin    = machine.A3
)

func main() {
	// Sleep to allow power to settle (needed on some STM32 boards)
	time.Sleep(100 * time.Millisecond)

	machine.I2C0.Configure(machine.I2CConfig{})

	sensor := bme280.New(machine.I2C0)
	sensor.Configure()

	screen := ssd1306.NewI2C(machine.I2C0)
	screen.Configure(ssd1306.Config{
		Address: ssd1306.Address_128_32,
	})

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
	})
	spiscreen := st7735.New(machine.SPI0, spiScreenResetPin, spiScreenDcPin, spiScreenCsPin, spiScreenBlPin)
	spiscreen.Configure(st7735.Config{
		Width:  128,
		Height: 160,
	})

	connected := sensor.Connected()
	if !connected {
		println("BME280 not detected")
		return
	}
	println("BME280 detected")

	i := 0
	spiscreen.FillScreen(color.RGBA{0, 0, 0, 0})
	for {
		temp, _ := sensor.ReadTemperature()
		pressure, _ := sensor.ReadPressure()

		// Draw temp and pressure to display
		screen.ClearBuffer()

		// Avoid using fmt package (for small devices) so divide
		// 'less' than required (by number of fraction digits)
		// and do custom formatting
		status := fmtD(uint32(temp/10), 2, 2) + " C  " + fmtD(uint32(pressure/10000), 4, 1) + " hPa"

		//		tinyfont.WriteLine(&screen, &proggy.TinySZ8pt7b, 0, 10, status, color.RGBA{255, 255, 255, 0})
		screen.Display()

		// Output temp and pressure to UART also
		println(status, "    ", i)
		i++

		spiscreen.FillRectangle(0, 0, 128, 20, color.RGBA{0, 0, 0, 0})
		tinyfont.WriteLine(&spiscreen, &proggy.TinySZ8pt7b, 0, 10, status, color.RGBA{255, 255, 255, 0})
		spiscreen.Display()

		time.Sleep(500 * time.Millisecond)
	}
}

func fmtD(val uint32, i int, f int) string {
	result := make([]byte, i+f+1)

	for p := len(result) - 1; p >= 0; p-- {
		result[p] = byte(itoc(val % 10))
		val = val / 10

		if p == i+1 && p > 0 {
			p--
			result[p] = '.'
		}
	}

	return string(result)
}

func itoc(v uint32) byte {
	return byte(uint32('0') + v)
}
