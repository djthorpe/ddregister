package superhub

import (
	// Frameworks
	"fmt"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register Timer
	gopi.RegisterModule(gopi.Module{
		Name: "sys/superhub",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("superhub.addr", "", "Superhub address")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			if addr, exists := app.AppFlags.GetString("superhub.addr"); exists == false {
				return nil, fmt.Errorf("Missing -superhub.addr argument")
			} else {
				return gopi.Open(Superhub{
					Addr: addr,
				}, app.Logger)
			}
		},
	})
}
