package gopkg_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	//"github.com/thecodedproject/gopkg/tmpl"
)

func TestWriteDeclVars(t *testing.T) {

	testCases := []struct{
		Name string
		Keyword string
		Vars []gopkg.DeclVar
		ImportAliases map[string]string
		ExpectedErr error
	}{
		{
			Name: "empty vars writes nothing",
		},
		{
			Name: "one var with empty name returns error",
			ExpectedErr: errors.New("WriteDeclVar: DeclVar.Name cannot be empty"),
			Vars: []gopkg.DeclVar{
				{
					Name: "a",
					Type: gopkg.TypeInt{},
				},
				{
				},
			},
		},
		{
			Name: "one var with empty type and literal value returns error",
			Vars: []gopkg.DeclVar{
				{
					Name: "a",
					Type: gopkg.TypeInt{},
				},
				{
					Name: "b",
				},
			},
			ExpectedErr: errors.New("WriteDeclVar: one of DeclVar.Type and DeclVar.LiteralValue must be set"),
		},
		{
			Name: "single var with built in type",
			Keyword: "var",
			Vars: []gopkg.DeclVar{
				{
					Name: "a",
					Type: gopkg.TypeInt{},
				},
			},
		},
		{
			Name: "single var with named type and import alias",
			Keyword: "var",
			Vars: []gopkg.DeclVar{
				{
					Name: "a",
					Type: gopkg.TypeNamed{
						Name: "MyStruct",
						Import: "some/import/pkg",
					},
				},
			},
			ImportAliases: map[string]string{
				"some/import/pkg": "pkg",
			},
		},
		{
			Name: "single var with no type and literal value",
			Keyword: "var",
			Vars: []gopkg.DeclVar{
				{
					Name: "myVal",
					LiteralValue: "123",
				},
			},
		},
		{
			Name: "single var with unknown literal type and literal value",
			Keyword: "var",
			Vars: []gopkg.DeclVar{
				{
					Name: "aVar",
					Type: gopkg.TypeUnnamedLiteral{},
					LiteralValue: `"my string"`,
				},
			},
		},
		{
			Name: "single var with type and literal value",
			Keyword: "var",
			Vars: []gopkg.DeclVar{
				{
					Name: "myOtherVal",
					Type: gopkg.TypeString{},
					LiteralValue: `"hello world"`,
				},
			},
		},
		{
			Name: "multiple vars",
			Keyword: "var",
			Vars: []gopkg.DeclVar{
				{
					Name: "myVal",
					LiteralValue: "123",
				},
				{
					Name: "someInt",
					Type: gopkg.TypeInt64{},
				},
				{
					Name: "cmdArg",
					Type: gopkg.TypeNamed{
						Name: "flag",
						Import: "flag",
					},
					LiteralValue: `flag.String("some_arg", "", "set an arg")`,
				},
			},
		},
		{
			Name: "multiple consts",
			Keyword: "const",
			Vars: []gopkg.DeclVar{
				{
					Name: "myVal",
					LiteralValue: "123",
				},
				{
					Name: "someInt",
					Type: gopkg.TypeInt64{},
				},
				{
					Name: "cmdArg",
					Type: gopkg.TypeNamed{
						Name: "flag",
						Import: "flag",
					},
					LiteralValue: `flag.String("some_arg", "", "set an arg")`,
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			buffer := bytes.NewBuffer(nil)

			err := gopkg.WriteDeclVars(
				buffer,
				test.Keyword,
				test.Vars,
				test.ImportAliases,
			)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), buffer.Bytes())
		})
	}

}
