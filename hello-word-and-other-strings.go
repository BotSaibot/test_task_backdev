package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 3; i++ {
		fmt.Printf("I\tlike\tGo!\n")
	}
	fmt.Println("Hello Go"[0])
	fmt.Println(string("Hello Go"[0]))
	fmt.Println(`a ...any
b ...any
c ...any`)
	// myv := 0
	// fmt.Fprintf(io.Writer, "format string %i\n", myv)
	// fmt.Fprintf(w io.Writer, format string, a ...any)
	myVar := 0
	fmt.Print(myVar)
}
