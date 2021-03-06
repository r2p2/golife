/**
 * User: r2p2
 * Date: 12/29/11
 * Time: 12:03 AM
 */
package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
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
	ruleSurvive           byte
	ruleBirth             byte
}

func NewField(width, height int32) *Field {
	rand.Seed(time.Now().UnixNano())

	cells := width * height
	return &Field{
		iteration:   0,
		width:       width,
		height:      height,
		currentArea: make([]byte, cells),
		nextArea:    make([]byte, cells),
		ruleSurvive: 2,
		ruleBirth:   4,
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

func NewFieldFromMap(name string) (*Field, error) {
	data, error := maps[name]
	if error != true {
		return nil, errors.New("Unknown map: " + name)
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

func (f *Field) SetRule(rule string) error {
	var newRuleSurvive, newRuleBirth byte

	ruleArray := strings.Split(rule, "/")
	if len(ruleArray) != 2 {
		return errors.New("Invalid rule format! Exactly one / is needed.")
	}

	for _, char := range ruleArray[1] {
		if char < '1' || char > '9' {
			return errors.New("Invalid rule argument! Nubers needs to be between 1 and 9.")
		}

		newRuleBirth += byte(1 << uint(char-'1'))
	}

	for _, char := range ruleArray[0] {
		if char < '1' || char > '9' {
			return errors.New("Invalid rule argument! Nubers needs to be between 1 and 9.")
		}

		bit := byte(1 << uint(char-'1'))
		if (newRuleBirth & bit) == 0 {
			newRuleSurvive += bit
		}
	}

	f.ruleSurvive = newRuleSurvive
	f.ruleBirth = newRuleBirth

	return nil
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
	var x, y int32
	for cellIndex := workerIndex - 1; cellIndex < f.CellCount(); cellIndex += MAXWORKERINDEX {
		x, y = f.toReal(cellIndex)
		neighbors := f.countNeighbors(x, y)

		if neighbors == 0 {
			f.nextArea[cellIndex] = 0
		} else {
			if neighborsBit := byte(1 << (neighbors - 1)); (f.ruleSurvive & neighborsBit) == neighborsBit {
				f.nextArea[cellIndex] = f.currentArea[cellIndex]
			} else if (f.ruleBirth & neighborsBit) == neighborsBit {
				f.nextArea[cellIndex] = 1
			} else {
				f.nextArea[cellIndex] = 0
			}
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
