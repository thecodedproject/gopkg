package custom_types

import (
	c "context"
)

type SingleEmbed struct {
	c.Context
}

type InterfaceEmbed interface {
	float64
	MyFunc()
}

type ManyEmbeds struct {
	error
	c.Context

	myVar string

	int32
}

type InterfaceManyEmbeds interface {
	SingleEmbed
	InterfaceEmbed
	c.Context
	error
}
