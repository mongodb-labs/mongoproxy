// Package modules contains interfaces and functions dealing with the creation
// of modules and the module pipeline within proxy core.
package modules

import (
	"github.com/mongodbinc-interns/mongoproxy/messages"
)

type Module interface {

	// Process is the function executed when a message is called in the pipeline.
	// It takes in a Requester from an upstream module (or proxy core), a
	// Responder that it writes a response to, and a PipelineFunc that should
	// be called to execute the next module in the pipeline.
	Process(messages.Requester, messages.Responder, PipelineFunc)
}
