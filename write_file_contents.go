package gopkg

import (
	"errors"
	"io"
)

func WriteFileContents(
	w io.Writer,
	c FileContents,
) error {

	if c.PackageName == "" {
		return errors.New("package name cannot be empty")
	}

	w.Write([]byte("package " + c.PackageName + "\n\n"))

	if len(c.Imports) > 0 {
		w.Write([]byte("import (\n"))
		for _, i := range c.Imports {
			w.Write([]byte("\t" + i.Alias + " \"" + i.Import + "\"\n"))
		}
		w.Write([]byte(")\n\n"))
	}

	importAliases := importsToImportAliasMap(c.Imports)

	err := WriteDeclVars(w, "const", c.Consts, importAliases)
	if err != nil {
		return err
	}

	err = WriteDeclVars(w, "var", c.Vars, importAliases)
	if err != nil {
		return err
	}

	for _, t := range c.Types {
		err := WriteDeclType(w, t, importAliases)
		if err != nil {
			return err
		}
		w.Write([]byte("\n"))
	}

	for _, f := range c.Functions {
		err := WriteDeclFunc(w, f, importAliases)
		if err != nil {
			return err
		}
		w.Write([]byte("\n"))
	}

	return nil
}

func importsToImportAliasMap(imports []ImportAndAlias) map[string]string {

	ret := make(map[string]string)
	for _, i := range imports {
		ret[i.Import] = i.Alias
	}
	return ret
}
