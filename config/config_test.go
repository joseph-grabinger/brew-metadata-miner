package config

import (
	"errors"
	"os"
	"testing"
)

func TestValidate_EmptyOutputDir(t *testing.T) {
	c := &Config{
		OutputDir: "",
	}
	err := c.Validate()
	if !errors.Is(err, ErrEmptyOutputDir) {
		t.Error("expected an error, got nil")
	}
}

func TestValidate_OutputDirIsNotADirectory(t *testing.T) {
	c := &Config{
		OutputDir: "config.go",
	}
	err := c.Validate()
	if err.Error() != ErrNotADirectory(c.OutputDir).Error() {
		t.Error("expected an ErrNotADirectory, got: ", err)
	}
}

func TestValidate_OutputDirIsNotEmpty(t *testing.T) {
	c := &Config{
		OutputDir: "../config",
	}
	err := c.Validate()
	if err.Error() != ErrDirectoryNotEmpty(c.OutputDir).Error() {
		t.Error("expected an ErrDirectoryNotEmpty, got: ", err)
	}
}

func TestValidate_EmptyCoreRepoDir(t *testing.T) {
	c := &Config{
		OutputDir: "./test_dir",
	}

	// clean up
	defer os.RemoveAll(c.OutputDir)

	err := c.Validate()
	if !errors.Is(err, ErrEmptyCoreRepoDir) {
		t.Error("expected an ErrEmptyCoreRepoURL, got:", err)
	}
}

func TestValidate_CoreRepoDirIsNotADirectory(t *testing.T) {
	c := &Config{
		OutputDir: "./test_dir",
	}
	c.CoreRepo.Dir = "config.go"

	// clean up
	defer os.RemoveAll(c.OutputDir)

	err := c.Validate()
	if err.Error() != ErrNotADirectory(c.CoreRepo.Dir).Error() {
		t.Error("expected an ErrNotADirectory, got: ", err)
	}
}

func TestValidate_CoreRepoURLIsEmpty(t *testing.T) {
	c := &Config{
		OutputDir: "./test_dir",
	}
	c.CoreRepo.Dir = "../config"
	c.CoreRepo.Clone = false

	// clean up
	defer os.RemoveAll(c.OutputDir)

	err := c.Validate()
	if !errors.Is(err, ErrEmptyCoreRepoURL) {
		t.Error("expected an ErrEmptyCoreRepoURL, got: ", err)
	}
}

func TestValidate_CoreRepoDirIsEmpty(t *testing.T) {
	c := &Config{
		OutputDir: "./test_dir",
	}
	c.CoreRepo.Dir = c.OutputDir
	c.CoreRepo.Clone = false
	c.CoreRepo.URL = "https://github.com/Homebrew/homebrew-core.git"

	// clean up
	defer os.RemoveAll(c.OutputDir)

	err := c.Validate()
	if err.Error() != ErrDirectoryIsEmpty(c.CoreRepo.Dir).Error() {
		t.Error("expected an ErrDirectoryIsEmpty, got: ", err)
	}
}

func TestValidate_CoreRepoBranchIsEmpty(t *testing.T) {
	c := &Config{
		OutputDir: "./test_dir",
	}
	c.CoreRepo.Dir = c.OutputDir
	c.CoreRepo.Clone = true
	c.CoreRepo.URL = "https://github.com/Homebrew/homebrew-core.git"

	// clean up
	defer os.RemoveAll(c.OutputDir)

	err := c.Validate()
	if !errors.Is(err, ErrEmptyCoreRepoBranch) {
		t.Error("expected an ErrEmptyCoreRepoBranch, got: ", err)
	}
}
