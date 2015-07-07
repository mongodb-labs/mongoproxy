package server

import (
	"github.com/mongodbinc-interns/mongoproxy/messages"
)

// PipelineFunc is the function type for the built pipeline, and is called
// to begin the pipeline.
type PipelineFunc func(messages.Requester, messages.Responder)

// A ChainFunc is a closure that wraps a module so that they can accept
// other modules as inputs and outputs for module chaining.
type ChainFunc func(PipelineFunc) PipelineFunc

// A ModuleChain consists of a chain of wrapped modules that can be built
// into a single pipeline function.
type ModuleChain struct {
	chain []ChainFunc
}

// AddModules adds the module mods to the end of a given module chain.
func (m *ModuleChain) AddModules(mods ...Module) *ModuleChain {
	for i := 0; i < len(mods); i++ {
		m.chain = append(m.chain, wrapModule(mods[i]))
	}

	return m
}

// wrapModule returns a closure ChainFunc that wraps over the module m, which
// can input and output PipelineFuncs to help with chaining.
func wrapModule(m Module) ChainFunc {

	return ChainFunc(func(next PipelineFunc) PipelineFunc {
		return PipelineFunc(func(r messages.Requester, w messages.Responder) {

			// if there is no next module in the pipeline, the pipeline terminates
			if next == nil {
				next = PipelineFunc(func(r messages.Requester, w messages.Responder) {
					return
				})
			}
			m.Process(r, w, next)
		})
	})
}

// CreateChain initializes and returns an empty module chain that can be used
// to build a pipeline
func CreateChain() *ModuleChain {
	return &ModuleChain{
		chain: make([]ChainFunc, 0),
	}
}

// BuildPipeline takes a module chain and creates a pipeline, returning
// a PipelineFunc that starts the pipeline when called.
// The proxy core manages the pipeline order by setting the PipelineFuncs of each
// module to the next module in the pipeline. The HandleFunc of the last module
// in the pipeline is set to nil to terminate the pipeline.
func BuildPipeline(m *ModuleChain) PipelineFunc {

	if len(m.chain) == 0 {
		return PipelineFunc(func(r messages.Requester, w messages.Responder) {
			return
		})
	}
	pipeline := m.chain[len(m.chain)-1](nil)
	for i := len(m.chain) - 2; i >= 0; i-- {
		pipeline = m.chain[i](pipeline)
	}

	return pipeline
}
