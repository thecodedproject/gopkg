package struct_with_tags

import ()

type AStruct struct {
	AField       int     `AKey:"some_value"`
	BField       bool    `BKey:"some_other_value"`
	privateField float32 `CKey:"some_third_value"`
}
