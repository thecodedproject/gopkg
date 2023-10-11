package composite_types

import (
	shopspring_decimal "github.com/shopspring/decimal"
)

type SomeArrayStruct struct {
	AOfInts       []int64
	AOfPToStrings []*string
}

type MyCustomArrayType [][][]float64

type SomeArrayInterface interface {
	ArrayMaker(n int64, vals string) []string
}

func SomeArrayFunc(
	a []shopspring_decimal.Decimal,
	b []float32,
) []SomeArrayStruct {

	return nil
}
