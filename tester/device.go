package tester

// I2CDevice represents a mock I2C device on a mock I2C bus.
//
// This is an extension of the machine.I2C interface
type I2CDevice interface {
	// Addr returns the Device address.
	Addr() uint8

	// ReadRegister implements I2C.ReadRegister.
	ReadRegister(r uint8, buf []byte) error

	// WriteRegister implements I2C.WriteRegister.
	WriteRegister(r uint8, buf []byte) error
}
