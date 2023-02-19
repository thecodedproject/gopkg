package gopkg

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// PackageImportPath returns the import path of a package given its path in the
// filesystem
//
// `path` may be either a relative or absoute path to the package directory or a
// file withing the package directory.
//
// **Note 1:** it assumes that the package is within a go module - if it is not it
// will return an error.
//
// **Note 2:** there are no checks that the given package does contain `.go`
// files (i.e. is a go package)
//
// **Note 3:** only tested on Unix... Windows users tread carefully.
//
// Deets: This method traverses up the directory tree from `path`, looking for
// a `go.mod` file in each directory.
// Upon finding one, it will return the module path within `go.mod` appended with
// the relative path of the given `path`.
func PackageImportPath(path string) (string, error) {

	if path == "" {
		return "", errors.New("cannot get import path for empty path")
	}

	fileInfo, err := os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return "", errors.New("no such path `" + path + "`")
	} else if err != nil {
		return "", err
	}

	if !fileInfo.IsDir() {
		path = filepath.Dir(path)
	}


	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	var goModFile *os.File
	goModSearchPath := absPath
	var subPackagePath string
	for {
		if len(goModSearchPath) < 2 {
			return "", errors.New("path `" + path + "` not within a go module")
		}

		goModFile, err = os.Open(filepath.Join(goModSearchPath, "go.mod"))
		if err != nil && errors.Is(err, os.ErrNotExist) {

			goModFile.Close()

			subPackagePath = filepath.Join(
				filepath.Base(goModSearchPath),
				subPackagePath,
			)

			goModSearchPath = filepath.Dir(goModSearchPath)
			continue
		} else if err != nil {
			return "", err
		}
		break
	}

	scanner := bufio.NewScanner(goModFile)
	scanner.Scan()
	goModFile.Close()
	line := scanner.Text()

	splitLine := strings.Split(line, " ")

	if len(splitLine) != 2 {
		return "", errors.New(
			"unexpected first line of go.mod - expected `module <path>` got " + line)
	}

	if splitLine[0] != "module" {
		return "", errors.New(
			"unexpected first line of go.mod - expected `module <path>` got " + line)
	}

	importPath := filepath.Join(splitLine[1], subPackagePath)

	return importPath, nil
}
