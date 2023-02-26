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

func LintAndGenerate(
	files []FileContents,
	extraLintRules ...func([]FileContents)error,
) error {

	err := Lint(files, extraLintRules...)
	if err != nil {
		return err
	}

	err = Generate(files)
	if err != nil {
		return err
	}

	return nil
}
