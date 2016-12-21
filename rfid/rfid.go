/*
*	package rfid, 用于与rfid进行串口数据通信，完成指令
 */
package rfid

import (
	"log"
	"time"

	"github.com/tarm/serial"
)

//串口设备
var SerialDevice *serial.Port

//package初始化
func init() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200, ReadTimeout: time.Nanosecond * 1}
	var err error
	SerialDevice, err = serial.OpenPort(c)
	log.Println(c.ReadTimeout.Nanoseconds())
	if err != nil {
		log.Println(err)
		return
	}
	// serialHandle()
}

//指令处理
func cmdHandle(cmd []byte) {
}

func read() []byte {
	b := make([]byte, 1)
	buf := make([]byte, 0)
	for {
		n, err := SerialDevice.Read(b)
		if err != nil {
			// log.Println(err)
			break
		}
		if n == 0 {
			log.Println("data length == 0")
			break
		}
		buf = append(buf, b[0])
	}
	return buf
}

func wait() []byte {
	/* timer := time.NewTimer(time.Second * 2) */
	// timeout := false
	// go func() {
	//     <-timer.C
	//     timeout = true
	// }()
	// if timeout == true {
	//     return nil
	/* } */
	return read()
}

//查找14443卡片
func SearchCard14443() {
	cmd := []byte{0xFF, 0xFE, 0x03, 0x00, 0x20, 0x20}
	time.Sleep(time.Second)
	SerialDevice.Write(cmd)
	buf := wait()
	log.Println("card", buf)
}

func Auth14443() {
	cmd := []byte{0xFF, 0x0FE, 0x03, 0x08, 0x23, 0x01, 0x60, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x86}
	SerialDevice.Write(cmd)
	time.Sleep(time.Second / 3)
	buf := wait()
	log.Printf("Auth cmd. Receive auth data from serial port: %X", buf)
}

//read data from section block
func Read14443(section byte, block byte) {
	cmd := []byte{0xFF, 0xFE, 0x03, 0x02, 0x21, section, block, 0x00}
	sum(cmd)
	log.Printf("cmd:%X\n", cmd)
	SerialDevice.Write(cmd)
	time.Sleep(time.Second / 5)
	buf := wait()
	log.Printf("Receive read data from serial port: %X\n", buf)
}

//write data to section block
func Write14443(section byte, block byte, data [16]byte) {
	cmd := make([]byte, 24)
	b := []byte{0xFF, 0xFE, 0x03, 0x12, 0x22, section, block}

	for i := 0; i < len(b); i++ {
		cmd[i] = b[i]
	}

	for i := 0; i < 16; i++ {
		cmd[len(b)+i] = data[i]
	}
	sum(cmd)
	log.Printf("cmd:%X\n", cmd)
	SerialDevice.Write(cmd)
	time.Sleep(time.Second / 5)
	buf := wait()
	log.Printf("Receive write data from serial port: %X", buf)
}

//calculate sum of slice
func sum(cmd []byte) {
	var sum byte = 0
	for _, v := range cmd[0 : cap(cmd)-1] {
		sum += v
	}
	cmd[cap(cmd)-1] = sum
}
