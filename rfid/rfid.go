/*
*	package rfid, 用于与rfid进行串口数据通信，完成指令
 */
package rfid

import (
	"fmt"
	"github.com/tarm/serial"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

//串口设备
var SerialDevice *serial.Port

var SearchCardCallBack func([]byte)

//package初始化
func init0() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200, ReadTimeout: time.Nanosecond * 1}
	var err error
	SerialDevice, err = serial.OpenPort(c)
	log.Println(c.ReadTimeout.Nanoseconds())
	if err != nil {
		log.Println(err)
		return
	}
}

func BoolReady() bool {
	if SerialDevice != nil {
		return true
	}
	return false
}

//return serial devices
//true: only ttyUSB*, false: tty*
func ListDevice(flag bool) []string {
	files, err := ioutil.ReadDir("/dev/")
	if err != nil {
		return nil
	} else {
		var s []string
		for _, v := range files {
			if flag {
				if strings.Contains(v.Name(), "ttyUSB") {
					s = append(s, v.Name())
				}
			} else {
				if strings.Contains(v.Name(), "tty") {
					s = append(s, v.Name())
				}
			}
		}
		return s
	}
}

func Connect(port string) error {
	if !(strings.Contains(port, "com") || strings.Contains(port, "com") || strings.Contains(port, "/dev")) {
		port = "/dev/" + port
	}
	c := &serial.Config{Name: port, Baud: 115200, ReadTimeout: time.Nanosecond * 1}
	var err error
	SerialDevice, err = serial.OpenPort(c)
	return err
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
func SearchCard14443() []byte {
	cmd := []byte{0xFF, 0xFE, 0x03, 0x00, 0x20, 0x20}
	SerialDevice.Write(cmd)
	time.Sleep(time.Second / 3)
	buf := wait()
	log.Println("card", buf)

	if len(buf) < 7 {
		return nil
	}
	log.Println("length:----------------------------->", len(buf))

	if buf[3] == 0x04 {
		return buf[6 : len(buf)-1]
	}
	return nil
}

var BoolAutoSearch bool

func AutoSearch14443() {
	if BoolAutoSearch == true {
		return
	}
	BoolAutoSearch = true
	go func() {
		for {
			if BoolAutoSearch {
				if buf := SearchCard14443(); SearchCardCallBack != nil {
					SearchCardCallBack(buf)
				}
				time.Sleep(time.Second / 50)
			} else {
				return
			}
		}
	}()
}

func StopAutoSearch14443() {
	BoolAutoSearch = false
}

func Auth14443(section byte) {
	cmd := []byte{0xFF, 0x0FE, 0x03, 0x08, 0x23, section, 0x60, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00}
	sum(cmd)
	SerialDevice.Write(cmd)
	time.Sleep(time.Second / 3)
	buf := wait()
	log.Printf("Auth cmd. Receive auth data from serial port: %X", buf)
}

//read data from section block
func Read14443(section byte, block byte) ([]byte, error) {
	cmd := []byte{0xFF, 0xFE, 0x03, 0x02, 0x21, section, block, 0x00}
	sum(cmd)
	log.Printf("cmd:%X\n", cmd)
	SerialDevice.Write(cmd)
	time.Sleep(time.Second / 50)
	buf := wait()
	log.Printf("Receive read data from serial port: %X\n", buf)

	//校验不通过
	if !check(buf) {
		return nil, fmt.Errorf("Check sum failed")
	}

	if buf[3] == 0 {
		return nil, fmt.Errorf("Error! NO CARD!")
	}

	return buf[12 : 12+16], nil
}

func ReadString14443(section byte, block byte) string {
	if b, err := Read14443(section, block); err == nil {
		for i, v := range b {
			if v == 0 {
				b = b[0:i]
				break
			}
		}
		return string(b)
	} else {
		return ""
	}
}

//write data to section block
func Write14443(section byte, block byte, data []byte) {
	cmd := make([]byte, 24)
	b := []byte{0xFF, 0xFE, 0x03, 0x12, 0x22, section, block}

	for i := 0; i < len(b); i++ {
		cmd[i] = b[i]
	}

	for i := 0; i < 16; i++ {
		if i >= len(data) {
			cmd[len(b)+i] = 0
			continue
		}
		cmd[len(b)+i] = data[i]
	}
	sum(cmd)
	log.Printf("cmd:%X\n", cmd)
	SerialDevice.Write(cmd)
	time.Sleep(time.Second / 50)
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

//check data from serial port
func check(cmd []byte) bool {
	var sum byte = 0
	if len(cmd) == 0 {
		return false
	}
	for _, v := range cmd[0 : len(cmd)-1] {
		sum += v
	}

	if cmd[len(cmd)-1] == sum {
		return true
	} else {
		return false
	}
}
