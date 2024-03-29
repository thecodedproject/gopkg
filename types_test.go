package gopkg_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestTypeDefaultInit(t *testing.T) {

	testCases := []struct {
		Def           gopkg.Type
		ImportAliases map[string]string
		Expected      string
		ExpectedErr   error
	}{
		{
			Def:      gopkg.TypeAny{},
			Expected: "nil",
		},
		{
			Def:      gopkg.TypeBool{},
			Expected: "false",
		},
		{
			Def:      gopkg.TypeByte{},
			Expected: "0",
		},
		{
			Def:      gopkg.TypeError{},
			Expected: "nil",
		},
		{
			Def:      gopkg.TypeFloat32{},
			Expected: "0",
		},
		{
			Def:      gopkg.TypeFloat64{},
			Expected: "0",
		},
		{
			Def:      gopkg.TypeFunc{},
			Expected: "nil",
		},
		{
			Def:      gopkg.TypeInt{},
			Expected: "0",
		},
		{
			Def:      gopkg.TypeInt32{},
			Expected: "0",
		},
		{
			Def:      gopkg.TypeInt64{},
			Expected: "0",
		},
		{
			Def:      gopkg.TypeString{},
			Expected: "\"\"",
		},
		{
			Def:      gopkg.TypeInterface{},
			Expected: "nil",
		},
		{
			Def:      gopkg.TypeStruct{},
			Expected: "{}",
		},
		{
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeNamed{
					Name:   "SomeType",
					Import: "some/import",
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeNamed{
					Name: "float32",
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeNamed{
					Name:   "MyType",
					Import: "my/import/path",
				},
				ValueType: gopkg.TypePointer{
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeNamed{
							Name:   "MyOtherType",
							Import: "other/import",
						},
					},
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeNamed{
				Name:   "MyType",
				Import: "my/import/path",
			},
			ExpectedErr: errors.New("cannot deduce default init for named type with no value type"),
		},
		{
			Def:         gopkg.TypeUnnamedLiteral{},
			ExpectedErr: errors.New("no default init for unnamed literal"),
		},
	}

	for _, test := range testCases {

		testName, err := test.Def.FullType(nil)
		require.NoError(t, err)

		t.Run(testName, func(t *testing.T) {
			actual, err := test.Def.DefaultInit(test.ImportAliases)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.Expected, actual)
		})

		t.Run(
			"TypeUnknownNamed_with_value_type"+testName,
			func(t *testing.T) {

				namedType := gopkg.TypeNamed{
					ValueType: test.Def,
				}

				actual, err := namedType.DefaultInit(test.ImportAliases)

				if test.ExpectedErr != nil {
					require.Equal(t, test.ExpectedErr, err)
					return
				}

				require.NoError(t, err)
				require.Equal(t, test.Expected, actual)
			})
	}
}

func TestTypeFullType(t *testing.T) {

	testCases := []struct {
		Def           gopkg.Type
		ImportAliases map[string]string
		Expected      string
	}{
		{
			Def:      gopkg.TypeAny{},
			Expected: "any",
		},
		{
			Def:      gopkg.TypeBool{},
			Expected: "bool",
		},
		{
			Def:      gopkg.TypeByte{},
			Expected: "byte",
		},
		{
			Def:      gopkg.TypeError{},
			Expected: "error",
		},
		{
			Def:      gopkg.TypeFloat32{},
			Expected: "float32",
		},
		{
			Def:      gopkg.TypeFloat64{},
			Expected: "float64",
		},
		{
			Def:      gopkg.TypeInt{},
			Expected: "int",
		},
		{
			Def:      gopkg.TypeInt32{},
			Expected: "int32",
		},
		{
			Def:      gopkg.TypeInt64{},
			Expected: "int64",
		},
		{
			Def:      gopkg.TypeString{},
			Expected: "string",
		},
		{
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeNamed{
					Name:   "SomeType",
					Import: "some/import",
				},
			},
			ImportAliases: map[string]string{
				"some/import": "some_alias",
			},
			Expected: "[]some_alias.SomeType",
		},
		{
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeNamed{
					Name: "float32",
				},
			},
			Expected: "*float32",
		},
		{
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeNamed{
					Name:   "MyType",
					Import: "my/import/path",
				},
				ValueType: gopkg.TypePointer{
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeNamed{
							Name:   "MyOtherType",
							Import: "other/import",
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"my/import/path": "path_alias",
				"other/import":   "other_alias",
			},
			Expected: "map[path_alias.MyType]*[]other_alias.MyOtherType",
		},
		{
			Def:      gopkg.TypeUnnamedLiteral{},
			Expected: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.Expected, func(t *testing.T) {
			fullType, err := test.Def.FullType(test.ImportAliases)
			require.NoError(t, err)
			require.Equal(t, test.Expected, fullType)
		})
	}
}

