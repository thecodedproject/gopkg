package interfaces

import (
	"context"
)

//go:generate go run ../main.go --interface MyInterface --type MyType
type MyInterface interface {
	SomeFunc(ctx context.Context) error
	SomeOtherFunc(i int, b string) []string
}
