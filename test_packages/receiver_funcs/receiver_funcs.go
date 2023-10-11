package receiver_funcs

import ()

type MyType struct {
}

func (m MyType) ValueReceiverFunc() {
	return
}

func (m *MyType) PointerRecFunc() {
	return
}

type OtherType bool

func (o *OtherType) OtherPRecFunc() {
	return
}

func (o OtherType) SomeOtherValRec() {
	return
}
