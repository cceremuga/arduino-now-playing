package main

import (
	"fmt"
	"github.com/tarm/goserial"
	"log"
	"io"
)

var portName string = "/dev/tty.usbserial-AM01VD4K"
var baudRate int = 9600
var playerType int = 1

func main() {
	//banner print
	displayBanner()

	if portName != "" && baudRate > 0 {
		serialConf := &serial.Config{Name: portName, Baud: baudRate}

		ser, err := serial.OpenPort(serialConf)

		if err != nil {
			log.Fatal(err)
		}

		sendToSerial(ser, "This is a test ")
	}
}

func sendToSerial(ser io.ReadWriteCloser, msg string) {
	ser.Write([]byte(msg))
}

func displayBanner() {
	fmt.Println("    _   __                 ____  __            _            ")
	fmt.Println("   / | / /___ _      __   / __ \\/ /___ ___  __(_)___  ____ _")
	fmt.Println("  /  |/ / __ \\ | /| / /  / /_/ / / __ `/ / / / / __ \\/ __ `/")
	fmt.Println(" / /|  / /_/ / |/ |/ /  / ____/ / /_/ / /_/ / / / / / /_/ / ")
	fmt.Println("/_/ |_/\\____/|__/|__/  /_/   /_/\\__,_/\\__, /_/_/ /_/\\__, /  ")
	fmt.Println("   __           _____           _    /____/        /____/   ")
	fmt.Println("  / /_____     / ___/___  _____(_)___ _/ /                  ")
	fmt.Println(" / __/ __ \\    \\__ \\/ _ \\/ ___/ / __ `/ /                   ")
	fmt.Println("/ /_/ /_/ /   ___/ /  __/ /  / / /_/ / /                    ")
	fmt.Println("\\__/\\____/   /____/\\___/_/  /_/\\__,_/_/                     \n\n")
}
