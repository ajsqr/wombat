package endpoint

import (
	"github.com/ajsqr/wombat/dispatcher"
	"github.com/ajsqr/wombat/frame"
)

// Endpoint defines the service an endpoint must perform.
// It defines how an endpoint device would behave with the dispatcher and frame it is managing.
type Endpoint interface {
	Handle(f *frame.Frame, dispatcher dispatcher.Dispatcher) error
}
