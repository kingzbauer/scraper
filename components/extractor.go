package components

import (
	"sync"

	"github.com/anaskhan96/soup"
	"github.com/trustmaster/goflow"
)

type extractor struct {
	In  <-chan string
	Out chan<- []string
}

// NewExtractor ...
func NewExtractor() goflow.Component {
	return &extractor{}
}

// Process ...
func (e *extractor) Process() {
	var wg sync.WaitGroup
	for html := range e.In {
		wg.Add(1)
		go func(html string) {
			defer wg.Done()

			root := soup.HTMLParse(html)
			e.Out <- e.findImageSrc(root)
		}(html)
	}

	wg.Wait()
}

func (e *extractor) findImageSrc(root soup.Root) []string {
	elements := root.FindAll("img")
	imgSrcs := make([]string, len(elements))

	for i, elem := range elements {
		imgSrcs[i] = elem.Attrs()["src"]
	}

	return imgSrcs
}
