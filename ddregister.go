package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/timer"
)

const (
	USER_AGENT = "github.com/djthorpe/ddregister"
	APIFY_URI  = "https://api.ipify.org"
	GOOGLE_URI = "https://domains.google.com/nic/update"
)

var (
	client = &http.Client{}
)

func NewRequest(method, url string) (*http.Request, error) {
	if req, err := http.NewRequest(method, url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Add("User-Agent", USER_AGENT)
		return req, nil
	}
}

func Do(req *http.Request) ([]byte, string, error) {
	if resp, err := client.Do(req); err != nil {
		return nil, "", err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, "", fmt.Errorf("Error: %v: %v", resp.StatusCode, resp.Status)
		} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, "", err
		} else {
			return body, resp.Header.Get("Content-Type"), nil
		}
	}
}

func GetExternalAddress() (net.IP, error) {
	// Obtain the current external IP address
	if req, err := NewRequest("GET", APIFY_URI); err != nil {
		return nil, err
	} else if body, content_type, err := Do(req); err != nil {
		return nil, err
	} else if content_type != "text/plain" {
		return nil, fmt.Errorf("Unexpected content type: '%v'", content_type)
	} else if ip := net.ParseIP(string(body)); ip == nil {
		return nil, fmt.Errorf("Unexpected response: '%v'", string(body))
	} else {
		return ip, nil
	}
}

func RegisterExternalAddress(user, passwd string, addr net.IP, hostname string) error {
	if req, err := NewRequest("GET", GOOGLE_URI); err != nil {
		return err
	} else {
		values := req.URL.Query()
		values.Set("hostname", hostname)
		values.Set("ip", addr.String())
		req.URL.RawQuery = values.Encode()
		req.URL.User = url.UserPassword(user, passwd)
		if body, _, err := Do(req); err != nil {
			return err
		} else {
			fmt.Printf("%s\n", string(body))
			return nil
		}
	}
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {	
	// Timer to discover IP address occasionally and re-register on change
	interval, _ := app.AppFlags.GetDuration('interval')
	if err := app.Timer.NewInterval(interval,nil,true); err != nil {
		return err
	}

	// Wait for CTRL+C
	app.Logger.Info("Press CTRL+C to exit")
	app.WaitForSignal()

	// Signal done, signal success
	done <- gopi.DONE
	return nil
}
/*
	// Return success
	if user, _ := app.AppFlags.GetString("user"); user == "" {
		return fmt.Errorf("Missing -user flag")
	} else if passwd, _ := app.AppFlags.GetString("passwd"); user == "" {
		return fmt.Errorf("Missing -passwd flag")
	} else if hostname, _ := app.AppFlags.GetString("hostname"); user == "" {
		return fmt.Errorf("Missing -hostname flag")
	} else if ip, err := GetExternalAddress(); err != nil {
		return err
	} else if err := RegisterExternalAddress(user, passwd, ip, hostname); err != nil {
		return err
	} else {
		fmt.Println(ip)
	}
	return nil
}
*/
func main() {
	config := gopi.NewAppConfig("timer")
	// Google username/password combination
	config.AppFlags.FlagDuration("interval", 60*time.Minute, "IP address discovery interval")
	config.AppFlags.FlagString("user", "", "Google Domains username")
	config.AppFlags.FlagString("passwd", "", "Google Domains password")
	config.AppFlags.FlagString("hostname", "", "Google Domains hostname")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
