package main

import (
	"github.com/changvvb/rfid/rfid"
	"github.com/changvvb/rfid/server"
	"time"
)

func main() {
	s := server.New()
	for {
		rfid.Auth14443()
		rfid.Read14443(1, 1)
		rfid.Auth14443()
		rfid.Write14443(1, 1, [16]byte{5, 6, 3, 5, 3, 2, 5, 6, 9})
		time.Sleep(time.Second / 5)
	}
	s.Run()
}
