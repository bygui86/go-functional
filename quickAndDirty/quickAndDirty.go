package quickAndDirty

import (
	"errors"
	"fmt"
)

/*
This is somewhat straightforward for two functions. But it can quickly become a burden when the number of functions increase
*/

func QuickAndDirty() {
	fmt.Println("POWER+1:", powerPlusOne(5))
	res, err := powerPlusOneHandlingError(5)
	if err != nil {
		panic(err)
	}
	fmt.Println("POWER+1:", res)
	fmt.Println("")
}

// Apply functions consequently
func powerPlusOne(x int) int {
	sqr := func(x int) int { return x * x }
	inc := func(x int) int { return x + 1 }
	return inc(sqr(x))
}

// Apply functions consequently handling functions returning multiple values and errors
func powerPlusOneHandlingError(x int) (int, error) {
	sqr := func(x int) (int, error) {
		if x < 0 {
			return 0, errors.New("x should not be negative")
		}
		return x * x, nil
	}
	inc := func(x int) int { return x + 1 }
	y, err := sqr(x)
	if err != nil {
		return 0, err
	}
	return inc(y), nil
}
