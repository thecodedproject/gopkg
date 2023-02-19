package gopkg_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/neurotempest/gopkg"
)

func TestPackageImportPath(t *testing.T) {

	testCases := []struct{
		Name string
		Path string
		Expected string
		ExpectedErr error
	}{
		{
			Name: "empty string returns error",
			ExpectedErr: errors.New("cannot get import path for empty path"),
		},
		{
			Name: "path which does not exist returns error",
			Path: "some/unknown/path",
			ExpectedErr: errors.New("no such path `some/unknown/path`"),
		},
		{
			Name: "current module",
			Path: ".",
			Expected: "github.com/neurotempest/gopkg",
		},
		{
			Name: "sub package of current module without dot",
			Path: "test_packages/composite_types",
			Expected: "github.com/neurotempest/gopkg/test_packages/composite_types",
		},
		{
			Name: "sub package of current module with dot",
			Path: "./test_packages/all_built_in_types",
			Expected: "github.com/neurotempest/gopkg/test_packages/all_built_in_types",
		},
		{
			Name: "existing file within sub package",
			Path: "./test_packages/proto_conversion/converters.go",
			Expected: "github.com/neurotempest/gopkg/test_packages/proto_conversion",
		},
		{
			Name: "sub package of current module with trailing slash",
			Path: "test_packages/composite_types/",
			Expected: "github.com/neurotempest/gopkg/test_packages/composite_types",
		},
		{
			Name: "sub package of current module double-dot in path",
			Path: "test_packages/../test_packages/composite_types/../composite_types",
			Expected: "github.com/neurotempest/gopkg/test_packages/composite_types",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := gopkg.PackageImportPath(test.Path)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, test.Expected, actual)
		})
	}
}
