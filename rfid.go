package rfid

import (
	"github.com/tarm/serial"
)

var serialDevice *serial.Port

func init() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	var err error
	serialDevice, err = serial.OpenPort(c)
	if err != nil {
		return
	}
}
