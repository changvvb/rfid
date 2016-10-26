package rfid

import (
	"os"
)

func writeCOMstring(com string, data string) error {
	f, err := os.Open(com)
	if err != nil {
		return err
	}
	f.WriteString(string)
}
