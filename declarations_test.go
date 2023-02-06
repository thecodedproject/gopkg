package gopkg_test

import (
	"testing"

	//"github.com/stretchr/testify/require"

	"github.com/neurotempest/gopkg"
)

func TestDeclVar(t *testing.T) {

	_ = gopkg.DeclVar{
		Name: "a",

		Type: gopkg.TypeUnknownNamed{
			Name: "string",
		},
	}
}
