package controller

// Pipeline manages controller handlers by providing a consistent way to
// propagate errors and handle web requests and responses
type Pipeline interface {
	// Next executes an action that may return an error that will be handled as a
	// http response
	Next(action func() error) Pipeline
}

type PipelineImpl struct {
	errorHandler func(err error)
	err          error
}

func (p PipelineImpl) Next(action func() error) Pipeline {
	if p.err != nil {
		return p
	}
	p.err = action()
	if p.err != nil {
		p.errorHandler(p.err)
	}
	return p
}

func NewPipelineImpl(errorHandler func(err error)) *PipelineImpl {
	return &PipelineImpl{errorHandler: errorHandler, err: nil}
}
