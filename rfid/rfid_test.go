package rfid

import "testing"
import "fmt"
import "time"

func TestFunc(t *testing.T) {
	sign := 0
	registerHandle("cmd", func(cmd string) {
		fmt.Println(cmd)
		sign = 1
	})
	serialDevice.Write([]byte("cmd"))
	time.Sleep(10 * time.Millisecond)
	if sign == 0 {
		t.Fatal("failed")
	}
}
