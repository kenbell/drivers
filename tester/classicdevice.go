package tester

// I2CClassicDevice represents a mock I2C device on a mock I2C bus.
type I2CClassicDevice struct {
	c Failer
	// addr is the i2c device address.
	addr uint8
	// Registers holds the device registers. It can be inspected
	// or changed as desired for testing.
	Registers map[uint8][]uint8

	// If err is non-nil, it will be returned as the error from the
	// I2C methods.
	Err error
}

// NewI2CDevice returns a new mock I2C device.
func NewI2CClassicDevice(c Failer, addr uint8) *I2CClassicDevice {
	return &I2CClassicDevice{
		c:    c,
		addr: addr,
	}
}

// Addr returns the Device address.
func (d *I2CClassicDevice) Addr() uint8 {
	return d.addr
}

// ReadRegister implements I2C.ReadRegister.
func (d *I2CClassicDevice) ReadRegister(r uint8, buf []byte) error {
	if d.Err != nil {
		return d.Err
	}
	d.assertRegisterRange(r, buf)
	copy(buf, d.Registers[r])
	return nil
}

// WriteRegister implements I2C.WriteRegister.
func (d *I2CClassicDevice) WriteRegister(r uint8, buf []byte) error {
	if d.Err != nil {
		return d.Err
	}
	d.assertRegisterRange(r, buf)
	copy(d.Registers[r], buf)
	return nil
}

// assertRegisterRange asserts that reading or writing the given
// register and subsequent registers is in range of the available registers.
func (d *I2CClassicDevice) assertRegisterRange(r uint8, buf []byte) {
	reg, ok := d.Registers[r]

	if !ok || reg == nil {
		d.c.Fatalf("register read/write [%#x] unknown register", r)
	}

	if len(buf) != len(reg) {
		d.c.Fatalf("register read/write [%#x] bad length (%#x vs %#x)", r, len(buf), len(reg))
	}
}
