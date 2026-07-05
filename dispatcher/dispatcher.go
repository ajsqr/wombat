package dispatcher

import (
	"github.com/ajsqr/wombat/frame"
)

//go:generate mockgen -source=dispatcher.go -destination=mocks/mock_dispatcher.go -package=mock_dispatcher
type Dispatcher interface {
	// Dispatch allows sending frames into the tunnel by the endpoint clients
	Dispatch(f *frame.Frame)

	// Stream begins streaming frames to and from the tunnel
	Stream() error
}
