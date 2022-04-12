package process

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type inMemoryJobStore struct {
	sync.RWMutex
	jobs   map[uuid.UUID]*Job
	closed bool
}

var _ JobStore = (*inMemoryJobStore)(nil)

func DefaultJobStore() JobStore {
	return &inMemoryJobStore{
		jobs: make(map[uuid.UUID]*Job),
	}
}

func (q *inMemoryJobStore) Save(job *Job) error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	if _, ok := q.jobs[job.ID]; ok {
		return fmt.Errorf("Job '%s' already in store", job.ID)
	}
	q.jobs[job.ID] = job
	return nil
}
func (q *inMemoryJobStore) Update(job *Job) error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	if _, ok := q.jobs[job.ID]; !ok {
		return fmt.Errorf("Job '%s' not found in store", job.ID)
	}
	q.jobs[job.ID] = job
	return nil
}
func (q *inMemoryJobStore) Retrive(id uuid.UUID) (*Job, error) {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return nil, errors.New("Queue closed")
	}
	if job, ok := q.jobs[id]; ok {
		return job, nil
	}
	return nil, fmt.Errorf("Job '%s' not found", id)
}
