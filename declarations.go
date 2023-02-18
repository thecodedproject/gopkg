package gopkg

import (
)

type DeclFunc struct {
	Name string
	Import string
	Receiver FuncReceiver
	Args []DeclVar
	// TODO make ReturnArgs a []DeclVar type and allow parsing/generating named return args
	ReturnArgs []Type
	BodyTmpl string
	BodyData any
}

type FuncReceiver struct {
	VarName string
	TypeName string
	IsPointer bool
}

type DeclStruct struct {
	Name string
	Import string
	// TODO: Maybe also add struct field descriptors
	Fields []DeclVar
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
}

