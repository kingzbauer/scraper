package components

import (
	"fmt"
	"sync"

	"github.com/trustmaster/goflow"
)

type imageRetriever struct {
	In    <-chan ImageSrc
	Out   chan<- ImageData
	Err   chan<- error
	Debug chan<- string
}

// NewImageRetriever ...
func NewImageRetriever() goflow.Component {
	return &imageRetriever{}
}

func (i *imageRetriever) Process() {
	var wg sync.WaitGroup
	for src := range i.In {
		wg.Add(1)
		go func(src ImageSrc) {
			defer wg.Done()

			data, err := src.GetData()
			if err != nil {
				i.Err <- err
			} else {
				i.Out <- data
				i.Debug <- fmt.Sprintf("Retrieved image: %s", data.Src)
			}
		}(src)
	}

	wg.Wait()
}
