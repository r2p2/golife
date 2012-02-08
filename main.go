/**
 * User: r2p2
 * Date: 12/29/11
 * Time: 12:02 AM
 */
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"
)

var delayMs *int = flag.Int("delay", 125, "delay between iterations in ms")
var rule *string = flag.String("rule", "23/3", "modifies the gol rule")
var width *int = flag.Int("width", 70, "width of a new generated map")
var height *int = flag.Int("height", 30, "height of a new generated map")
var fillRate *int = flag.Int("fill", 30, "fill rate of a new generated map 0-100")
var loadFile *string = flag.String("file", "", "load a map from file")
var loadMap *string = flag.String("map", "", "load a preinstalled map")
var listMaps *bool = flag.Bool("list", false, "list the name of preinstalled maps")

func fpsCounter() func() uint32 {
	var nextFPS, currentFPS uint32
	var timestamp time.Time
	fpsDelay, _ := time.ParseDuration("1s")

	nextFPS = 1
	currentTime := time.Now()
	timestamp = currentTime.Add(fpsDelay)

	return func() uint32 {
		if now := time.Now(); timestamp.Sub(now) <= 0 {
			currentFPS = nextFPS
			nextFPS = 1
			timestamp = now.Add(fpsDelay)
		} else {
			nextFPS++
		}
		return currentFPS
	}
}

func main() {
	flag.Parse()
	var gol *Field
	delayNs := int64(*delayMs * 1e6)

	if *listMaps == true {
		for name := range maps {
			fmt.Println(name)
		}
		return
	}

	if *loadMap != "" {
		var error error
		gol, error = NewFieldFromMap(*loadMap)
		if error != nil {
			fmt.Println(error)
			return
		}
	} else if *loadFile != "" {
		contents, error := ioutil.ReadFile(*loadFile)
		if error != nil {
			fmt.Println(error)
			return
		}
		gol = NewFieldFromString(string(contents))
	} else {
		gol = NewField(int32(*width), int32(*height))
		gol.Initialize(float32(*fillRate) / 100)
	}

	if error := gol.SetRule(*rule); error != nil {
		fmt.Println(error)
		return
	}

	fmt.Print("\033[2J")
	fps := fpsCounter()
	for {
		fmt.Printf("\033[%dA", gol.Height()+2)
		fmt.Println("iteration: ", gol.Iteration(), " iterations per second: ", fps(), "   ")
		fmt.Println(gol)
		gol.Step()
		time.Sleep(time.Duration(delayNs))
	}
}
