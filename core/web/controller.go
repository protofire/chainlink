package web

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

// controller wraps a chainlink.Application and logger.Logger to embed in a Controller.
// lggr is derived from app.GetLogger() and may be extended (see namedController).
type controller struct {
	app  chainlink.Application
	lggr logger.Logger
}

func newController(app chainlink.Application) controller {
	return controller{app: app, lggr: app.GetLogger()}
}

// namedController returns a new controller with a named logger.
// If the logger is only called from the *_controller.go file, then naming it may be redundant and newController
// should be considered instead.
func namedController(app chainlink.Application, name string) controller {
	return controller{app: app, lggr: app.GetLogger().Named(name)}
}
