package config

import (
	"fmt"
)

var (
	// ErrNotADirectory is returned when a given dir is not a directory.
	ErrNotADirectory = func(dir string) error {
		return fmt.Errorf("%s is not a directory", dir)
	}

	// ErrDirectoryNotEmpty is returned when a given dir is not empty.
	ErrDirectoryNotEmpty = func(dir string) error {
		return fmt.Errorf("%s is not empty", dir)
	}

	// ErrDirectoryIsEmpty is returned when a given dir is empty.
	ErrDirectoryIsEmpty = func(dir string) error {
		return fmt.Errorf("%s is empty", dir)
	}

	// ErrEmptyOutputDir is returned when the output directory is empty.
	ErrEmptyOutputDir = fmt.Errorf("the output directory is empty")

	// ErrEmptyCoreRepoDir is returned when the core repository directory is empty.
	ErrEmptyCoreRepoDir = fmt.Errorf("the core repository directory is empty")

	// ErrEmptyCoreRepoDir is returned when the core repository URL is empty.
	ErrEmptyCoreRepoURL = fmt.Errorf("the core repository URL is empty")

	// ErrEmptyCoreRepoBranch is returned when the core repository branch is empty.
	ErrEmptyCoreRepoBranch = fmt.Errorf("the core repository branch is empty")

	// ErrInvalidMaxWorkers is returned when the number of workers is invalid.
	ErrInvalidMaxWorkers = fmt.Errorf("invalid number of workers")
)
