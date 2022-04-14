package process

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type registry map[string]Process

type Runner interface {
	Handler
	Register(Process)
	Unregister(processName string)
	Run(ctx context.Context)
	Close()
}
type basicRunner struct {
	registry  registry
	queue     Queue
	store     JobStore
	close     chan bool
	closeOnce sync.Once
}

func DefaultRunner(queue Queue, store JobStore) Runner {
	return &basicRunner{
		registry: make(registry),
		queue:    queue,
		store:    store,
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

func (r *basicRunner) RunJob(job *Job) error {
	fmt.Printf("RunJob[Id: %s, ProcessName: %s] Input (%v)\n", job.ID, job.ProcessName, job.Input)
	if !job.IsPending() {
		return fmt.Errorf("Can not run job '%s', expected PENDING status, got '%s' status\n", job.ID, job.Status)
	}
	out, err := r.RunProces(job.ProcessName, job.Input)
	job.Output = out
	if err != nil {
		job.Errors = append(job.Errors, err)
	}
	r.store.Update(job)
	fmt.Printf("RunJob[Id: %s, ProcessName: %s] Ouput (%v) Errors(%v)\n", job.ID, job.ProcessName, job.Output, job.Errors)
	return err
}
func (r *basicRunner) RunJobId(id uuid.UUID) error {
	job, err := r.store.Retrive(id)
	if err != nil {
		return err
	}
	return r.RunJob(job)
}

func (r *basicRunner) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-r.close:
			break loop
		default:
			// TODO user chan and select instead of sleep
			if n, err := r.queue.Len(); err == nil && n == 0 {
				break loop
			} else if id, err := r.queue.GetJobId(); err == nil {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := r.RunJobId(id)
					if err != nil {
						fmt.Printf("RubJob '%s', got error: %v\n", id, err)
					}
				}()
			} else {
				fmt.Printf("Error: %s\n", err)
				time.Sleep(1)
			}
		}
	}
	wg.Wait()
}
func (r *basicRunner) Close() {
	r.closeOnce.Do(func() {
		close(r.close)
	})
}

var _ Runner = (*basicRunner)(nil)
