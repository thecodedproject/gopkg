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

	fullType, err := decl.Type.FullType(importAliases)
	if err != nil {
		return err
	}

	w.Write([]byte(
		"type " + decl.Name + " " + fullType + "\n",
	))

	return nil
}
