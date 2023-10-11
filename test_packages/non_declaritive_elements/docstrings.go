// A package level doc string
// with
//
// multiple lines
package non_declaritive_elements

import (
)

// singleVar has a docstring
// with multiple lines
var singleVar int64

const (
	// a doc string on a group of consts
	firstConstant int32 = 1
	// another with
	// several lines
	secondConstant string  = "hello"
	// some comment on multiple values
	thirdC, fourthC = 10, 12
)

// comment on var grounds are ignored
var (
	// only docstrings inside var groups are kept
	someVar int
)

// adocumentedType with a docstring
// and multiple lines
// of text
type adocumentedType string

// a docstring not starting with the method name
func ExportedMethod() {
}

// unexported method docstring
func unexportedMethod() {
}
