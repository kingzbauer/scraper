package components

import (
	"github.com/trustmaster/goflow"
)

type spreader struct {
	In  <-chan []string
	Out chan<- string
}

// NewSpreader returns a spreader,
//
// A spreader receivers a slice of values and returns each entry in the slice as a
// single value to the output channel
func NewSpreader() goflow.Component {
	return &spreader{}
}

// Process ...
func (s *spreader) Process() {
	for slices := range s.In {
		for _, value := range slices {
			s.Out <- value
		}
	}
}
