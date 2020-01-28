package pipeline

import (
	"errors"
	"fmt"
	"reflect"
)

/*
	As Go does not support overloading and/or defining new infix operators, we have to implement the pipe operator as a function.
	We start by defining Pipe function and Pipeline type for its result.

	This presents a rather simple method to implement a pipe operator in Go. It works, but it also eliminates
	type-safety. It casts all arguments into interface{} before the execution of the pipeline, and it uses reflection
	to apply functions. So, there is no way to check functions and their arguments at compile-time.

	The benchmark in this repo shows that (on my machine): while the direct function call results 2/3 ns/op,
	Pipe runs in ~600/700 ns/op which is 300–350× worst.
	This makes it impractical for repetitive small tasks, as it has a relatively big overhead because of the
	reflection. Although it is specially useful for medium-sized tasks.
*/

/*
	A Pipeline instance is just another function that does the actual work. It accepts zero or more inputs and gives an error.

	The number of its input arguments must match the input arguments of the first function in fs.
	But its output may or may not match the last function, and that is because Go does not have variadic return values.
*/
type Pipeline func(...interface{}) error

/*
	errType is the type of error interface.
*/
var errType = reflect.TypeOf((*error)(nil)).Elem()

/*
	Pipe uses variadic arguments to take zero or more functions as inputs and produces a Pipeline.
*/
func Pipe(fns ...interface{}) Pipeline {
	if len(fns) == 0 {
		return emptyFn
	}

	return func(args ...interface{}) (err error) {
		/*
			?
		*/
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("pipeline does panic: %v", r)
			}
		}()

		/*
			We start by processing the input arguments.
			Go reflection enables us to call functions dynamically, but for that, we need to pass the input arguments
			as an slice of reflect.Value values. So first, we need to do the conversion.
		*/
		var inputs []reflect.Value
		for _, arg := range args {
			inputs = append(inputs, reflect.ValueOf(arg))
		}

		/*
			Secondly, we have to solve the nested function calls, f(g(..)). We can always unwind nested
			functions using a for-loop.
			We don't have to forget to handle possible errors, without passing them down to the pipeline.
		*/
		for fnIndex, fn := range fns {
			// call the function
			outputs := reflect.ValueOf(fn).Call(inputs)
			// ?
			inputs = inputs[:0]

			fnType := reflect.TypeOf(fn)

			for outputIndex, output := range outputs {
				if fnType.Out(outputIndex).Implements(errType) {
					/*
						Pay attention on how to check return values for errors, neglecting the nil ones.
					*/
					if !output.IsNil() {
						err = fmt.Errorf("%s func failed: %w", ord(fnIndex), output.Interface().(error))
						return
					}
				} else {
					inputs = append(inputs, output)
				}
			}
		}

		return
	}
}

/*
	emptyFn returning nil error
*/
func emptyFn(...interface{}) error {
	return nil
}

/*
	?
*/
func ord(index int) string {
	order := index + 1
	switch {
	case order > 10 && order < 20:
		return fmt.Sprintf("%dth", order)
	case order%10 == 1:
		return fmt.Sprintf("%dst", order)
	case order%10 == 2:
		return fmt.Sprintf("%dnd", order)
	case order%10 == 3:
		return fmt.Sprintf("%drd", order)
	default:
		return fmt.Sprintf("%dth", order)
	}
}

func powerPlusOne(x int) (int, error) {
	var result int

	/*
		Here you see how easy it is to compose functions with Pipe.
		There is no need to handle errors for each function and the list of composing functions can go on indefinitely.
	*/
	err := Pipe(
		// power
		func(x int) (int, error) {
			if x < 0 {
				return 0, errors.New("x should not be negative")
			}
			return x * x, nil
		},
		// PlusOne
		func(x int) int { return x + 1 },
		/*
			We call the last function of the pipeline a sink, i.e. its job is to gather the results and clean the pipeline.
			The sink should not return any values other than an optional error.
		*/
		func(x int) { result = x },
	)(x) // the execution of pipeline

	if err != nil {
		return -1, err
	}

	return result, nil
}

func PossibleSolution() {
	res, err := powerPlusOne(5)
	if err != nil {
		panic(err)
	}
	fmt.Println("POWER+1:", res)
	fmt.Println("")
}
