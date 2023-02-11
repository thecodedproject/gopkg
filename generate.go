package gopkg

import (
	"os"
)

func Generate(files []FileContents) error {

	for _, file := range files {
		writer, err := os.Create(file.Filepath)
		if err != nil {
			return err
		}
		err = WriteFileContents(writer, file)
		if err != nil {
			return err
		}
	}

	return nil
}
