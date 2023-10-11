package non_declaritive_elements

import ()

type someType struct{}

func SomeFunc() int64 {
	// A comment...
	var a int64
	a = 1234
	b := a
	return b
}

func (s someType) receiverMethod() bool {
	//some reciever method comment...
	return true
}

func unexportedFunc() string {
	// some other comment...
	return "foobar"
}
