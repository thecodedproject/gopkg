package gopkg

import (
	"errors"
	"io"
)

func WriteImports(
	w io.Writer,
	imports []ImportAndAlias,
) error {

	if len(imports) == 0 {
		return nil
	}

	if len(imports) == 1  && imports[0].Alias == ""{
		w.Write([]byte("import \"" + imports[0].Import + "\"\n\n"))
		return nil
	}


	w.Write([]byte("import (\n"))
	lastImportGroup := imports[0].Group
	for _, i := range imports {

		if i.Group < lastImportGroup {
			return errors.New("WriteImports: import groups are not in order")
		}

		if i.Group > lastImportGroup {
			w.Write([]byte("\n"))
		}

		w.Write([]byte("\t"))

		if i.Alias != "" {
			w.Write([]byte(i.Alias + " "))
		}

		w.Write([]byte("\"" + i.Import + "\"\n"))

		lastImportGroup = i.Group
	}
	w.Write([]byte(")\n\n"))

	return nil
}
