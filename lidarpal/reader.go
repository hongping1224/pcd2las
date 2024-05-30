package lidarpal

import (
	"bufio"
	"fmt"
	"sync"

	"github.com/hongping1224/lidario"
)

// Reader read Pointcloud
type Reader struct {
	scanner *bufio.Scanner
	wg      *sync.WaitGroup
}

// NewReader Create a new Reader
func NewReader(scanner *bufio.Scanner, wg *sync.WaitGroup) *Reader {
	return &Reader{scanner: scanner, wg: wg}
}

// Read point into channel
func (read *Reader) Read(input chan<- lidario.LasPointer) {
	parsercount := 4
	parsers := make([]*Parser, parsercount)
	var wg sync.WaitGroup
	stringchan := make(chan string, 4*parsercount)

	for i := range parsercount {
		wg.Add(1)
		parsers[i] = NewParser(stringchan, &wg)
		parsers[i].Serve(input)
	}
	//Skip Header
	for i := 0; i < 13; i++ {
		read.scanner.Scan()
	}
	for read.scanner.Scan() {
		data := read.scanner.Text()
		stringchan <- data
	}
	wg.Wait()
	read.wg.Done()
	fmt.Println("Reader Done")
}

// Serve read concerrently
func (read *Reader) Serve(input chan<- lidario.LasPointer) {
	go read.Read(input)
}
