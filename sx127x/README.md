# SX127x Driver

This driver is for the Semtech SX127x RF low-power, long-range transceivers.

The HopeRF RFM95W is a module built using the SX127x, which is available as a break-out board from Adafruit.


Datasheets:
 * https://www.hoperf.com/data/upload/portal/20190801/RFM95W-V2.0.pdf
 * https://semtech.my.salesforce.com/sfc/p/E0000000JelG/a/2R0000001Rbr/6EfVZUorrpoKFfvaF_Fkpgp5kzjiNyiAbqcpqh9qSjE

This code is based on the Circuit Python implementation, itself adapted from the Radiohead library RF95 code: https://github.com/adafruit/Adafruit_CircuitPython_RFM9x/blob/main/adafruit_rfm9x.py


## Interface Notes

The current implementation tries to be essentially compatible with the sx126x implementation.  There are some areas that might be improved in the APIs:

1. LoraRx always allocates, instead it would be good if caller could provide output buffer.  I.e. change API from `func (d *Device) LoraRx(timeoutMs uint32) ([]uint8, error)` to `func (d *Device) LoraRx(timeoutMs uint32, buf []uint8) (error)`

