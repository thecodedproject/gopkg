package dependent_types

import (
	"math/big"

	"github.com/shopspring/decimal"
)

type AStruct struct {
	One big.Int
	Two decimal.Decimal
}
