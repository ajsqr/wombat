package endpoint

import (
	"github.com/ajsqr/wombat/frame"
)

type Endpoint interface {
	Handle(f *frame.Frame) error
}
