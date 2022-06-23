// This example shows how to use sx127x LoRa modules
// Tested on Feather-RP2040 https://www.adafruit.com/product/4884
//
// An example of the sx127x is https://www.adafruit.com/product/3072
//

package main

import (
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/sx127x"
)

// Pins available on feather-rp2040
const (
	resetPin = machine.GPIO13
	csPin    = machine.GPIO9
	dio0Pin  = machine.LORA_DIO0
)

func main() {

	time.Sleep(100 * time.Millisecond)

	println("starting tx")

	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 1000000, // 1MHz (10MHz should be possible)
	})

	modem := sx127x.New(machine.SPI1, resetPin, csPin, dio0Pin)

	err := modem.Configure(sx127x.Config{
		Frequency: 868100000, // 868.1MHz
		CRC:       sx127x.CrcModeOn,
	})

	if err != nil {
		panic(err)
	}

	if !modem.Detect() {
		panic("failed to detect sx127x")
	}

	println("detected sx127x")

	ticker := time.NewTicker(time.Second)

	count := 0
	for {
		err = modem.LoraTx([]byte{0xA0, 0xA1, 0xA2, 0xA3}, 1000)
		if err != nil {
			panic(err)
		}

		count++

		println("tx complete: " + strconv.Itoa(count))

		// Wait for ticker
		select {
		case <-ticker.C:
		}
	}
}
