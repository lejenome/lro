package process

import (
	"context"
	"time"
)

type ProcessOptions struct {
	Ctx   context.Context
	Delay time.Duration
}
type Process struct {
	Name    string
	Handler Handler
}
