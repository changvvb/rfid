package main

// import "github.com/changvvb/rfid/rfid"
import "github.com/changvvb/rfid/server"

func main() {
	s := server.New()
	//
	// for {
	//     rfid.Auth14443()
	//     rfid.Read14443(1, 1)
	//
	//     rfid.Auth14443()
	//     rfid.Write14443(1, 1, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	//     // time.Sleep(time.Second / 5)
	// }
	//
	s.Run()
}
