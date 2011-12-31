/**
 * User: r2p2
 * Date: 12/29/11
 * Time: 12:02 AM
 */
package main

import (
	"fmt"
	"time"
	"flag"
	"io/ioutil"
	"os"
)

var delayMs *int = flag.Int("delay", 125, "delay between iterations in ms")
var width *int = flag.Int("width", 70, "width of a new generated map")
var height *int = flag.Int("height", 30, "height of a new generated map")
var fillRate *int = flag.Int("fill", 30, "fill rate of a new generated map 0-100")
var loadFile *string = flag.String("file", "", "load a map from file")
var loadMap *string = flag.String("map", "", "load a preinstalled map")
var listMaps *bool = flag.Bool("list", false, "list the name of preinstalled maps")

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
		var error os.Error
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

	fmt.Print("\033[2J")
	for {
		fmt.Printf("\033[%dA", gol.Height()+2)
		fmt.Println("Iteration: ", gol.Iteration())
		fmt.Println(gol)
		gol.Step()
		time.Sleep(delayNs)

	}
}
