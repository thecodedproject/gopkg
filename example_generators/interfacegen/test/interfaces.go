package interfaces

import (
	"context"
)

//go:generate go run ../main.go --interface MyInterface --type MyType --import github.com/neurotempest/gopkg/example_generators/interfacegen/test
type MyInterface interface {
	SomeFunc(ctx context.Context) error
	SomeOtherFunc(i int, b string) []string
}
