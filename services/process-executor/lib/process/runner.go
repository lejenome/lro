package process

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type registry map[string]Process

type Runner interface {
	Handler
	Register(Process)
	Unregister(processName string)
	Run()
	Close()
}
type basicRunner struct {
	registry  registry
	queue     Queue
	close     chan bool
	closeOnce sync.Once
}

func DefaultRunner(queue Queue) Runner {
	return &basicRunner{
		registry: make(registry),
		queue:    queue,
		close:    make(chan bool),
	}
}

func (r *basicRunner) Register(p Process) {
	r.registry[p.Name] = p
}

func (r *basicRunner) Unregister(name string) {
	delete(r.registry, name)
}

func (r *basicRunner) Handle(in map[string]interface{}) (interface{}, error) {
	_pName, ok := in["processName"]
	if !ok {
		return nil, errors.New("processName field is required")
	}
	pName := _pName.(string)
	_data, ok := in["data"]
	if !ok {
		return nil, errors.New("data field is required")
	}
	data, ok := _data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data field; expected type map[string]interface{}; got %T", data)
	}
	return r.RunProces(pName, data)
}

func (r *basicRunner) RunProces(name string, in map[string]interface{}) (interface{}, error) {
	process, ok := r.registry[name]
	if !ok {
		return nil, fmt.Errorf("process '%s' was not found", name)
	}
	return process.Handler.Handle(in)
}

func (r *basicRunner) RunJob(job *Job) {
	out, err := r.RunProces(job.ProcessName, job.Input)
	job.Output = out
	job.Errors = append(job.Errors, err)
}

func (r *basicRunner) Run() {
	for {
		select {
		case <-r.close:
			return
		default:
			// TODO user chan and select instead of sleep
			if n, err := r.queue.Len(); err == nil && n == 0 {
				return
			} else if job, err := r.queue.Get(); err == nil {
				go r.RunJob(job)
			} else {
				fmt.Printf("Error: %s\n", err)
				time.Sleep(1)
			}
		}
	}
}
func (r *basicRunner) Close() {
	r.closeOnce.Do(func() {
		close(r.close)
	})
}

var _ Runner = (*basicRunner)(nil)