// TestTypeFuncFullType tests the TypeFunc.FullType functionality
// It is in it's own test fixutre as it's more complicated than the other
// types
func TestTypeFuncFullType(t *testing.T) {

	testCases := []struct {
		Name          string
		Def           gopkg.TypeFunc
		ImportAliases map[string]string
		Expected      string
		ExpectedErr   error
	}{
		{
			Name:     "empty func",
			Def:      gopkg.TypeFunc{},
			Expected: "func()",
		},
		{
			Name: "func with built in unnamed args",
			Def: gopkg.TypeFunc{
				Args: tmpl.UnnamedReturnArgs(
					gopkg.TypeAny{},
					gopkg.TypeString{},
					gopkg.TypeInt64{},
					gopkg.TypeError{},
				),
			},
			Expected: "func(any, string, int64, error)",
		},
		{
			Name: "func with built in named args",
			Def: gopkg.TypeFunc{
				Args: []gopkg.DeclVar{
					{
						Name: "first",
						Type: gopkg.TypeArray{
							ValueType: gopkg.TypeInt32{},
						},
					},
					{
						Name: "second",
						Type: gopkg.TypeBool{},
					},
				},
			},
			Expected: "func(first []int32, second bool)",
		},
		{
			Name: "func with built in unnamed return arg",
			Def: gopkg.TypeFunc{
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeString{},
				),
			},
			Expected: "func() string",
		},
		{
			Name: "func with built in named return args",
			Def: gopkg.TypeFunc{
				ReturnArgs: []gopkg.DeclVar{
					{
						Name: "first",
						Type: gopkg.TypeArray{
							ValueType: gopkg.TypeInt32{},
						},
					},
					{
						Name: "err",
						Type: gopkg.TypeError{},
					},
				},
			},
			Expected: "func() (first []int32, err error)",
		},
		{
			Name: "func with args and return args with import aliases",
			Def: gopkg.TypeFunc{
				Args: []gopkg.DeclVar{
					{
						Name: "first",
						Type: gopkg.TypeNamed{
							Name:   "SomeType",
							Import: "some/path",
						},
					},
					{
						Name: "second",
						Type: gopkg.TypeNamed{
							Name:   "SomeOtherType",
							Import: "some/other",
						},
					},
				},
				ReturnArgs: []gopkg.DeclVar{
					{
						Name: "retVal",
						Type: gopkg.TypeNamed{
							Name:   "SomeThirdType",
							Import: "some/third",
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"some/path":  "path",
				"some/other": "otheralias",
				"some/third": "third",
			},
			Expected: "func(first path.SomeType, second otheralias.SomeOtherType) (retVal third.SomeThirdType)",
		},
		{
			Name: "variadic func with single arg",
			Def: gopkg.TypeFunc{
				Args: []gopkg.DeclVar{
					{
						Name: "v",
						Type: gopkg.TypeInt64{},
					},
				},
				VariadicLastArg: true,
			},
			Expected: "func(v ...int64)",
		},
		{
			Name: "variadic func with multiple args and return args",
			Def: gopkg.TypeFunc{
				Args: []gopkg.DeclVar{
					{
						Name: "a",
						Type: gopkg.TypeInt64{},
					},
					{
						Name: "b",
						Type: gopkg.TypeInt32{},
					},
					{
						Name: "v",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypeString{},
						},
					},
				},
				VariadicLastArg: true,
			},
			Expected: "func(a int64, b int32, v ...*string)",
		},
		{
			Name: "mix of named and unnamed params returns error",
			Def: gopkg.TypeFunc{
				Args: []gopkg.DeclVar{
					{
						Name: "first",
						Type: gopkg.TypeArray{
							ValueType: gopkg.TypeInt32{},
						},
					},
					{
						Type: gopkg.TypeBool{},
					},
				},
			},
			ExpectedErr: errors.New("mix of named and unnamed func args"),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			fullType, err := test.Def.FullType(test.ImportAliases)
			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.Expected, fullType)
		})
	}
}

