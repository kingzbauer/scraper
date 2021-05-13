package components

import (
	"sync"

	"github.com/trustmaster/goflow"
)

// Retriever provides the mechanics of retrieving a html page given a url
type Retriever func(string) (string, error)

type retriever struct {
	// In will receiver the url to fetch
	In <-chan string
	// Out will carry the HTML content of the url
	Out    chan<- string
	getter Retriever
}

// NewRetriever ...
func NewRetriever(getter Retriever) goflow.Component {
	return &retriever{getter: getter}
}

// Process is the entrypoint of the component
func (r *retriever) Process() {
	var wg sync.WaitGroup

	for url := range r.In {
		wg.Add(1)
		// For now the goroutines will be unbounded
		go func(url string) {
			defer wg.Done()
			// TODO: Add an error process
			if html, err := r.getter(url); err == nil {
				r.Out <- html
			}
		}(url)
	}

	wg.Wait()
}
