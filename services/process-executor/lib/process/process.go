package process

import (
	"context"
	"time"
)

type ProcessConfig struct {
	Ctx   context.Context
	Delay time.Duration
}
type Process struct {
	Name    string
	Handler Handler
}
