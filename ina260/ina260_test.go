package ina260

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"tinygo.org/x/drivers/tester"
)

func TestDefaultI2CAddress(t *testing.T) {
	c := qt.New(t)
	bus := tester.NewI2CBus(c)
	dev := New(bus, 0)
	c.Assert(dev.Address, qt.Equals, uint16(Address))
}

func TestConnected(t *testing.T) {
	c := qt.New(t)
	bus := tester.NewI2CBus(c)
	fake := tester.NewI2CClassicDevice(c, Address)
	fake.Registers = defaultRegisters()
	bus.AddDevice(fake)

	dev := New(bus, 0)
	c.Assert(dev.Connected(), qt.Equals, true)
}

func TestVoltage(t *testing.T) {
	c := qt.New(t)
	bus := tester.NewI2CBus(c)
	fake := tester.NewI2CClassicDevice(c, Address)
	fake.Registers = defaultRegisters()
	fake.Registers[REG_BUSVOLTAGE] = []byte{0x25, 0x70}
	bus.AddDevice(fake)

	dev := New(bus, 0)
	// Datasheet: 2570h = 11.98V = 11980mV = 11980000uV
	c.Assert(dev.Voltage(), qt.Equals, int32(11980000))
}

func TestCurrent(t *testing.T) {
	c := qt.New(t)
	bus := tester.NewI2CBus(c)
	fake := tester.NewI2CClassicDevice(c, Address)
	fake.Registers = defaultRegisters()
	fake.Registers[REG_CURRENT] = []byte{0x27, 0x10}
	bus.AddDevice(fake)

	dev := New(bus, 0)
	// Datasheet: 2710h = 12.5A = 12500mA = 12500000uA
	c.Assert(dev.Current(), qt.Equals, int32(12500000))
}

func TestPower(t *testing.T) {
	c := qt.New(t)
	bus := tester.NewI2CBus(c)
	fake := tester.NewI2CClassicDevice(c, Address)
	fake.Registers = defaultRegisters()
	fake.Registers[REG_POWER] = []byte{0x3A, 0x7F}
	bus.AddDevice(fake)

	dev := New(bus, 0)
	// 3A7Fh = 149.75W = 149750mW = 149750000uW
	c.Assert(dev.Power(), qt.Equals, int32(149750000))
}

// defaultRegisters returns the default values for all of the device's registers.
// set TI INA260 datasheet for power-on defaults
func defaultRegisters() map[uint8][]uint8 {
	return map[uint8][]uint8{
		REG_CONFIG:     {0x61, 0x27},
		REG_CURRENT:    {0x00, 0x00},
		REG_BUSVOLTAGE: {0x00, 0x00},
		REG_POWER:      {0x00, 0x00},
		REG_MASKENABLE: {0x00, 0x00},
		REG_ALERTLIMIT: {0x00, 0x00},
		REG_MANF_ID:    {0x54, 0x49},
		REG_DIE_ID:     {0x22, 0x70},
	}
}
