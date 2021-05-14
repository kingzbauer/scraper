package components

import (
	"sync"

	"github.com/fatih/color"
	"github.com/trustmaster/goflow"
)

type debugger struct {
	In map[string]<-chan string
}

// NewDebugger ...
func NewDebugger() goflow.Component {
	return &debugger{}
}

func (d *debugger) Process() {
	var wg sync.WaitGroup

	for component, ch := range d.In {
		wg.Add(1)
		go func(component string, ch <-chan string) {
			defer wg.Done()
			for str := range ch {
				color.Blue("DEBUG [%s]: %s", component, str)
			}
		}(component, ch)
	}
	wg.Wait()
}
