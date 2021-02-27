package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	serial "github.com/tarm/goserial"
	"github.com/zmb3/spotify"
)

var settings struct {
	Device       string        `json:"device"`
	BaudRate     int           `json:"baudRate"`
	PollInterval time.Duration `json:"pollInterval"`
}

const logo = `
███    ██  ██████  ██     ██                           
████   ██ ██    ██ ██     ██                           
██ ██  ██ ██    ██ ██  █  ██                           
██  ██ ██ ██    ██ ██ ███ ██                           
██   ████  ██████   ███ ███                            
                                                       
                                                       
██████  ██       █████  ██    ██ ██ ███    ██  ██████  
██   ██ ██      ██   ██  ██  ██  ██ ████   ██ ██       
██████  ██      ███████   ████   ██ ██ ██  ██ ██   ███ 
██      ██      ██   ██    ██    ██ ██  ██ ██ ██    ██ 
██      ███████ ██   ██    ██    ██ ██   ████  ██████  
`

func main() {
	config()

	if settings.Device == "" || settings.BaudRate == 0 {
		fatal("baudRate and / or device specified in config file. Cannot continue.", "")
	}

	fmt.Println(logo)

	ser, client := connect()

	// Give Arduino adevice few seconds to catch up.
	time.Sleep(5 * time.Second)

	start(ser, client)
}

func connect() (io.ReadWriteCloser, *spotify.Client) {
	return serialConnect(), spotifyConnect()
}

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotify.NewAuthenticator(redirectURI,
		spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		info(fmt.Sprintf("Please navigate to the following URL in a browser: %s", url))
	}

	if err != nil {
		log.Fatal(err)
	}
}

func spotifyConnect() *spotify.Client {
	// Start a callback HTTP server.
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		info(fmt.Sprintf("Got request for: %s", r.URL.String()))
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	openBrowser(url)

	// Wait for auth to complete.
	client := <-ch

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	info(fmt.Sprintf("Successfully authorized as: %s.", user.ID))
	return client
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)

	if err != nil {
		http.Error(w, "Couldn't get token.", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Authorization successful, you may close this window.")
	ch <- &client
}

func serialConnect() io.ReadWriteCloser {
	serialConf := &serial.Config{Name: settings.Device, Baud: settings.BaudRate}

	ser, err := serial.OpenPort(serialConf)

	if err != nil {
		fatal("Could not connect to serial port.", err.Error())
	}

	info("Arduino connected successfully.")

	return ser
}

func start(ser io.ReadWriteCloser, client *spotify.Client) {
	// Poll for current status, then start timer.
	elapsed(ser, client)

	ticker := time.NewTicker(time.Second * settings.PollInterval)
	for _ = range ticker.C {
		elapsed(ser, client)
	}
}

func elapsed(ser io.ReadWriteCloser, client *spotify.Client) {
	current, err := client.PlayerCurrentlyPlaying()

	if err != nil {
		log.Fatal(err)
	}

	if !current.Playing {
		convertSend(ser, "Spotify - player idle.")
		return
	}

	convertSend(ser, fmt.Sprintf("%s - %s",
		current.Item.Artists[0].Name, current.Item.Name))
}

var lastMsg string = "Artist - Track"

func convertSend(ser io.ReadWriteCloser, msg string) {
	if msg == "" || lastMsg == msg {
		return
	}

	info(msg)
	send(ser, strings.Replace(msg, " - ", "<~>", 1))
	lastMsg = msg
}

const configFilePath string = "config/config.json"

func config() {
	configFile, err := os.Open(configFilePath)

	if err != nil {
		fatal("Couldn't open config file.", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&settings); err != nil {
		fatal("Couldn't parse config file.", err.Error())
	}
}

func fatal(msg string, err string) {
	log.Fatal(fmt.Printf("%s FULL ERROR: %s \n\nExiting.\n\n", msg, err))
}

func info(msg string) {
	log.Println(msg)
}

func send(ser io.ReadWriteCloser, msg string) {
	ser.Write([]byte(msg))
}
