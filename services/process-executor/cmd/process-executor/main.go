package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lejenome/lro/pkg/config"
	executor "github.com/lejenome/lro/services/process-executor"
	"github.com/lejenome/lro/services/process-executor/examples"
	"github.com/lejenome/lro/services/process-executor/lib/process"
	"github.com/lejenome/lro/services/process-executor/lib/process/db"
	"github.com/lejenome/lro/services/process-executor/lib/process/queues"
	"github.com/lejenome/lro/services/process-executor/lib/process/redis"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	var conf executor.ProcessExecutorConfig
	if err := config.Load(&conf); err != nil {
		panic(fmt.Errorf("Config error: %w", err))
	}
	cache := redis.RedisJobCache(conf.Redis.URL, conf.Redis.Username, conf.Redis.Password, conf.Redis.DB)
	queue := queues.NatsSubscriber(conf.Nats.URL, "lro", cache)
	queuePub := queues.NatsPublisher(conf.Nats.URL, "lro", cache)
	store := db.DBJobStore(conf.Database.URL)
	runner := process.DefaultRunner(queue, store)
	runner.Register(process.Process{
		Name:    "greeting:v1",
		Handler: process.HandlerFunc(examples.GreetingProcessV1),
	})
	runner.Register(process.Process{
		Name:    "greeting:v2",
		Handler: process.HandlerFunc(examples.GreetingProcessV2),
	})
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
		fmt.Printf("Running Handler with input : (%v)\n", test.In)
		out, err := runner.Handle(test.In)
		fmt.Printf("= Expected: %v, %v\n", test.Out, test.Err)
		fmt.Printf("= Got: %v, %v\n", out, err)
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
	runner.Run()
}
