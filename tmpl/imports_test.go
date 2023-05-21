package tmpl_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestNnnamedImports(
	t *testing.T,
) {

	testCases := []struct{
		Name string
		Input []string
		Expected []gopkg.ImportAndAlias
	}{
		{
			Name: "Empty input",
		},
		{
			Name: "Several inputs",
			Input: []string{
				"strings",
				"some/import/path",
				"github.com/mypkg",
			},
			Expected: []gopkg.ImportAndAlias{
				{
					Import: "strings",
				},
				{
					Import: "some/import/path",
				},
				{
					Import: "github.com/mypkg",
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			actual := tmpl.UnnamedImports(test.Input...)

			if len(test.Expected) == 0 {
				require.Equal(t, len(test.Expected), len(actual))
				return
			}

			require.Equal(t, test.Expected, actual)
		})
	}
}
