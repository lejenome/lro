package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/lejenome/lro/services/process-executor/lib/process"
)

var ctx = context.Background()

type redisJobCache struct {
	sync.RWMutex
	rdb    *redis.Client
	closed bool
}

var _ process.JobCache = (*redisJobCache)(nil)

func RedisJobCache(url string, username string, password string, db int) process.JobCache {
	return &redisJobCache{
		rdb: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: password,
			Username: username,
			DB:       db,
		}),
	}
}

func (q *redisJobCache) Len() (int, error) {
	panic("Not supported")
}
func (q *redisJobCache) Add(job *process.Job) error {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return errors.New("Queue closed")
	}
	_, err := q.rdb.Get(ctx, "job:"+job.ID.String()).Result()
	if err != nil && err != redis.Nil {
		return err
	} else if err == nil {
		return fmt.Errorf("Job '%s' already added", job.ID)
	}
	err = q.rdb.Set(ctx, "job:"+job.ID.String(), job, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
func (q *redisJobCache) Get(id uuid.UUID) (*process.Job, error) {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return nil, errors.New("Queue closed")
	}
	data, err := q.rdb.Get(ctx, "job:"+id.String()).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("Job '%s' not found", id)
	} else if err != nil {
		return nil, err
	}
	job := &process.Job{}
	job.UnmarshalBinary([]byte(data))
	return job, nil
}
