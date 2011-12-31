/**
 * User: r2p2
 * Date: 12/29/11
 * Time: 12:03 AM
 */
package main

import (
	"fmt"
	"bytes"
	"rand"
	"time"
	"strings"
	"os"
)

const MAXWORKERINDEX = 4

var maps = map[string]string{
	"glider": "           \n" +
		"   #       \n" +
		"    #      \n" +
		"  ###      \n" +
		"           \n" +
		"           \n" +
		"           \n" +
		"           \n",
}

type Field struct {
	iteration             uint64
	width, height         int32
	currentArea, nextArea []byte
}

func NewField(width, height int32) *Field {
	rand.Seed(time.Nanoseconds())

	cells := width * height
	return &Field{
		iteration:   0,
		width:       width,
		height:      height,
		currentArea: make([]byte, cells),
		nextArea:    make([]byte, cells),
	}
}

func NewFieldFromString(field string) *Field {
	var width, height int32
	lines := strings.Split(field, "\n")
	height = int32(len(lines))
	for _, line := range lines {
		if linewidth := int32(len(line)); linewidth > width {
			width = linewidth
		}
	}

	newField := NewField(width, height)
	for y, _ := range lines {
		for x, _ := range lines[y] {
			if lines[y][x] != ' ' {
				newField.Set(int32(x), int32(y), 1)
			}
		}
	}

	return newField
}

func NewFieldFromMap(name string) (*Field, os.Error) {
	data, error := maps[name]
	if error != true {
		return nil, os.NewError("Unknown map: " + name)
	}
	return NewFieldFromString(data), nil
}

func (f *Field) Initialize(fillRate float32) {
	f.iteration = 0
	for index := range f.currentArea {
		if rand.Float32() <= fillRate {
			f.currentArea[index] = 1
		}
	}
}

func (f *Field) Clear() {
	f.iteration = 0
	for index := range f.currentArea {
		f.currentArea[index] = 0
	}
}

func (f *Field) Set(x, y int32, value byte) {
	f.currentArea[f.toArea(x, y)] = value
}

func (f *Field) Step() {
	f.iteration++

	resChan := make(chan byte, MAXWORKERINDEX)
	for workerIndex := int32(1); workerIndex <= MAXWORKERINDEX; workerIndex++ {
		go f.worker(workerIndex, resChan)
	}

	for workerIndex := int32(1); workerIndex <= MAXWORKERINDEX; workerIndex++ {
		<-resChan
	}
	f.swapFields()
}

func (f *Field) CellCount() int32 {
	return f.width * f.height
}

func (f *Field) Iteration() uint64 {
	return f.iteration
}

func (f *Field) String() string {
	sbuffer := bytes.NewBufferString("")
	for y := int32(0); y < f.height; y++ {
		for x := int32(0); x < f.width; x++ {
			index := f.toArea(x, y)
			if f.currentArea[index] == 1 {
				fmt.Fprint(sbuffer, "#")
			} else {
				fmt.Fprint(sbuffer, " ")
			}
		}
		fmt.Fprint(sbuffer, "\n")
	}
	return string(sbuffer.Bytes())
}

func (f *Field) StringNeighborMap() string {
	sbuffer := bytes.NewBufferString("")
	for y := int32(0); y < f.height; y++ {
		for x := int32(0); x < f.width; x++ {
			fmt.Fprint(sbuffer, f.countNeighbors(x, y))
		}
		fmt.Fprint(sbuffer, "\n")
	}
	return string(sbuffer.Bytes())
}

func (f *Field) Height() int32 {
	return f.height
}

func (f *Field) Width() int32 {
	return f.width
}

func (f *Field) worker(workerIndex int32, resChan chan byte) {
	var neighbors byte
	var x, y int32
	for cellIndex := workerIndex - 1; cellIndex < f.CellCount(); cellIndex += workerIndex {
		x, y = f.toReal(cellIndex)
		neighbors = f.countNeighbors(x, y)

		if neighbors == 3 {
			f.nextArea[cellIndex] = 1
		} else if neighbors == 2 {
			f.nextArea[cellIndex] = f.currentArea[cellIndex]
		} else {
			f.nextArea[cellIndex] = 0
		}
	}
	resChan <- 0
}

func (f *Field) toReal(index int32) (x, y int32) {
	y = index / f.width
	x = index - y*f.width
	return
}

func (f *Field) toArea(x, y int32) int32 {
	if x < 0 {
		x = f.width - 1
	} else if x >= f.width {
		x = 0
	}
	if y < 0 {
		y = f.height - 1
	} else if y >= f.height {
		y = 0
	}

	return y*f.width + x
}

func (f *Field) countNeighbors(x, y int32) (neighbors byte) {
	neighbors += f.currentArea[f.toArea(x-1, y-1)]
	neighbors += f.currentArea[f.toArea(x, y-1)]
	neighbors += f.currentArea[f.toArea(x+1, y-1)]
	neighbors += f.currentArea[f.toArea(x-1, y)]
	neighbors += f.currentArea[f.toArea(x+1, y)]
	neighbors += f.currentArea[f.toArea(x-1, y+1)]
	neighbors += f.currentArea[f.toArea(x, y+1)]
	neighbors += f.currentArea[f.toArea(x+1, y+1)]
	return
}

func (f *Field) swapFields() {
	backupPointer := f.nextArea
	f.nextArea = f.currentArea
	f.currentArea = backupPointer
}
