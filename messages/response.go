package messages

// A ResponderError is used to represent an error in a module response.
type ResponderError struct {
	ErrorCode int32
	Message   string
}

// A Responder is the interface that are used to record responses from modules
// to other modules or the proxy core.
type Responder interface {
	Type() string

	// Write takes a ResponseWriter to be sent to an upstream module, or to proxy core.
	Write(ResponseWriter)

	// Error indicates that the response failed, and takes in an int32 error code
	// and a string for an error message.
	Error(int32, string)
}

// Struct that records the responses from modules to be handled by proxy core.
// implements a Responder
type ModuleResponse struct {
	CommandError *ResponderError
	Writer       ResponseWriter
}

func (r *ModuleResponse) Type() string {
	return "response"
}

func (r *ModuleResponse) Write(writer ResponseWriter) {
	r.Writer = writer
}

func (r *ModuleResponse) Error(code int32, message string) {
	r.CommandError = &ResponderError{code, message}
}
