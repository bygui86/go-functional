package main_test

import (
	"math/rand"
	"testing"

	"go-functional/pipeline"
)

var result int // to prevent compiler optimization

func BenchmarkPipe(b *testing.B) {
	sqr := func(x int) int { return x * x }
	inc := func(x int) int { return x + 1 }
	sink := func(x int) { result = x }

	x := rand.Intn(1000)

	b.Run("PowerPlusOneDirect", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			result = inc(sqr(x))
		}
	})

	b.Run("PowerPlusOnePipe", func(b *testing.B) {
		pipe := pipeline.Pipe(sqr, inc, sink)
		for n := 0; n < b.N; n++ {
			pipe(x)
		}
	})
}
