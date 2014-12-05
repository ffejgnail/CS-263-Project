package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	export    = flag.String("o", "new.gif", "output file")
	brainData = flag.String("b", "data.json", "brain json file")
)

func main() {
	flag.Parse()
	env := NewEnvironment()
	for i := 0; i < Iteration; i++ {
		if Iteration-i <= RecordIteration {
			fmt.Printf("[%d/%d]\n", i+1, Iteration)
		}
		env.Run(i)
	}
	f, err := os.Create(*export)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	env.WriteTo(f)
	fmt.Println("exported to", *export)
}
