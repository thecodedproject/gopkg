package tmpl_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestFuncWithContextAndError(t *testing.T) {

	testCases := []struct{
		Name string
		FuncName string
		Args []gopkg.DeclVar
		RetArgs []gopkg.DeclVar
		Expected gopkg.DeclFunc
	}{
		{
			Name: "empty returns func with ctx and unnamed err",
			Expected: gopkg.DeclFunc{
				Args: []gopkg.DeclVar{
					{
						Name: "ctx",
						Type: gopkg.TypeUnknownNamed{
							Name: "Context",
							Import: "context",
						},
					},
				},
				ReturnArgs: []gopkg.DeclVar{
					{
						Type: gopkg.TypeError{},
					},
				},
			},
		},
		{
			Name: "with some args and unnamed return args",
			FuncName: "MyMethod",
			Args: []gopkg.DeclVar{
				{
					Name: "a",
					Type: gopkg.TypeInt{},
				},
				{
					Name: "b",
					Type: gopkg.TypeString{},
				},
			},
			RetArgs: []gopkg.DeclVar{
				{
					Type: gopkg.TypeFloat32{},
				},
				{
					Type: gopkg.TypeBool{},
				},
			},
			Expected: gopkg.DeclFunc{
				Name: "MyMethod",
				Args: []gopkg.DeclVar{
					{
						Name: "ctx",
						Type: gopkg.TypeUnknownNamed{
							Name: "Context",
							Import: "context",
						},
					},
					{
						Name: "a",
						Type: gopkg.TypeInt{},
					},
					{
						Name: "b",
						Type: gopkg.TypeString{},
					},
				},
				ReturnArgs: []gopkg.DeclVar{
					{
						Type: gopkg.TypeFloat32{},
					},
					{
						Type: gopkg.TypeBool{},
					},
					{
						Type: gopkg.TypeError{},
					},
				},
			},
		},
		{
			Name: "with named return args",
			FuncName: "MyMethod",
			RetArgs: []gopkg.DeclVar{
				{
					Name: "a",
					Type: gopkg.TypeFloat32{},
				},
			},
			Expected: gopkg.DeclFunc{
				Name: "MyMethod",
				Args: []gopkg.DeclVar{
					{
						Name: "ctx",
						Type: gopkg.TypeUnknownNamed{
							Name: "Context",
							Import: "context",
						},
					},
				},
				ReturnArgs: []gopkg.DeclVar{
					{
						Name: "a",
						Type: gopkg.TypeFloat32{},
					},
					{
						Name: "err",
						Type: gopkg.TypeError{},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			require.Equal(
				t,
				test.Expected,
				tmpl.FuncWithContextAndError(
					test.FuncName,
					test.Args,
					test.RetArgs,
				),
			)
		})
	}
}

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
