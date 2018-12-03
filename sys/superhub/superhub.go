package superhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Values struct {
	SNMPBase string
	Keys     map[string]string
}

type Superhub struct {
	Addr string
}

type superhub struct {
	log    gopi.Logger
	base   *url.URL
	client *http.Client
}

type response map[string]string

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	USER_AGENT = "github.com/djthorpe/ddregister"
)

const (
	SUPERHUB_DOWNSTREAM = iota
	SUPERHUB_UPSTREAM
	SUPERHUB_UPSTREAM_EXT
	SUPERHUB_UPSTREAM_STATUS
	SUPERHUB_SIGNAL_QUALITY
	SUPERHUB_QOS
	SUPERHUB_QOS_FLOWS
)

var (
	downstream = Values{
		SNMPBase: "1.3.6.1.2.1.10.127.1.1.1",
		Keys: map[string]string{
			"1.1": "chanid",
			"1.2": "freq",
			"1.3": "width",
			"1.4": "modulation",
			"1.5": "interleave",
			"1.6": "power",
			"1.7": "annex",
			"1.8": "storage",
		},
	}
	upstream = Values{
		SNMPBase: "1.3.6.1.2.1.10.127.1.1.2",
		Keys: map[string]string{
			"1.1":  "chanid",
			"1.2":  "freq",
			"1.3":  "width",
			"1.4":  "modulation",
			"1.5":  "slotsize",
			"1.6":  "timingofs",
			"1.7":  "backoffstart",
			"1.8":  "backoffend",
			"1.9":  "txbackoffstart",
			"1.10": "txbackoffend",
			"1.11": "scdmaactivecodes",
			"1.12": "scdmacodesperslot",
			"1.13": "scdmaframesize",
			"1.14": "scdmahoppingspeed",
			"1.15": "type",
			"1.16": "clonefrom",
			"1.17": "update",
			"1.18": "status",
			"1.19": "preeqenable",
		},
	}
	upstreamext = Values{
		SNMPBase: "1.3.6.1.4.1.4115.1.3.4.1.9.2",
		Keys: map[string]string{
			"1.1": "chanid",
			"1.2": "symrate",
			"1.3": "modulation",
		},
	}
	upstreamstatus = Values{
		SNMPBase: "1.3.6.1.4.1.4491.2.1.20.1.2",
		Keys: map[string]string{
			"1.1": "power",
			"1.2": "t3timeouts",
			"1.3": "t4timeouts",
			"1.4": "rangingaborteds",
			"1.5": "modulation",
			"1.6": "eqdata",
			"1.7": "t3exceededs",
			"1.8": "ismuted",
			"1.9": "ranging",
		},
	}
	signalqualityext = Values{
		SNMPBase: "1.3.6.1.4.1.4491.2.1.20.1.24",
		Keys: map[string]string{
			"1.1": "rxmer",
			"1.2": "rxmersamples",
		},
	}
	qos = Values{
		SNMPBase: "1.3.6.1.4.1.4491.2.1.21.1.2.1.6",
		Keys: map[string]string{
			"2.1": "maxrate",
			"2.2": "",
			"2.3": "",
		},
	}
	qosflows = Values{
		SNMPBase: "1.3.6.1.4.1.4491.2.1.21.1.3.1",
		Keys: map[string]string{
			"6.2":  "sfsid",
			"7.2":  "direction",
			"8.2":  "primary",
			"9.2":  "flowparam",
			"10.2": "chansetid",
			"11.2": "flowattrsuccess",
			"12.2": "sfdsid",
			"13.2": "",
			"14.2": "",
			"15.2": "",
			"16.2": "",
			"17.2": "",
		},
	}
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Superhub) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.superhub.Open>{ Addr='%v' }", config.Addr)

	this := new(superhub)
	this.log = log
	this.client = &http.Client{}

	if url, err := url.Parse("http://" + config.Addr + "/walk"); err != nil {
		return nil, err
	} else {
		this.base = url
	}

	if err := this.get(&downstream); err != nil {
		return nil, err
	}

	return this, nil
}

func (this *superhub) Close() error {
	this.log.Debug("<sys.superhub.Close>{ base=%v }", this.base)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *superhub) String() string {
	return fmt.Sprintf("<sys.superhub>{ base=%v }", this.base)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func newrequest(method, url string) (*http.Request, error) {
	if req, err := http.NewRequest(method, url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Add("User-Agent", USER_AGENT)
		return req, nil
	}
}

func (this *superhub) do(req *http.Request) (*response, error) {
	var data response
	if resp, err := this.client.Do(req); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Error: %v: %v", resp.StatusCode, resp.Status)
		} else if err := decoder.Decode(&data); err != nil {
			return nil, err
		} else {
			return &data, nil
		}
	}
}

func (this *superhub) get(keys *Values) error {
	this.log.Debug2("<sys.superhub.get>{}")

	if req, err := newrequest("GET", this.base.String()); err != nil {
		return err
	} else {
		values := req.URL.Query()
		values.Set("oids", keys.SNMPBase)
		req.URL.RawQuery = values.Encode()
		if data, err := this.do(req); err != nil {
			return err
		} else {
			this.log.Debug2("<sys.superhub.get>{ response='%v' }", data)
			return nil
		}
	}
}

/*
func (this *superhub) get(values *Values) error {
	url = this.url
}
    r = requests.get(URLBASE + 'walk?oids=' + keymap[page]['snmpbase'])
    jdata = r.json()
    data = {}
    for key in jdata:
        if jdata[key] == 'Finish':
            continue
        keyext = key[len(keymap[page]['snmpbase']) + 1:key.rfind('.')]
        index = key[key.rfind('.') + 1:]
        if index not in data:
            data[index] = {}
        if keyext not in keymap[page]['keys']:
            print("Unknown:", keymap[page]['snmpbase'], keyext, index, ':' +
                  jdata[key] + ':')
        elif keymap[page]['keys'] is None:
            1  # Ignore
        else:
            data[index][keymap[page]['keys'][keyext]] = jdata[key]

    if flatten is None:
        return data

    newdata = []
    for channel in sorted(data.values(), key=lambda x: int(x[flatten])):
        newdata.append(channel)

    return newdata

*/
