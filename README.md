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

## Limitations

### Chaining calls
This is not possible to write:
```go
str := iterator.Range(0, 10, 1).
	Filter(func(i int) bool {
		return i%2 == 0
	}).
	Map(func(i int) string {
		return fmt.Sprint(i)
	}).
	Join(", ")
fmt.Println(str) // 0, 2, 4, 6, 8
```
Because that would require methods with generic types that are not tied to the Iterator interfaces
itself to be present on the Iterator interface, which is not possible because templated interfaces
are still dynamic types. It's like why you can not have a `dyn Hash` in Rust.

Making this possible requires a change to the language to permit interfaces to have methods declared
on them. Either by reusing the existing syntax for methods or by providing some kind of reverse
Method Expression.
