package gopkg

import (
	path "path/filepath"
	"os"
)

// CreatePathAndOpen creats all contains directors in filepath if they do not
// exist and opens a file for writing at filepath.
//
// If filepath already exists and is a file then the file will be overwritten.
// If filepath already exists and is a directory then an error will be returned.
func CreatePathAndOpen(
	filepath string,
) (*os.File, error) {

	dir, _ := path.Split(filepath)

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	return os.Create(filepath)
}
