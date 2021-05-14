package components

import (
	"sync"

	"github.com/anaskhan96/soup"
	"github.com/trustmaster/goflow"
)

type extractor struct {
	In  <-chan Page
	Out chan<- PageImages
}

// NewExtractor ...
func NewExtractor() goflow.Component {
	return &extractor{}
}

// Process ...
func (e *extractor) Process() {
	var wg sync.WaitGroup
	for page := range e.In {
		wg.Add(1)
		go func(page Page) {
			defer wg.Done()

			root := soup.HTMLParse(*page.HTML)
			e.Out <- NewPageImages(e.findImageSrc(root), page.URL)
		}(page)
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
