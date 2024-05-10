package nested_pkg

import (
	"context"
)

type SomeType *int32

type AnotherType struct {
	A string
	B context.Context
}
