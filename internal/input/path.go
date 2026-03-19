package input

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FilePaths struct {
	InputPath   string
	InputDir    string
	InputBase   string
	OutputPath  string
	DisplayPath string
}

func Resolve(input string) (FilePaths, error) {
	if strings.TrimSpace(input) == "" {
		return FilePaths{}, fmt.Errorf("%w: input path is empty", cliPathError{})
	}
	if filepath.Ext(input) != ".md" {
		return FilePaths{}, fmt.Errorf("%w: input must have a .md extension: %s", cliPathError{}, input)
	}

	absInput, err := filepath.Abs(input)
	if err != nil {
		return FilePaths{}, fmt.Errorf("%w: resolve input path: %v", cliPathError{}, err)
	}

	info, err := os.Stat(absInput)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return FilePaths{}, fmt.Errorf("%w: input file does not exist: %s", cliPathError{}, input)
		}
		return FilePaths{}, fmt.Errorf("%w: read input file info: %v", cliPathError{}, err)
	}
	if info.IsDir() {
		return FilePaths{}, fmt.Errorf("%w: input path is a directory: %s", cliPathError{}, input)
	}

	inputDir := filepath.Dir(absInput)
	inputBase := strings.TrimSuffix(filepath.Base(absInput), filepath.Ext(absInput))
	outputPath, err := NextOutputPath(inputDir, inputBase)
	if err != nil {
		return FilePaths{}, fmt.Errorf("%w: %v", cliPathError{}, err)
	}

	return FilePaths{
		InputPath:   absInput,
		InputDir:    inputDir,
		InputBase:   inputBase,
		OutputPath:  outputPath,
		DisplayPath: outputPath,
	}, nil
}

func NextOutputPath(dir, base string) (string, error) {
	if strings.TrimSpace(dir) == "" {
		return "", errors.New("output directory is empty")
	}
	if strings.TrimSpace(base) == "" {
		return "", errors.New("output base name is empty")
	}

	primary := filepath.Join(dir, base+".pdf")
	if _, err := os.Stat(primary); errors.Is(err, os.ErrNotExist) {
		return primary, nil
	} else if err != nil {
		return "", fmt.Errorf("check existing output %s: %w", primary, err)
	}

	for i := 1; i < 1000; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s (%02d).pdf", base, i))
		if _, err := os.Stat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate, nil
		} else if err != nil {
			return "", fmt.Errorf("check existing output %s: %w", candidate, err)
		}
	}

	return "", fmt.Errorf("could not determine an available output filename for %s", base)
}

type cliPathError struct{}

func (cliPathError) Error() string { return "path error" }
