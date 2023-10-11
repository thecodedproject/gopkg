package proto_conversion

import (
	shopspring_decimal "github.com/shopspring/decimal"
	"strconv"
)

func IntAsStringFromProto(v *IntAsString) (int, error) {

	return strconv.Atoi(v.Value)
}

func IntAsStringToProto(i int) (*IntAsString, error) {

	return &IntAsString{
		Value: strconv.Itoa(i),
	}, nil
}

func ShopspringDecimalFromProto(v *ShopspringDecimal) (shopspring_decimal.Decimal, error) {

	return shopspring_decimal.NewFromString(v.Value)
}

func ShopspringDecimalToProto(v shopspring_decimal.Decimal) (*ShopspringDecimal, error) {

	return &ShopspringDecimal{
		Value: v.String(),
	}, nil
}
