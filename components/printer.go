package components

import (
	"fmt"

	"github.com/trustmaster/goflow"
)

type printer struct {
	In <-chan ImageSrc
}

// NewPrinter ...
func NewPrinter() goflow.Component {
	return &printer{}
}

// Process
func (p *printer) Process() {
	for value := range p.In {
		fmt.Printf("Src: %s\n", value)
	}
}
