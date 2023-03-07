package custom_types

import (
	c "context"
)

type SingleEmbed struct {
	c.Context
}

type ManyEmbeds struct {
	error
	c.Context

	myVar string

	int32
}
