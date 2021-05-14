package components

import (
	"sync"

	"github.com/fatih/color"
	"github.com/trustmaster/goflow"
)

type errorHandler struct {
	In map[string]<-chan error
}

// NewErrorHandler creates a, well, error handler
func NewErrorHandler() goflow.Component {
	return &errorHandler{}
}

// Process ...
func (e *errorHandler) Process() {
	var wg sync.WaitGroup

	for component, ch := range e.In {
		wg.Add(1)
		go func(component string, ch <-chan error) {
			defer wg.Done()
			for err := range ch {
				color.Red("Error [%s]: %s", component, err)
			}
		}(component, ch)
	}
	wg.Wait()
}