// TestTypeStructFullType tests the TypeStruct.FullType functionality
// It is in it's own test fixutre as it's morecomplicated than the other
// types
func TestTypeStructFullType(t *testing.T) {

	testCases := []struct {
		Name          string
		Def           gopkg.TypeStruct
		ImportAliases map[string]string
		Expected      string
	}{
		{
			Name:     "empty struct",
			Def:      gopkg.TypeStruct{},
			Expected: "struct {}",
		},
		{
			Name: "struct with fields no import aliases",
			Def: gopkg.TypeStruct{
				Fields: []gopkg.DeclVar{
					{
						Name: "SomeValue",
						Type: gopkg.TypeInt{},
					},
					{
						Name: "SomeOtherValue",
						Type: gopkg.TypeString{},
					},
				},
			},
			Expected: `struct {
	SomeValue int
	SomeOtherValue string
}`,
		},
		{
			Name: "struct with tags",
			Def: gopkg.TypeStruct{
				Fields: []gopkg.DeclVar{
					{
						Name:      "AField",
						Type:      gopkg.TypeBool{},
						StructTag: "some_unconvensional_value",
					},
					{
						Name:      "bField",
						Type:      gopkg.TypeInt32{},
						StructTag: "json:\"-\" key:\"other,tags\"",
					},
				},
			},
			Expected: `struct {
	AField bool ` + "`some_unconvensional_value`" + `
	bField int32 ` + "`json:\"-\" key:\"other,tags\"`" + `
}`,
		},
		{
			Name: "struct with embedded types only",
			Def: gopkg.TypeStruct{
				Embeds: []gopkg.Type{
					gopkg.TypeNamed{
						Name:   "MyType",
						Import: "github.com/myrepo",
					},
					gopkg.TypeNamed{
						Name:   "MyTypeTwo",
						Import: "github.com/myotherrepo",
					},
					gopkg.TypeError{},
				},
			},
			ImportAliases: map[string]string{
				"github.com/myrepo":      "myrepo_alias",
				"github.com/myotherrepo": "myotherrepo_alias",
			},
			Expected: `struct {
	myrepo_alias.MyType
	myotherrepo_alias.MyTypeTwo
	error
}`,
		},
		{
			Name: "struct with fields with import aliases",
			Def: gopkg.TypeStruct{
				Fields: []gopkg.DeclVar{
					{
						Name: "MyVal",
						Type: gopkg.TypeNamed{
							Name:   "SomeImportedType",
							Import: "github.com/myrepo",
						},
					},
					{
						Name: "MyOtherVal",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypeNamed{
								Name:   "SomeOtherImportedType",
								Import: "github.com/myotherrepo",
							},
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"github.com/myrepo":      "myrepo",
				"github.com/myotherrepo": "myotherrepo",
			},
			Expected: `struct {
	MyVal myrepo.SomeImportedType
	MyOtherVal *myotherrepo.SomeOtherImportedType
}`,
		},
		{
			Name: "struct with fields and embedded types",
			Def: gopkg.TypeStruct{
				Embeds: []gopkg.Type{
					gopkg.TypeInt32{},
					gopkg.TypeError{},
				},
				Fields: []gopkg.DeclVar{
					{
						Name: "MyVal",
						Type: gopkg.TypeNamed{
							Name:   "SomeImportedType",
							Import: "github.com/myrepo",
						},
					},
					{
						Name: "MyOtherVal",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypeNamed{
								Name:   "SomeOtherImportedType",
								Import: "github.com/myotherrepo",
							},
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"github.com/myrepo":      "myrepo",
				"github.com/myotherrepo": "myotherrepo",
			},
			Expected: `struct {
	int32
	error

	MyVal myrepo.SomeImportedType
	MyOtherVal *myotherrepo.SomeOtherImportedType
}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			fullType, err := test.Def.FullType(test.ImportAliases)
			require.NoError(t, err)
			require.Equal(t, test.Expected, fullType)
		})
	}
}

// TestTypeInterfaceFullType tests the TypeInterface.FullType functionality
// It is in it's own test fixutre as it's morecomplicated than the other
// types
func TestTypeInterfaceFullType(t *testing.T) {

	testCases := []struct {
		Name          string
		Def           gopkg.TypeInterface
		ImportAliases map[string]string
		Expected      string
	}{
		{
			Name:     "empty interface",
			Def:      gopkg.TypeInterface{},
			Expected: "interface{}",
		},
		{
			Name: "with functions no import aliases",
			Def: gopkg.TypeInterface{
				Funcs: []gopkg.DeclFunc{
					{
						Name: "FirstMethod",
						Args: []gopkg.DeclVar{
							{
								Name: "a",
								Type: gopkg.TypeInt32{},
							},
							{
								Name: "b",
								Type: gopkg.TypeFloat64{},
							},
						},
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeString{},
							gopkg.TypeError{},
						),
					},
					{
						Name: "SecondMethod",
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeNamed{
								Name: "MyType",
							},
						),
					},
					{
						Name: "VariadicMethod",
						Args: []gopkg.DeclVar{
							{
								Type: gopkg.TypeFloat64{},
							},
						},
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeString{},
							gopkg.TypeError{},
						),
						VariadicLastArg: true,
					},
				},
			},
			Expected: `interface {
	FirstMethod(a int32, b float64) (string, error)
	SecondMethod() MyType
	VariadicMethod(...float64) (string, error)
}`,
		},
		{
			Name: "with functions and import aliases",
			Def: gopkg.TypeInterface{
				Funcs: []gopkg.DeclFunc{
					{
						Name: "MyMethod",
						Args: []gopkg.DeclVar{
							{
								Name: "val",
								Type: gopkg.TypeNamed{
									Name:   "AStruct",
									Import: "some/path/toa",
								},
							},
						},
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeNamed{
								Name:   "BStruct",
								Import: "some/other/path/tob",
							},
							gopkg.TypeError{},
						),
					},
					{
						Name: "OtherMethod",
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeNamed{
								Name:   "CStruct",
								Import: "some/third/path/toc",
							},
						),
					},
				},
			},
			ImportAliases: map[string]string{
				"some/path/toa":       "alias_a",
				"some/other/path/tob": "alias_b",
				"some/third/path/toc": "alias_c",
			},
			Expected: `interface {
	MyMethod(val alias_a.AStruct) (alias_b.BStruct, error)
	OtherMethod() alias_c.CStruct
}`,
		},
		{
			Name: "with embedded types only",
			Def: gopkg.TypeInterface{
				Embeds: []gopkg.Type{
					gopkg.TypeString{},
					gopkg.TypeNamed{
						Name:   "SomeT",
						Import: "path/to/some",
					},
					gopkg.TypeNamed{
						Name:   "Other",
						Import: "path/to/other",
					},
				},
			},
			ImportAliases: map[string]string{
				"path/to/some":  "alias_some",
				"path/to/other": "other",
			},
			Expected: `interface {
	string
	alias_some.SomeT
	other.Other
}`,
		},
		{
			Name: "with embedded types and functions",
			Def: gopkg.TypeInterface{
				Embeds: []gopkg.Type{
					gopkg.TypeString{},
					gopkg.TypeBool{},
				},
				Funcs: []gopkg.DeclFunc{
					{
						Name: "One",
						Args: tmpl.UnnamedReturnArgs(
							gopkg.TypeInt64{},
							gopkg.TypeFloat64{},
						),
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeInt64{},
						),
					},
					{
						Name: "Two",
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeInt32{},
						),
					},
				},
			},
			ImportAliases: map[string]string{
				"path/to/some":  "alias_some",
				"path/to/other": "other",
			},
			Expected: `interface {
	string
	bool

	One(int64, float64) int64
	Two() int32
}`,
		},
		{
			Name: "functions with named returns and import aliases",
			Def: gopkg.TypeInterface{
				Funcs: []gopkg.DeclFunc{
					{
						Name: "One",
						ReturnArgs: []gopkg.DeclVar{
							{
								Name: "val",
								Type: gopkg.TypeNamed{
									Name:   "BStruct",
									Import: "some/other/path/tob",
								},
							},
							{
								Name: "err",
								Type: gopkg.TypeError{},
							},
						},
					},
					{
						Name: "Two",
						ReturnArgs: []gopkg.DeclVar{
							{
								Name: "secondVal",
								Type: gopkg.TypeNamed{
									Name:   "CStruct",
									Import: "some/third/path/toc",
								},
							},
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"some/path/toa":       "alias_a",
				"some/other/path/tob": "alias_b",
				"some/third/path/toc": "alias_c",
			},
			Expected: `interface {
	One() (val alias_b.BStruct, err error)
	Two() (secondVal alias_c.CStruct)
}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			fullType, err := test.Def.FullType(test.ImportAliases)
			require.NoError(t, err)
			require.Equal(t, test.Expected, fullType)
		})
	}
}

func TestTypeUnknownNamedDefaultInit(t *testing.T) {

	testCases := []struct {
		// TODO add names to these tests
		Def           gopkg.TypeNamed
		ImportAliases map[string]string
		Expected      string
		ExpectedErr   error
	}{
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeByte{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeError{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeFloat32{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeFloat64{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeInt{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeInt32{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeInt64{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeString{},
			},
			Expected: "\"\"",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeInterface{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeNameWithImport",
				Import:    "some/other/import",
				ValueType: gopkg.TypeInterface{},
			},
			ImportAliases: map[string]string{
				"some/other/import": "some_other_alias",
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeName",
				ValueType: gopkg.TypeStruct{},
			},
			Expected: "SomeName{}",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeOtherName",
				Import:    "github.com/somepkg",
				ValueType: gopkg.TypeStruct{},
			},
			ImportAliases: map[string]string{
				"github.com/somepkg": "somepkg",
			},
			Expected: "somepkg.SomeOtherName{}",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeType",
				ValueType: gopkg.TypeArray{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeType",
				ValueType: gopkg.TypePointer{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeNamed{
				Name:      "SomeType",
				ValueType: gopkg.TypeMap{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeNamed{
				Name:   "MyType",
				Import: "my/import/path",
			},
			ExpectedErr: errors.New("cannot deduce default init for named type with no value type"),
		},
	}

	for _, test := range testCases {
		testName, err := test.Def.FullType(nil)
		require.NoError(t, err)

		t.Run(testName, func(t *testing.T) {
			actual, err := test.Def.DefaultInit(test.ImportAliases)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.Expected, actual)
		})
	}
}

func TestTypeRequiredImports_SimpleTypes(t *testing.T) {

	simpleTypes := []gopkg.Type{
		gopkg.TypeBool{},
		gopkg.TypeByte{},
		gopkg.TypeError{},
		gopkg.TypeFloat32{},
		gopkg.TypeFloat64{},
		gopkg.TypeInt{},
		gopkg.TypeInt32{},
		gopkg.TypeInt64{},
		gopkg.TypeString{},
	}

	for _, simpleType := range simpleTypes {
		testName, err := simpleType.FullType(nil)
		require.NoError(t, err)

		t.Run(testName, func(t *testing.T) {
			require.Equal(
				t,
				map[string]bool(nil),
				simpleType.RequiredImports(),
			)
			return
		})
	}
}

func TestTypeRequiredImports_CompositeTypes(t *testing.T) {

	testCases := []struct {
		Name     string
		Def      gopkg.Type
		Expected map[string]bool
	}{
		{
			Name: "array of simple type",
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeString{},
			},
			Expected: map[string]bool(nil),
		},
		{
			Name: "array of composite type",
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeMap{
					KeyType: gopkg.TypeNamed{
						Import: "array/other/import",
					},
					ValueType: gopkg.TypeNamed{
						Import: "array/an/import",
					},
				},
			},
			Expected: map[string]bool{
				"array/an/import":    true,
				"array/other/import": true,
			},
		},
		{
			Name: "pointer of simple type",
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeString{},
			},
			Expected: map[string]bool(nil),
		},
		{
			Name: "pointer of composite type",
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeMap{
					KeyType: gopkg.TypeNamed{
						Import: "pointer/other/import",
					},
					ValueType: gopkg.TypeNamed{
						Import: "pointer/an/import",
					},
				},
			},
			Expected: map[string]bool{
				"pointer/an/import":    true,
				"pointer/other/import": true,
			},
		},
		{
			Name: "map with simple types",
			Def: gopkg.TypeMap{
				KeyType:   gopkg.TypeInt{},
				ValueType: gopkg.TypeString{},
			},
			Expected: map[string]bool(nil),
		},
		{
			Name: "map with named types",
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeNamed{
					Import: "other/import",
				},
				ValueType: gopkg.TypeNamed{
					Import: "an/import",
				},
			},
			Expected: map[string]bool{
				"an/import":    true,
				"other/import": true,
			},
		},
		{
			Name:     "named without import",
			Def:      gopkg.TypeNamed{},
			Expected: map[string]bool(nil),
		},
		{
			Name: "named with import",
			Def: gopkg.TypeNamed{
				Import: "some/import",
			},
			Expected: map[string]bool{
				"some/import": true,
			},
		},
		{
			Name:     "struct without fields",
			Def:      gopkg.TypeStruct{},
			Expected: map[string]bool(nil),
		},
		{
			Name: "struct with named fields",
			Def: gopkg.TypeStruct{
				Fields: []gopkg.DeclVar{
					{
						Type: gopkg.TypeNamed{
							Import: "import/a",
						},
					},
					{
						Type: gopkg.TypeString{},
					},
					{
						Type: gopkg.TypeNamed{
							Import: "import/b",
						},
					},
				},
			},
			Expected: map[string]bool{
				"import/a": true,
				"import/b": true,
			},
		},
		{
			Name: "struct with named embedded types",
			Def: gopkg.TypeStruct{
				Embeds: []gopkg.Type{
					gopkg.TypeNamed{
						Import: "import/a",
					},
					gopkg.TypeString{},
					gopkg.TypeNamed{
						Import: "import/b",
					},
				},
			},
			Expected: map[string]bool{
				"import/a": true,
				"import/b": true,
			},
		},
		{
			Name:     "interface without funcs",
			Def:      gopkg.TypeInterface{},
			Expected: map[string]bool(nil),
		},
		{
			Name: "interface with named embedded types",
			Def: gopkg.TypeInterface{
				Embeds: []gopkg.Type{
					gopkg.TypeNamed{
						Import: "import/aa",
					},
					gopkg.TypeString{},
					gopkg.TypeInt64{},
					gopkg.TypeNamed{
						Import: "import/bb",
					},
				},
			},
			Expected: map[string]bool{
				"import/aa": true,
				"import/bb": true,
			},
		},
		{
			Name: "interface with named funcs",
			Def: gopkg.TypeInterface{
				Funcs: []gopkg.DeclFunc{
					{
						Args: []gopkg.DeclVar{
							{
								Type: gopkg.TypeNamed{
									Import: "import/aa",
								},
							},
							{
								Type: gopkg.TypeNamed{
									Import: "import/bb",
								},
							},
						},
					},
					{
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeNamed{
								Import: "import/bb",
							},
							gopkg.TypeNamed{
								Import: "import/cc",
							},
						),
					},
				},
			},
			Expected: map[string]bool{
				"import/aa": true,
				"import/bb": true,
				"import/cc": true,
			},
		},
		{
			Name: "TypeFunc empty",
			Def:  gopkg.TypeFunc{},
		},
		{
			Name: "TypeFunc with built in type args and return args",
			Def: gopkg.TypeFunc{
				Args: tmpl.UnnamedReturnArgs(
					gopkg.TypeString{},
					gopkg.TypeInt32{},
				),
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeString{},
					gopkg.TypeInt32{},
				),
			},
		},
		{
			Name: "TypeFunc named types in type args and return args",
			Def: gopkg.TypeFunc{
				Args: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{
						Name:   "A",
						Import: "import/aaa",
					},
					gopkg.TypeNamed{
						Name:   "B",
						Import: "import/bbb",
					},
				),
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{
						Name:   "A",
						Import: "import/ccc",
					},
					gopkg.TypeNamed{
						Name:   "B",
						Import: "import/ddd",
					},
				),
			},
			Expected: map[string]bool{
				"import/aaa": true,
				"import/bbb": true,
				"import/ccc": true,
				"import/ddd": true,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			actual := test.Def.RequiredImports()

			if len(test.Expected) == 0 && len(actual) == 0 {
				// Allow actual to be either nil or an empty map
				return
			}

			require.Equal(t, test.Expected, test.Def.RequiredImports())
		})
	}

}
