package composite_types

import (
	"github.com/shopspring/decimal"
)

type SomeMapStruct struct {
	MOfInts       map[int64]int64
	MOfPToStrings map[string]*string
	MOfDToArrayOfInt map[decimal.Decimal][]int32
}

type MyCustomMapType map[int]float64

type SomeMapInterface interface {
	MapMaker(n int64, vals string) map[int64]string
}

func SomeMapFunc(
	a map[int64]decimal.Decimal,
	b map[*string][]float32,
) map[string]SomeMapStruct {
	return nil
}

