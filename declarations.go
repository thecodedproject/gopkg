package gopkg

import (
)

type DeclFunc struct {
	Name string
	Import string
	Receiver FuncReceiver
	Args []DeclVar
	ReturnArgs []Type
	BodyTmpl string
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
	Fields []DeclVar//map[string]Type
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
}

