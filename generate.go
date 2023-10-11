package gopkg

import (
	"errors"
)

func Generate(files []FileContents) error {

	for _, file := range files {

		if file.Filepath == "" {
			return errors.New("gopkg.Generate: empty FileContents.Filepath - this is required")
		}

		writer, err := CreatePathAndOpen(file.Filepath)
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
	extraLintRules ...func([]FileContents) error,
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
