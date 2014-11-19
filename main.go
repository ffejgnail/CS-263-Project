package main

import (
	"fmt"
)

func main() {
	var env Environment
	env.setup()
	for i := 0; i < numOfIterations; i++ {
		env.run(i)
	}
	env.print()
	fmt.Println("end")
}
