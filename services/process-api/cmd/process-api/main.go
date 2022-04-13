package main

import (
	"github.com/google/uuid"
	api "github.com/lejenome/lro/services/process-api"
	"github.com/lejenome/lro/services/process-executor/lib/process"
)

func main() {
	service := api.Service{}
	service.Init()
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
		job := process.Job{
			ID:          uuid.Nil,
			ProcessName: _pName.(string),
			Input:       _data.(map[string]interface{}),
		}
		service.ScheduleJob(&job)
	}
	service.Run()
}
