package main

import (
	"flag"
	"fmt"
)

var (
	export     = flag.String("o", "new.gif", "output file")
	brainData  = flag.String("b", "data.json", "brain json file")
	trainBrain = flag.Bool("t", false, "train brain")
)

func main() {
	flag.Parse()
	if *trainBrain {
		brain := NewRBMBrain()

		brain.Default()
		brain.Dump(*brainData)
		return
	}
	env := NewEnvironment()
	for i := 0; i < Iteration; i++ {
		fmt.Printf("[%d/%d]\n", i+1, Iteration)
		env.Run(i)
	}
	env.WriteFile(*export)
	fmt.Println("exported to", *export)
}
