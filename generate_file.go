package gopkg

import (
	"errors"
	"io"
)

func GenerateFileToWriter(
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

	for _, t := range c.Types {
		err := GenerateDeclType(w, t, importAliases)
		if err != nil {
			return err
		}
		w.Write([]byte("\n"))
	}

	for _, f := range c.Functions {
		err := GenerateDeclFunc(w, f, importAliases)
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
