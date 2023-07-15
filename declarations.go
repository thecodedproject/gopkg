package gopkg

import (
	"reflect"
)

type DeclFunc struct {
	Name string
	Import string
	Receiver FuncReceiver
	Args []DeclVar
	// TODO make ReturnArgs a []DeclVar type and allow parsing/generating named return args
	ReturnArgs []DeclVar
	BodyTmpl string
	BodyData any
}

type FuncReceiver struct {
	VarName string
	TypeName string
	IsPointer bool
}

type DeclType struct {
	Name string
	Import string
	Type Type
}

type DeclVar struct {
	Type
	Name string
	Import string

	// LiteralValue is the value of the literal assigned to this variable declaration
	// (if one was assigned - otherwise it will be empty)
	//
	// e.g. for `var MyVar int = 123`, LiteralValue will be `123`
	LiteralValue string

	// StructTag holds the tags for a struct field if this DeclVar represents a
	// field within a struct.
	// If this DeclVar is not within a struct then it is not used.
	StructTag reflect.StructTag
}

func (d DeclFunc) RequiredImports() map[string]bool {

	ret := make(map[string]bool)
	for _, arg := range d.Args {
		ret = union(ret, arg.RequiredImports())
	}
	for _, retArg := range d.ReturnArgs {
		ret = union(ret, retArg.RequiredImports())
	}
	return ret
}
