package process

import (
	"time"
)

type Job struct {
	Id          string
	ProcessName string
	Input       map[string]interface{}
	Output      interface{}
	Errors      []error
	Meta        struct {
		CreatedAt time.Time
		UpdatedAt time.Time
		StartedAt time.Time
		EndedAt   time.Time
	}
	Options struct {
		Retry uint8
	}
}

type JobI interface {
	Id() string
	ProcessName() string
	Input() map[string]interface{}
	SetOutput(interface{})
	Output() interface{}
	Errors() []error
}
