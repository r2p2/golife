/**
 * User: r2p2
 * Date: 12/29/11
 * Time: 12:02 AM
 */
package main

import (
	"fmt"
	"time"
)

func main() {
	gol, error := NewFieldFromMap("glider")
	if error != nil {
		fmt.Println(error)
		return
	}
	//gol := NewField(25, 25)
	fmt.Println(gol)
	//Sgol.Initialize(0.5)

	fmt.Println(gol)

	for {
		gol.Step()
		fmt.Print("\033[2J")
		fmt.Println(gol)
		time.Sleep(0.125e9)
	}
}
