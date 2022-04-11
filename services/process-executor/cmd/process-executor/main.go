package main

import (
	"fmt"

	"github.com/lejenome/lro/services/process-executor/examples"
	"github.com/lejenome/lro/services/process-executor/lib/process"
)

func main() {
	queue := process.DefaultQueue()
	runner := process.DefaultRunner(queue)
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
}
