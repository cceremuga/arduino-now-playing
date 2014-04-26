package main

import (
	"fmt"
	"github.com/tarm/goserial"
	"log"
	"io"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
	"net/http"
	"strings"
)

var configFilePath string = "config.json"

var settings struct {
	PortName string `json:"portName"`
	BaudRate int `json:"baudRate"`
	PlayerType int `json:"playerType"`
	PollRateSeconds time.Duration `json:"pollRateSeconds"`
	VlcWebUrl string `json:"vlcWebUrl"`
	VlcWebPassword string `json:"vlcWebPassword"`
}

type VlcNowPlaying struct {
	State string `json:"state"`
	Information struct {
		Category struct {
			Meta struct {
				NowPlaying string `json:"now_playing"`
			}
		}
	}

	NowPlaying string `json:"information.category.meta.now_playing"`
}

var lastSentMessage string = "Artist - Track"

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

		//start our X seconds timer
		startTicker(ser)
	} else {
		//something is not configured, back out.
		endEarly("No baudRate and / or portName specified in config file. Cannot continue.", "")
	}
}

//kicks off our ticker, fires the elapsed once to start
func startTicker(ser io.ReadWriteCloser) {
	infoMessage("Ticker initialized.")

	tickerElapsed(ser)

	ticker := time.NewTicker(time.Second * settings.PollRateSeconds)
	for _ = range ticker.C {
		tickerElapsed(ser)
	}
}

//determines which player to poll, then acts
func tickerElapsed(ser io.ReadWriteCloser) {
	if settings.PlayerType == 1 {
		pollVlc(ser)
	} else {
		pollSpotify(ser)
	}
}

//poll VLC for now playing info
func pollVlc(ser io.ReadWriteCloser) {
	content, err := getResponseContent(settings.VlcWebUrl, settings.VlcWebPassword)

	if err != nil {
		infoMessage("Unable to get data from VLC. Double check your url and password.")
		infoMessage(err.Error())
	} else {
		//deserialize jSON
		var vlcStatus VlcNowPlaying
		err = json.Unmarshal(content, &vlcStatus)

		if err != nil {
			infoMessage("Unable to parse data returned from VLC.")
			infoMessage(err.Error())
		} else {
			convertForArduino(ser, vlcStatus.Information.Category.Meta.NowPlaying)
		}
	}
}

//convert to format our arduino sketch is expecting
func convertForArduino(ser io.ReadWriteCloser, msg string) {
	if msg != "" && lastSentMessage != msg {
		sendToSerial(ser, strings.Replace(msg, " - ", "<~>", -1))

		lastSentMessage = msg
	}
}

//poll Spotify for now playing info
func pollSpotify(ser io.ReadWriteCloser) {
	infoMessage("Spotify polling not implemented yet.")
}

// Code adapted from http://www.codingcookies.com/2013/03/21/consuming-json-apis-with-go/
// This function fetch the content of a URL will return it as an
// array of bytes if retrieved successfully.
func getResponseContent(url string, password string) ([]byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)

	if password != "" {
		req.SetBasicAuth("", password)
	}

	if err != nil {
		return nil, err
	}

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Read the content into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// At this point we're done - simply return the bytes
	return body, nil
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

//prints a fatal error message which in turn exits immediately
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

	infoMessage(msg)
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
