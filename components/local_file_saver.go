package components

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/trustmaster/goflow"
)

type localFileSaver struct {
	baseDir string
	In      <-chan ImageData
	Err     chan<- error
	Debug   chan<- string
}

// NewLocalFileSaver ...
func NewLocalFileSaver(baseDir string) goflow.Component {
	return &localFileSaver{baseDir: baseDir}
}

func (l *localFileSaver) Process() {
	var wg sync.WaitGroup
	for data := range l.In {
		wg.Add(1)
		go func(data ImageData) {
			defer wg.Done()
			filename := path.Base(string(data.Src))
			filepath := filepath.Join(l.baseDir, filename)
			l.Debug <- fmt.Sprintf("Writing to file: %s", filepath)
			if file, err := os.Create(filepath); err == nil {
				// TODO: probably use a buffered copy:w
				if n, err := io.Copy(file, data.File); err == nil {
					l.Debug <- fmt.Sprintf("Wrote to file: %s. %d bytes written", filepath, n)
				} else {
					l.Err <- err
				}
				// Close both files
				if closer, ok := data.File.(io.Closer); ok {
					if err := closer.Close(); err != nil {
						l.Err <- err
					}
				}
				if err := file.Close(); err != nil {
					l.Err <- err
				}
			} else {
				l.Err <- err
			}
		}(data)
	}

	wg.Wait()
}
