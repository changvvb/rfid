package rfid

import "testing"

func TestWriteCOMstring(t *testing.T) {
	_, err := WriteCOMstring("/dev/ttyUSB0", "ddffdfd")
	if err != nil {
		t.Error(err)
	}
}
