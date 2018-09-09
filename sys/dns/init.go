package dns

import (
	"fmt"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register Timer
	gopi.RegisterModule(gopi.Module{
		Name: "sys/dns",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("dns.user", "", "DNS registry username")
			config.AppFlags.FlagString("dns.passwd", "", "DNS registry password")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			if user, exists := app.AppFlags.GetString("dns.user"); exists == false {
				return nil, fmt.Errorf("Missing -dns.user flag")
			} else if passwd, exists := app.AppFlags.GetString("dns.passwd"); exists == false {
				return nil, fmt.Errorf("Missing -dns.passwd flag")
			} else {
				return gopi.Open(DNS{User: user, Passwd: passwd}, app.Logger)
			}
		},
	})
}
