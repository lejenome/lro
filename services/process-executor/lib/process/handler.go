package process

// Handler is an interface for processing an input map and returning an output
// or an error
type Handler interface {
	Handle(map[string]interface{}) (interface{}, error)
}

// HandlerFunc is a func type that process an input map and returns an output
// or an error
type HandlerFunc func(map[string]interface{}) (interface{}, error)

// Handle process an input and returns an output or an error
func (fn HandlerFunc) Handle(in map[string]interface{}) (interface{}, error) {
	return fn(in)
}

var _ Handler = (*HandlerFunc)(nil)

// NewHandler wrap any function of signature `func(struct) (struct, error)`
// and returns a Handler object.
// TODO: implement NewHandler
// TODO: implement map[string]interface{} to struct marshal and checks
func NewHandler(fn interface{}) Handler {
	panic("Not implemented")
}
