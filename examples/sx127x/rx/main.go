// This example shows how to use sx127x LoRa modules
// Tested on Feather-RP2040 https://www.adafruit.com/product/4884
//
// An example of the sx127x is https://www.adafruit.com/product/3072
//

package main

import (
	"encoding/hex"
	"fmt"
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/sx127x"
)

// Pins available on feather-rp2040
const (
	resetPin = machine.A2
	csPin    = machine.A3
	dio0Pin  = machine.GPIO6
)

func main() {

	machine.GPIO13.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.GPIO13.Set(false)

	time.Sleep(100 * time.Millisecond)

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 1000000, // 1MHz (10MHz should be possible)
	})

	modem := sx127x.New(machine.SPI0, resetPin, csPin, dio0Pin)

	err := modem.Configure(sx127x.Config{
		Frequency: 868100000,
		CRC:       sx127x.CrcModeOff,
	})

	if err != nil {
		panic(err)
	}

	if !modem.Detect() {
		panic("failed to detect sx127x")
	}

	println("detected sx127x")

	count := 0
	for {
		buf, err := modem.LoraRx(12345)
		if err != nil {
			println(err)
			continue
		}

		if buf == nil {
			print("TO,")
			continue
		}

		count++

		if len(buf) > 0 {
			println("dump: (" + strconv.Itoa(len(buf)) + ")")
			println(fmt.Sprintf("%v", time.Now()))
			println(hex.Dump(buf))
		}
	}
}
