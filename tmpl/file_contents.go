package tmpl

import (
	"github.com/thecodedproject/gopkg"
)

// AppendFileContents will append new `gopkg.FileContents` returned from methods which may return an error
//
// It is a convenience function which makes it easier to build up a list of all
// the required FileContents objects for a package by breaking up their
// construction into multiple methods, while still allowing proper error handling.
func AppendFileContents(
	files []gopkg.FileContents,
	fileFuncs ...func() ([]gopkg.FileContents, error),
) ([]gopkg.FileContents, error) {

	for _, fileFunc := range fileFuncs {
		newFiles, err := fileFunc()
		if err != nil {
			return nil, err
		}

		files = append(files, newFiles...)
	}
	return files, nil
}
