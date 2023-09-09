package non_declaritive_elements

import (
)

// singleVar has a docstring
// with multiple lines
var singleVar int64

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
