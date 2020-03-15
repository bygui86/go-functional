package main

import (
	"fmt"

	"go-functional/pipeline"
	"go-functional/quickAndDirty"
)

func main() {
	fmt.Println("*** Quick & dirty ***")
	quickAndDirty.QuickAndDirty()

	fmt.Println("*** Pipeline ***")
	pipeline.PossibleSolution()
}
