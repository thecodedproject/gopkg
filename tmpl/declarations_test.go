package tmpl_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/neurotempest/gopkg"
	"github.com/neurotempest/gopkg/tmpl"
)

func TestUnnamedReturnArgs(t *testing.T) {

	testCases := []struct{
		Name string
		RetArgs []gopkg.Type
		Expected []gopkg.DeclVar
	}{
		{
			Name: "empty returns empty list",
			Expected: []gopkg.DeclVar{},
		},
		{
			Name: "list of args",
			RetArgs: []gopkg.Type{
				gopkg.TypeInt64{},
				gopkg.TypeBool{},
				gopkg.TypeString{},
			},
			Expected: []gopkg.DeclVar{
				{Type: gopkg.TypeInt64{}},
				{Type: gopkg.TypeBool{}},
				{Type: gopkg.TypeString{}},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			require.Equal(
				t,
				test.Expected,
				tmpl.UnnamedReturnArgs(test.RetArgs...),
			)
		})
	}
}
