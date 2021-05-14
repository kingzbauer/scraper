package components

import (
	"github.com/trustmaster/goflow"
)

// imageParser will takes images from a page, make sure they are well formatted and
// output them as individual image sources
type imageParser struct {
	In  <-chan PageImages
	Out chan<- ImageSrc
}

func (i *imageParser) Process() {
	for page := range i.In {
		if sources, err := page.ToImageSlice(); err == nil {
			for _, src := range sources {
				if len(src) == 0 {
					continue
				}
				i.Out <- src
			}
		}
	}
}

// NewImageParser will takes images from a page, make sure they are well formatted and
// output them as individual image sources
func NewImageParser() goflow.Component {
	return &imageParser{}
}
