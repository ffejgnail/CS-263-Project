package main

import (
	"fmt"
	"os"
)

func main() {
	var env Environment
	env.setup()
	for i := 0; i < numOfIterations; i++ {
		env.run(i)
	}
	f, err := os.Create("new.gif")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	env.WriteTo(f)
	fmt.Println("end")
}
