package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	export = flag.String("o", "new.gif", "output file")
)

func main() {
	flag.Parse()
	var env Environment
	env.setup()
	for i := 0; i < numOfIterations; i++ {
		env.run(i)
		fmt.Printf("[%d/%d]\n", i+1, numOfIterations)
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
