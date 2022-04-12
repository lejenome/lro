package process

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Queue interface {
	Len() (int, error)
	Add(job *Job) error
	SafeAdd(job *Job) error
	AddJobId(id uuid.UUID) error
	Get() (*Job, error)
	GetJobId() (uuid.UUID, error)
	// LockJob(job *Job) error
	// UnlockJob(job *Job) error
	// Delete(job *Job) error
	Purge() error
	Close() error
}

type inMemoryQueue struct {
	sync.RWMutex
	jobs   JobCache
	jobIds []uuid.UUID
	closed bool
}

var _ Queue = (*inMemoryQueue)(nil)

func DefaultQueue(jc JobCache) Queue {
	return &inMemoryQueue{jobs: jc}
}

func (q *inMemoryQueue) Len() (int, error) {
	return len(q.jobIds), nil
}
func (q *inMemoryQueue) Add(job *Job) error {
	if err := q.jobs.Add(job); err != nil {
		return err
	}
	return q.AddJobId(job.ID)
}
func (q *inMemoryQueue) SafeAdd(job *Job) error {
	q.jobs.Add(job)
	return q.AddJobId(job.ID)
}
func (q *inMemoryQueue) AddJobId(jobid uuid.UUID) error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	q.jobIds = append(q.jobIds, jobid)
	return nil
}
func (q *inMemoryQueue) GetJobId() (uuid.UUID, error) {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return uuid.Nil, errors.New("Queue closed")
	}
	if len(q.jobIds) > 0 {
		job := q.jobIds[0]
		q.jobIds = q.jobIds[1:]
		return job, nil
	}
	return uuid.Nil, errors.New("Queue is empty")
}
func (q *inMemoryQueue) Get() (*Job, error) {
	id, err := q.GetJobId()
	if err != nil {
		return nil, err
	}
	return q.jobs.Get(id)
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
