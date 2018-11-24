package dns

import (
	// Frameworks
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type DNS struct {
	User, Passwd string
}

type dns struct {
	log          gopi.Logger
	user, passwd string
	client       *http.Client
	addr         net.IP
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	USER_AGENT = "github.com/djthorpe/ddregister"
	APIFY_URI  = "https://api.ipify.org"
	GOOGLE_URI = "https://domains.google.com/nic/update"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config DNS) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.dns.Open>{ user='%v' passwd='%v' }", config.User, config.Passwd)

	this := new(dns)
	this.log = log
	this.user = config.User
	this.passwd = config.Passwd
	this.client = &http.Client{}

	return this, nil
}

func (this *dns) Close() error {
	this.log.Debug("<sys.dns.Close>{ user='%v' passwd='%v' }", this.user, this.passwd)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *dns) String() string {
	return fmt.Sprintf("<sys.dns>{ addr=%v }", this.addr)
}

////////////////////////////////////////////////////////////////////////////////
// REGISTER

func (this *dns) Register(host string) error {
	this.log.Debug2("<sys.dns.Register>{ host='%v' }", host)

	// Obtain the current IP address
	if addr, err := this.GetExternalAddress(); err != nil {
		return err
	} else if this.addr == nil || addr.Equal(this.addr) == false {
		// Register DNS
		if err := this.RegisterExternalAddress(addr, host); err != nil {
			return err
		}
		this.addr = addr
	}

	// Success
	return nil
}

func (this *dns) RegisterExternalAddress(addr net.IP, host string) error {
	this.log.Debug2("<sys.dns.RegisterExternalAddress>{ addr=%v host='%v' }", addr, host)

	if req, err := NewRequest("GET", GOOGLE_URI); err != nil {
		return err
	} else {
		values := req.URL.Query()
		values.Set("hostname", host)
		values.Set("ip", addr.String())
		req.URL.RawQuery = values.Encode()
		req.URL.User = url.UserPassword(this.user, this.passwd)
		if body, _, err := this.Do(req); err != nil {
			return err
		} else {
			this.log.Debug2("<sys.dns.RegisterExternalAddress>{ response='%v' }", string(body))
			return nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func NewRequest(method, url string) (*http.Request, error) {
	if req, err := http.NewRequest(method, url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Add("User-Agent", USER_AGENT)
		return req, nil
	}
}

func (this *dns) Do(req *http.Request) ([]byte, string, error) {
	if resp, err := this.client.Do(req); err != nil {
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

func (this *dns) GetExternalAddress() (net.IP, error) {
	// Obtain the current external IP address
	if req, err := NewRequest("GET", APIFY_URI); err != nil {
		return nil, err
	} else if body, content_type, err := this.Do(req); err != nil {
		return nil, err
	} else if content_type != "text/plain" {
		return nil, fmt.Errorf("Unexpected content type: '%v'", content_type)
	} else if ip := net.ParseIP(string(body)); ip == nil {
		return nil, fmt.Errorf("Unexpected response: '%v'", string(body))
	} else {
		return ip, nil
	}
}
