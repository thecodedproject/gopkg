package dependent_types

import (
	"math/big"

	"github.com/shopspring/decimal"

	"github.com/thecodedproject/gopkg/test_packages/dependent_types/nested_pkg"
)

type AStruct struct {
	One big.Int
	Two decimal.Decimal
	Three nested_pkg.SomeType
	Four nested_pkg.AnotherType
}
