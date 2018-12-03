package main

import (
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/ddregister/sys/superhub"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/timer"
)

////////////////////////////////////////////////////////////////////////////////
// MAIN

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Get superhub
	superhub := app.ModuleInstance("sys/superhub")
	app.Logger.Info("superhub=%v", superhub)

	// Wait for CTRL+C
	app.Logger.Info("Press CTRL+C to exit")
	app.WaitForSignal()

	// Signal done, signal success
	done <- gopi.DONE
	return nil
}

func main() {
	config := gopi.NewAppConfig("timer", "sys/superhub")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
