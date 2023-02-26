package main_test

import (
	"testing"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {

	testDirs := []string{
		"example_single_enum",
	}

	for _, testDir := range testDirs {
		t.Run(testDir, func(t *testing.T) {

			generatedFiles := runGenerateAndGetGeneratedFileBuffers(t, testDir)

			runGoTestAndCheckOutput(t, testDir)

			checkGeneratedFilePaths(t, generatedFiles)

			checkGeneratedFileBuffers(t, testDir, generatedFiles)

			removeGeneratedFiles(t, generatedFiles)
		})
	}
}

func runGenerateAndGetGeneratedFileBuffers(
	t *testing.T,
	testDir string,
) map[string][]byte {

	initalFiles, err := listFilesRecursively(testDir)
	require.NoError(t, err)

	cmd := exec.Command("go", "generate", "./" + testDir)

	_, err = cmd.Output()
	// TODO output cmd stderr output if fails
	require.NoError(t, err, "go generate fail")

	postGenFiles, err := listFilesRecursively(testDir)
	require.NoError(t, err)

	generatedFiles := make(map[string][]byte, 0)
	for f := range postGenFiles {
		if !initalFiles[f] {
			fileBuffer, err := os.ReadFile(f)
			require.NoError(t, err)

			generatedFiles[f] = fileBuffer
		}
	}

	return generatedFiles
}

func runGoTestAndCheckOutput(
	t *testing.T,
	testDir string,
) {

	testCmd := exec.Command("go", "test", "-v", "-count=1", "./" + testDir)

	testOutput, err := testCmd.Output()
	require.NoError(t, err)

	timeRegex1, err := regexp.Compile(`\([0-9]\.[0-9][0-9]s\)`)
	require.NoError(t, err)
	timeRegex2, err := regexp.Compile(`[0-9]\.[0-9][0-9][0-9]s`)
	require.NoError(t, err)

	testOutput = timeRegex1.ReplaceAll(testOutput, []byte("(X.XXs)"))
	testOutput = timeRegex2.ReplaceAll(testOutput, []byte("X.XXXs"))

	t.Run("go_test", func(t *testing.T) {
		g := goldie.New(t)
		g.Assert(t, t.Name(), testOutput)
	})
}

func checkGeneratedFilePaths(
	t *testing.T,
	generatedFiles map[string][]byte,
) {

	var generatedFilesPaths []string
	for path := range generatedFiles {
		generatedFilesPaths = append(generatedFilesPaths, path)
	}

	sort.Slice(generatedFilesPaths, func(i, j int) bool {
		return generatedFilesPaths[i] < generatedFilesPaths[j]
	})

	var generatedFilesBuffer string
	for _, f := range generatedFilesPaths {
		generatedFilesBuffer += f + "\n"
	}

	t.Run("generated_file_paths", func(t *testing.T) {
			g := goldie.New(t)
			g.Assert(t, t.Name(), []byte(generatedFilesBuffer))
	})
}

func checkGeneratedFileBuffers(
	t *testing.T,
	testDir string,
	generatedFiles map[string][]byte,
) {

	for filePath, fileBuffer := range generatedFiles {

		testName, err := filepath.Rel(testDir, filePath)
		require.NoError(t, err)

		t.Run(testName, func(t *testing.T) {
			g := goldie.New(t)
			g.Assert(t, t.Name(), fileBuffer)
		})
	}
}

func removeGeneratedFiles(
	t *testing.T,
	generatedFiles map[string][]byte,
) {

	for path := range generatedFiles {
		err := os.Remove(path)
		require.NoError(t, err)
	}
}

func listFilesRecursively(path string) (map[string]bool, error) {

	files := make(map[string]bool)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			files[path] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
