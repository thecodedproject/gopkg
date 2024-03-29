package gopkg

import (
	"errors"
	"io"
)

func WriteDeclVars(
	w io.Writer,
	keyword string,
	decls []DeclVar,
	importAliases map[string]string,
) error {

	if len(decls) == 0 {
		return nil
	}

	for _, d := range decls {
		if d.Name == "" {
			return errors.New("WriteDeclVar: DeclVar.Name cannot be empty")
		}

		if d.Type == nil && d.LiteralValue == "" {
			return errors.New("WriteDeclVar: one of DeclVar.Type and DeclVar.LiteralValue must be set")
		}
	}

	if len(decls) == 1 {
		if decls[0].DocString != "" {
			w.Write([]byte(decls[0].DocString + "\n"))
		}
		w.Write([]byte(keyword + " "))
		writeDeclVar(w, decls[0], importAliases)
		w.Write([]byte("\n"))
		return nil
	}

	w.Write([]byte(keyword + " (\n"))
	for i, d := range decls {
		if d.DocString != "" {
			// vanity space var/const declarations if there is a docstring
			if i != 0 {
				w.Write([]byte("\n"))
			}
			w.Write([]byte("\t" + d.DocString + "\n"))
		}
		w.Write([]byte("\t"))
		writeDeclVar(w, d, importAliases)
	}

	w.Write([]byte(")\n\n"))

	return nil
}

func writeDeclVar(
	w io.Writer,
	d DeclVar,
	importAliases map[string]string,
) error {

	w.Write([]byte(d.Name))

	if d.Type != nil {
		if _, isLiteral := d.Type.(TypeUnnamedLiteral); !isLiteral {
			dType, err := d.Type.FullType(importAliases)
			if err != nil {
				return err
			}

			w.Write([]byte(" " + dType))
		}
	}

	if d.LiteralValue != "" {
		w.Write([]byte(" = " + d.LiteralValue))
	}

	w.Write([]byte("\n"))

	return nil
}
