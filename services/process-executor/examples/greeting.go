package examples

import (
	"errors"
	"time"

	"github.com/lejenome/lro/services/process-executor/lib/process"
	"github.com/mitchellh/mapstructure"

	"github.com/fatih/structs"
)

func greeting(name string) string {
	time.Sleep(1)
	return "Hello " + name
}

// GreetingProcessV1 is a process example that accepts a map containing a name
// and returns a map containing a greeting
func GreetingProcessV1(in map[string]interface{}) (interface{}, error) {
	_name, ok := in["name"]
	if !ok {
		return nil, errors.New("name not found")
	}
	name, ok := _name.(string)
	if !ok {
		return nil, errors.New("name should be a string")
	}
	out := make(map[string]interface{})
	out["greeting"] = greeting(name)
	return out, nil
}

var _ process.HandlerFunc = GreetingProcessV1

type greetingInput struct {
	Name string
}
type greetingOutput struct {
	Greeting string
}

// GreetingProcessV2 is a process example that uses internally structs for
// input and output instead of just plain maps used on GreetingProcessV1.
func GreetingProcessV2(in map[string]interface{}) (interface{}, error) {
	var data greetingInput
	mapstructure.Decode(in, &data)
	var out greetingOutput
	out.Greeting = greeting(data.Name)
	return structs.Map(&out), nil
}

var _ process.HandlerFunc = GreetingProcessV2
