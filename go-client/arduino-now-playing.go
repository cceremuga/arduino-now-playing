package main

import (
	"fmt"
	"github.com/tarm/goserial"
	"log"
	"io"
	"os"
	"encoding/json"
	"time"
)

var configFilePath string = "config.json"

var settings struct {
	PortName string `json:"portName"`
	BaudRate int `json:"baudRate"`
	PlayerType int `json:"playerType"`
	PollRateSeconds time.Duration `json:"pollRateSeconds"`
}

func main() {
	//banner print
	displayBanner()

	//load config jSON into struct
	loadConfig()

	//verify port name and baud rate set before continuing
	if settings.PortName != "" && settings.BaudRate > 0 {
		//config & open serial connection
		serialConf := &serial.Config{Name: settings.PortName, Baud: settings.BaudRate}

		ser, err := serial.OpenPort(serialConf)

		if err != nil {
			endEarly("Could not connect to serial port. Cannot continue.", err.Error())
		} else {
			infoMessage("Connected to serial port successfully.")
		}

		sendToSerial(ser, "Test")

		//start our X seconds timer
		startTicker()
	} else {
		//something is not configured, back out.
		endEarly("No baudRate and / or portName specified in config file. Cannot continue.", "")
	}
}

//kicks off our ticker, fires the elapsed once to start
func startTicker() {
	tickerElapsed()

	ticker := time.NewTicker(time.Second * settings.PollRateSeconds)
    for _ = range ticker.C {
    	tickerElapsed()
	}
}

//determines which player to poll, then acts
func tickerElapsed() {
	infoMessage("Test timer elapsed")
}

//loads a jSON config file, parses it into a struct
func loadConfig() {
	configFile, err := os.Open(configFilePath)

	if err != nil {
		endEarly("Couldn't open config file. Cannot continue.", err.Error())
	} else{
		infoMessage("Opened config file successfully.")
	}

	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&settings); err != nil {
		endEarly("Couldn't parse config file. Cannot continue.", err.Error())
	} else {
		infoMessage("Parsed config file successfully.")
	}
}

//prints a fatal error message, then waits for a keypress prior to exit
func endEarly(msg string, err string) {
	log.Fatal(msg, " FULL ERROR: ", err, "\n\nExiting immediately.\n\n")
}

//logs a standard string message
func infoMessage(msg string) {
	log.Print(msg)
}

//sends string data to the serial port passed in
func sendToSerial(ser io.ReadWriteCloser, msg string) {
	ser.Write([]byte(msg))
}

//prints a welcome banner to the start of the app
func displayBanner() {
	fmt.Println("\n    _   __                 ____  __            _            ")
	fmt.Println("   / | / /___ _      __   / __ \\/ /___ ___  __(_)___  ____ _")
	fmt.Println("  /  |/ / __ \\ | /| / /  / /_/ / / __ `/ / / / / __ \\/ __ `/")
	fmt.Println(" / /|  / /_/ / |/ |/ /  / ____/ / /_/ / /_/ / / / / / /_/ / ")
	fmt.Println("/_/ |_/\\____/|__/|__/  /_/   /_/\\__,_/\\__, /_/_/ /_/\\__, /  ")
	fmt.Println("   __           _____           _    /____/        /____/   ")
	fmt.Println("  / /_____     / ___/___  _____(_)___ _/ /                  ")
	fmt.Println(" / __/ __ \\    \\__ \\/ _ \\/ ___/ / __ `/ /                   ")
	fmt.Println("/ /_/ /_/ /   ___/ /  __/ /  / / /_/ / /                    ")
	fmt.Println("\\__/\\____/   /____/\\___/_/  /_/\\__,_/_/                     \n")
}
