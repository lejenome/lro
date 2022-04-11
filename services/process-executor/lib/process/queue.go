package process

import (
	"errors"
	"sync"
)

type Queue interface {
	Len() (int, error)
	Add(job *Job) error
	Get() (*Job, error)
	// LockJob(job *Job) error
	// UnlockJob(job *Job) error
	// Delete(job *Job) error
	Purge() error
	Close() error
}

type inMemoryQueue struct {
	sync.RWMutex
	jobs   []*Job
	closed bool
}

var _ Queue = (*inMemoryQueue)(nil)

func DefaultQueue() Queue {
	return &inMemoryQueue{}
}

func (q *inMemoryQueue) Len() (int, error) {
	return len(q.jobs), nil
}
func (q *inMemoryQueue) Add(job *Job) error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	q.jobs = append(q.jobs, job)
	return nil
}
func (q *inMemoryQueue) Get() (*Job, error) {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return nil, errors.New("Queue closed")
	}
	if len(q.jobs) > 0 {
		job := q.jobs[0]
		q.jobs = q.jobs[1:]
		return job, nil
	}
	return nil, errors.New("Queue is empty")
}

// func (q *inMemoryQueue) Lock(job *Job) error {
// 	if q.closed {
// 		return errors.New("Queue closed")
// 	}
// 	panic("Not implemented")
// }
// func (q *inMemoryQueue) Unlock(job *Job) error {
// 	if q.closed {
// 		return errors.New("Queue closed")
// 	}
// 	panic("Not implemented")
// }
func (q *inMemoryQueue) Purge() error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	q.jobs = nil
	return nil
}
func (q *inMemoryQueue) Close() error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	q.closed = true
	return nil
}
