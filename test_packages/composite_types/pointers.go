package composite_types

import (
)

type SomePointerStruct struct{
	PToInt *int32
}

func SomePointerFunc(a *float32, b *SomePointerStruct) *string {

	return nil
}
