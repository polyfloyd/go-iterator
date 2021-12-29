Iterators for Go
================

This is a little experiment to see how iterators would look like with the generics introduced in Go
1.18. The design has been mostly taken from Rust's Iterator with some minor alterations to make it
more Go-like.

This library should be considered an Alpha version. Feel free to toy around, but please do not use
it in things that are deemed important.

It has functions to create iterators from Go types such as slices, channels, integer ranges, etc. as
well as the standard operations one would expect such as Map, Reduce, Filter, Flatten and more.

```go
package main

import (
	"fmt"

	"github.com/polyfloyd/go-iterator"
)

func main() {
	numbers := iterator.Range(0, 10, 1)
	evenNumbers := iterator.Filter(numbers, func(i int) bool {
		return i%2 == 0
	})
	numberStrings := iterator.Map(evenNumbers, func(i int) string {
		return fmt.Sprint(i)
	})
	str := iterator.Join(numberStrings, ", ")
	fmt.Println(str) // 0, 2, 4, 6, 8
}
```

The syntax of chaining iterators is somewhat unwieldy because Go does not allow member functions of
interfaces to introduce new type parameters.
