package internal

import (
	"flag"
	"github.com/thecodedproject/gopkg"
	//tmpl "github.com/thecodedproject/gopkg/tmpl"
	"os"
)

var writeInplace = flag.Bool("w", false, "write result to (source) file instead of stdout")

func Generate() error {

	flag.Parse()

	if flag.NArg() != 1 {
		return nil
	}

	// TODO: Allow parsing files recursively with `./...`
	path := flag.Arg(0)

	pkgFiles, err := gopkg.Parse(path)
	if err != nil {
		return err
	}

	if *writeInplace {
		return gopkg.LintAndGenerate(pkgFiles)
	}

	for _, f := range pkgFiles {

		// TODO: run linting here

		err := gopkg.WriteFileContents(os.Stdout, f)
		if err != nil {
			return err
		}
	}
	return nil
}

