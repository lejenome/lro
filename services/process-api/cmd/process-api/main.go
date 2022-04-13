package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lejenome/lro/services/process-api/config"
	"github.com/lejenome/lro/services/process-executor/lib/process"
	"github.com/lejenome/lro/services/process-executor/lib/process/db"
	"github.com/lejenome/lro/services/process-executor/lib/process/queues"
	"github.com/lejenome/lro/services/process-executor/lib/process/redis"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	config, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("Config error: %w", err))
	}
	cache := redis.RedisJobCache(config.Redis.URL, config.Redis.Username, config.Redis.Password, config.Redis.DB)
	queuePub := queues.NatsPublisher(config.Nats.URL, "lro", cache)
	store := db.DBJobStore(config.Database.URL)
	tests := []struct {
		In  map[string]interface{}
		Out interface{}
		Err error
	}{
		{
			map[string]interface{}{
				"processName": "greeting:v1",
				"data": map[string]interface{}{
					"name": "World",
				},
			},

			map[string]interface{}{
				"greeting": "Hello World",
			},
			nil,
		},
		{
			map[string]interface{}{
				"processName": "greeting:v2",
				"data": map[string]interface{}{
					"name": "World",
				},
			},

			map[string]interface{}{
				"greeting": "Hello World",
			},
			nil,
		},
	}
	for _, test := range tests {
		_pName := test.In["processName"]
		_data := test.In["data"]
		job := &process.Job{
			ID:          uuid.Nil,
			ProcessName: _pName.(string),
			Input:       _data.(map[string]interface{}),
		}
		if err := store.Save(job); err != nil {
			fmt.Printf("%s\n", err)
		}
		if err := queuePub.SafeAdd(job); err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
