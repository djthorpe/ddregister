package main

import (
	"fmt"
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/ddregister/sys/dns"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/timer"
)

////////////////////////////////////////////////////////////////////////////////
// MODULE INTERFACE

type DNS interface {
	// Obtain external IP address and register under hostname
	Register(host string) error
}

////////////////////////////////////////////////////////////////////////////////
// EVENT HANDLER

func HandleEvent(app *gopi.AppInstance, evt gopi.TimerEvent) error {
	dns := app.ModuleInstance("sys/dns").(DNS)
	host, _ := app.AppFlags.GetString("host")
	if err := dns.Register(host); err != nil {
		return err
	}

	// Return success
	return nil
}

func EventHandler(app *gopi.AppInstance, done <-chan struct{}) error {
	// register for events from the timer
	events := app.Timer.Subscribe()

	// Wait for maturing timers or done signal
FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case evt := <-events:
			if err := HandleEvent(app, evt.(gopi.TimerEvent)); err != nil {
				return err
			}
		}
	}

	// return success
	app.Timer.Unsubscribe(events)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	// Timer to discover IP address occasionally and re-register on change
	interval, _ := app.AppFlags.GetDuration("interval")
	if err := app.Timer.NewInterval(interval, nil, true); err != nil {
		return err
	}

	// Check for hostname flag
	if _, exists := app.AppFlags.GetString("host"); exists == false {
		done <- gopi.DONE
		return fmt.Errorf("Missing -host flag")
	}

	// Wait for CTRL+C
	app.Logger.Info("Press CTRL+C to exit")
	app.WaitForSignal()

	// Signal done, signal success
	done <- gopi.DONE
	return nil
}

func main() {
	config := gopi.NewAppConfig("timer", "sys/dns")

	// Google username/password combination
	config.AppFlags.FlagDuration("interval", 60*time.Minute, "IP address discovery interval")
	config.AppFlags.FlagString("host", "", "Hostname to register")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main, EventHandler))
}
