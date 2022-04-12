package process

import (
	"encoding"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

type JobStatus string

const (
	JOB_PENDING JobStatus = "PENDING"
	JOB_RUNNING JobStatus = "RUNNING"
	JOB_SUCCESS JobStatus = "SUCCESS"
	JOB_FAILURE JobStatus = "FAILURE"
	JOB_CANCEL  JobStatus = "CANCEL"
	JOB_TIMEOUT JobStatus = "TIMEOUT"
)

type Job struct {
	ID          uuid.UUID
	ProcessName string
	Input       map[string]interface{}
	Output      interface{}
	Errors      []error
	Status      JobStatus
	/*
		Meta        struct {
			CreatedAt time.Time
			UpdatedAt time.Time
			StartedAt time.Time
			EndedAt   time.Time
		}
		Options struct {
			Retry uint8
		}
	*/
}

func (j *Job) IsPending() bool {
	return j.Status == "" || j.Status == JOB_PENDING
}

type rawJob Job

func (j *Job) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal((*rawJob)(j))
}
func (j *Job) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, (*rawJob)(j))

}

var _ encoding.BinaryMarshaler = (*Job)(nil)
var _ encoding.BinaryUnmarshaler = (*Job)(nil)

/*
type JobI interface {
	Id() string
	ProcessName() string
	Input() map[string]interface{}
	SetOutput(interface{})
	Output() interface{}
	Errors() []error
}
*/

type JobCache interface {
	Add(job *Job) error
	Get(id uuid.UUID) (*Job, error)
}
type JobStore interface {
	Save(job *Job) error
	Update(job *Job) error
	Retrive(id uuid.UUID) (*Job, error)
}

type inMemoryJobCache struct {
	sync.RWMutex
	jobs   map[uuid.UUID]*Job
	closed bool
}

var _ JobCache = (*inMemoryJobCache)(nil)

func DefaultJobCache() JobCache {
	return &inMemoryJobCache{
		jobs: make(map[uuid.UUID]*Job),
	}
}

func (q *inMemoryJobCache) Len() (int, error) {
	return len(q.jobs), nil
}
func (q *inMemoryJobCache) Add(job *Job) error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	if _, ok := q.jobs[job.ID]; ok {
		return fmt.Errorf("Job '%s' already in cache", job.ID)
	}
	q.jobs[job.ID] = job
	return nil
}
func (q *inMemoryJobCache) Get(id uuid.UUID) (*Job, error) {
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
