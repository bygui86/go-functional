package main

import (
	"fmt"

	. "go-functional/pipeline"
	. "go-functional/quickAndDirty"
)

func main() {
	fmt.Println("*** Quick & dirty ***")
	QuickAndDirty()

	fmt.Println("*** Pipeline ***")
	PossibleSolution()
}
