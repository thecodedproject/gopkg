package gopkg

import (
	"errors"
	"io"
)

func GenerateFileToWriter(
	w io.Writer,
	c Contents,
) error {

	if c.PackageName == "" {
		return errors.New("package name cannot be empty")
	}

	w.Write([]byte("package " + c.PackageName + "\n\n"))

	if len(c.ImportsAndAliases) > 0 {
		w.Write([]byte("import (\n"))
		for importPath, alias := range c.ImportsAndAliases {
			w.Write([]byte("\t" + alias + " \"" + importPath + "\"\n"))
		}
		w.Write([]byte(")\n\n"))
	}

	for _, t := range c.Types {
		err := GenerateDeclType(w, t, c.ImportsAndAliases)
		if err != nil {
			return err
		}
		w.Write([]byte("\n"))
	}

	for _, f := range c.Functions {
		err := GenerateFunc(w, f, c.ImportsAndAliases)
		if err != nil {
			return err
		}
		w.Write([]byte("\n"))
	}

	return nil
}
