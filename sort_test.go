package gopkg_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
)

func TestSortFuncs(t *testing.T) {

	testCases := []struct{
		Name string
		Funcs []gopkg.DeclFunc
		Expected []gopkg.DeclFunc
	}{
		{
			Name: "empty list",
		},
		{
			Name: "normal funcs are sorted by name with exported funcs coming first",
			Funcs: []gopkg.DeclFunc{
				{Name: "cFunc"},
				{Name: "CFunc"},
				{Name: "AFunc"},
				{Name: "bFunc"},
				{Name: "BFunc"},
				{Name: "aFunc"},
			},
			Expected: []gopkg.DeclFunc{
				{Name: "AFunc"},
				{Name: "BFunc"},
				{Name: "CFunc"},
				{Name: "aFunc"},
				{Name: "bFunc"},
				{Name: "cFunc"},
			},
		},
		{
			Name: "receiver funs are sorted by receiver types then func name - irrespective of pointer receivers",
			Funcs: []gopkg.DeclFunc{
				{
					Name: "CFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "BType",
					},
				},
				{
					Name: "aFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "AType",
					},
				},
				{
					Name: "AFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "AType",
						IsPointer: true,
					},
				},
				{
					Name: "BFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "BType",
						IsPointer: true,
					},
				},
			},
			Expected: []gopkg.DeclFunc{
				{
					Name: "AFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "AType",
						IsPointer: true,
					},
				},
				{
					Name: "aFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "AType",
					},
				},
				{
					Name: "BFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "BType",
						IsPointer: true,
					},
				},
				{
					Name: "CFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "BType",
					},
				},
			},
		},
		{
			Name: "receiver funs come before normal funcs",
			Funcs: []gopkg.DeclFunc{
				{Name: "anotherFunc"},
				{Name: "AFunc"},
				{
					Name: "bFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "AType",
						IsPointer: true,
					},
				},
				{
					Name: "BFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "AType",
					},
				},
			},
			Expected: []gopkg.DeclFunc{
				{
					Name: "BFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "AType",
					},
				},
				{
					Name: "bFunc",
					Receiver: gopkg.FuncReceiver{
						TypeName: "AType",
						IsPointer: true,
					},
				},
				{Name: "AFunc"},
				{Name: "anotherFunc"},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			gopkg.SortFuncs(test.Funcs)
			require.Equal(t, test.Expected, test.Funcs)
		})
	}
}
