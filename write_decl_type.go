package gopkg

import (
	"errors"
	"io"
)

func WriteDeclType(
	w io.Writer,
	decl DeclType,
	importAliases map[string]string,
) error {

	if decl.Name == "" {
		return errors.New("type decl name cannot be empty")
	}
	if decl.Type == nil {
		return errors.New("type decl type cannot be nil")
	}

	w.Write([]byte(
		"type " + decl.Name + " " + decl.Type.FullType(importAliases) + "\n",
	))

	return nil
}
