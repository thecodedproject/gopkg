package pkg_with_tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMyCoolLogic(t *testing.T) {

	testCases := []struct{
		Name string
		I int
		J int
		Expected int
	}{
		{
			Name: "empty returns zero",
		},
		{
			Name: "adds inputs together",
			I: 1,
			J: 2,
			Expected: 3,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			require.Equal(
				t,
				test.Expected,
				pkg_with_tests_test.MyCoolLogic(i, j),
			)
		})
	}
}
