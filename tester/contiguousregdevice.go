package tester

// MaxRegisters is the maximum number of registers supported for a Device.
const MaxRegisters = 200

// I2CContiguousRegisterDevice represents a mock I2C device on a mock I2C bus.
//
// This models the device registers as single-byte registers in a contiguous
// address space where the address auto-increments on read.
type I2CContiguousRegisterDevice struct {
	c Failer
	// addr is the i2c device address.
	addr uint8
	// Registers holds the device registers. It can be inspected
	// or changed as desired for testing.
	Registers [MaxRegisters]uint8
	// If Err is non-nil, it will be returned as the error from the
	// I2C methods.
	Err error
}

// NewI2CContiguousRegisterDevice returns a new mock I2C device.
func NewI2CContiguousRegisterDevice(c Failer, addr uint8) *I2CContiguousRegisterDevice {
	return &I2CContiguousRegisterDevice{
		c:    c,
		addr: addr,
	}
}

// Addr returns the Device address.
func (d *I2CContiguousRegisterDevice) Addr() uint8 {
	return d.addr
}

// ReadRegister implements I2C.ReadRegister.
func (d *I2CContiguousRegisterDevice) ReadRegister(r uint8, buf []byte) error {
	if d.Err != nil {
		return d.Err
	}
	d.assertRegisterRange(r, buf)
	copy(buf, d.Registers[r:])
	return nil
}

// WriteRegister implements I2C.WriteRegister.
func (d *I2CContiguousRegisterDevice) WriteRegister(r uint8, buf []byte) error {
	if d.Err != nil {
		return d.Err
	}
	d.assertRegisterRange(r, buf)
	copy(d.Registers[r:], buf)
	return nil
}

// assertRegisterRange asserts that reading or writing the given
// register and subsequent registers is in range of the available registers.
func (d *I2CContiguousRegisterDevice) assertRegisterRange(r uint8, buf []byte) {
	if int(r) >= len(d.Registers) {
		d.c.Fatalf("register read/write [%#x, %#x] start out of range", r, int(r)+len(buf))
	}
	if int(r)+len(buf) > len(d.Registers) {
		d.c.Fatalf("register read/write [%#x, %#x] end out of range", r, int(r)+len(buf))
	}
}
