package main

import (
	executor "github.com/lejenome/lro/services/process-executor"
	"github.com/lejenome/lro/services/process-executor/examples"
	"github.com/lejenome/lro/services/process-executor/lib/process"
)

func main() {
	service := executor.Service{}
	service.Init()
	service.Run()
	service.RegisterProcesses(process.Process{
		Name:    "greeting:v1",
		Handler: process.HandlerFunc(examples.GreetingProcessV1),
	})
	service.RegisterProcesses(process.Process{
		Name:    "greeting:v2",
		Handler: process.HandlerFunc(examples.GreetingProcessV2),
	})
	service.Run()
	/*
		queuePub := queues.NatsPublisher(s.config.Nats.URL, "lro", s.cache)
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
	*/
}
