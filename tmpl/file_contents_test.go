package tmpl_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestAppendFileContents(t *testing.T) {

	testCases := []struct{
		Name string
		F []gopkg.FileContents
		Methods []func() ([]gopkg.FileContents, error)
		Expected []gopkg.FileContents
		ExpectedErr error
	}{
		{
			Name: "Empty input",
		},
		{
			Name: "No file methods leaves files unchanged",
			F: []gopkg.FileContents{
				{Filepath: "a"},
				{Filepath: "b"},
				{Filepath: "c"},
			},
			Expected: []gopkg.FileContents{
				{Filepath: "a"},
				{Filepath: "b"},
				{Filepath: "c"},
			},
		},
		{
			Name: "When methods return no errors appends new file contents",
			F: []gopkg.FileContents{
				{Filepath: "a"},
				{Filepath: "b"},
				{Filepath: "c"},
			},
			Methods: []func()([]gopkg.FileContents, error){
				func()([]gopkg.FileContents, error){
					return []gopkg.FileContents{
						{Filepath: "d"},
						{Filepath: "e"},
					}, nil
				},
				func()([]gopkg.FileContents, error){
					return []gopkg.FileContents{
					}, nil
				},
				func()([]gopkg.FileContents, error){
					return []gopkg.FileContents{
						{Filepath: "f"},
					}, nil
				},
			},
			Expected: []gopkg.FileContents{
				{Filepath: "a"},
				{Filepath: "b"},
				{Filepath: "c"},
				{Filepath: "d"},
				{Filepath: "e"},
				{Filepath: "f"},
			},
		},
		{
			Name: "When one method returns error the error is returned",
			F: []gopkg.FileContents{
				{Filepath: "a"},
				{Filepath: "b"},
				{Filepath: "c"},
			},
			Methods: []func()([]gopkg.FileContents, error){
				func()([]gopkg.FileContents, error){
					return []gopkg.FileContents{
						{Filepath: "d"},
						{Filepath: "e"},
					}, nil
				},
				func()([]gopkg.FileContents, error){
					return []gopkg.FileContents{
					}, errors.New("someerror")
				},
			},
			ExpectedErr: errors.New("someerror"),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			actual, err := tmpl.AppendFileContents(test.F, test.Methods...)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, test.Expected, actual)
		})
	}
}
