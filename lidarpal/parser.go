package lidarpal

import (
	"strconv"
	"strings"
	"sync"

	"github.com/hongping1224/lidario"
)

// Reader read Pointcloud
type Parser struct {
	input chan string
	wg    *sync.WaitGroup
}

// NewReader Create a new Reader
func NewParser(input chan string, wg *sync.WaitGroup) *Parser {
	return &Parser{input: input, wg: wg}
}

// Read point into channel
func (parser *Parser) Parse(input chan<- lidario.LasPointer) {
	for {
		s, open := <-parser.input
		if !open {
			//fmt.Println("Writer Closing")
			break
		}
		data := strings.Split(s, " ")
		if len(data) < 4 {
			continue
		}
		x, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			continue
		}
		y, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			continue
		}
		z, err := strconv.ParseFloat(data[2], 64)
		if err != nil {
			continue
		}
		intensity, err := strconv.ParseFloat(data[3], 64)
		if err != nil {
			continue
		}
		p := lidario.PointRecord0{X: x, Y: y, Z: z, Intensity: uint16(intensity)}
		input <- &p
	}
	parser.wg.Done()
}

// Serve read concerrently
func (parser *Parser) Serve(input chan<- lidario.LasPointer) {
	go parser.Parse(input)
}
