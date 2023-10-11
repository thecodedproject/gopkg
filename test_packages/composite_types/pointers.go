package composite_types

import ()

type MyCustomPointer **float32

type SomePointerInterface interface {
	Something() *int64
	PointerMaker(val *string) *float64
}

type SomePointerStruct struct {
	PToInt *int32
}

func SomePointerFunc(a *float32, b *SomePointerStruct) *string {

	return nil
}
