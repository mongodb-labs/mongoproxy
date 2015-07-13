// Package server contains interfaces and functions dealing with setting up proxy core,
// including code construct the module pipeline.
package server

import (
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"gopkg.in/mgo.v2/bson"
)

type Module interface {

	// Name returns the name to identify this module when registered.
	Name() string

	// Configure configures this module with the given configuration object.
	Configure(bson.M) error

	// Process is the function executed when a message is called in the pipeline.
	// It takes in a Requester from an upstream module (or proxy core), a
	// Responder that it writes a response to, and a PipelineFunc that should
	// be called to execute the next module in the pipeline.
	Process(messages.Requester, messages.Responder, PipelineFunc)

	// New creates a new instance of this module
	New() Module
}
