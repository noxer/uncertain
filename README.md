# uncertain
Deal with uncertain data structures in Go.

This package allows you to traverse unknown data structures without complicated type assertions or dealing with the reflect package. It provides a simple `Get` function to retrieve nested data.

## docs
https://pkg.go.dev/github.com/noxer/uncertain?tab=doc

## example

```go
package main

import (
	"fmt"

	"github.com/noxer/uncertain"
)

type MapElement struct {
	List []int
}

type ComplexStruct struct {
	Label   string
	private string
	Values  map[string]interface{}
}

func main() {
	// create a complex struct
	cs := &ComplexStruct{
		Label:   "Easy",
		private: "inaccessible!",
		Values: map[string]interface{}{
			"first": MapElement{
				List: []int{0, 1, 2, 3, 4, 5},
			},
			"second": &MapElement{
				List: []int{5, 4, 3, 2, 1, 0},
			},
			"third": "something else entirely",
		},
	}

	// the "<nil>" in the output indicates that a nil was returned

	fmt.Println(uncertain.Get(cs, "Label"))
	// prints "Easy <nil>"

	fmt.Println(uncertain.Get(cs, "Label", 1))
	// prints "97 <nil>"
	// 97 is the Unicode value of the 'a' in "Easy"

	fmt.Println(uncertain.Get(cs, "private"))
	// prints "<nil> can't access private field"
	// private fields of structs can't be accessed with this library

	fmt.Println(uncertain.Get(cs, "Values", "first", "List", 1))
	// prints "1 <nil>"

	fmt.Println(uncertain.Get(cs, "Values", "second", "List", 1))
	// prints "4 <nil>"

	fmt.Println(uncertain.Get(cs, "Values", "third", "List", 1))
	// prints "<nil> path segment can't be interpreted as a number"
	// the third map element is a string, so Get tries to interpret "List" as an index for the string

	fmt.Println(uncertain.Get(cs, "Values", "third", "2"))
	// prints "109 <nil>"
	// a string of "2" will be parsed as a decimal number and used as an index for the string
	// 109 is the Unicode value for the 'm' in "something"
}
```
