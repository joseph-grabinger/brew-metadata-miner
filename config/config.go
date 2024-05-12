package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Config struct for the application.
type Config struct {
	// The directory where the extracted meta data will be stored.
	OutputDir string `yaml:"output_dir"`

	CoreRepo struct {
		// The URL of the core repository.
		URL string `yaml:"url"`

		// The branch of the core repository.
		Branch string `yaml:"branch"`

		// The path to the core repository.
		Dir string `yaml:"dir"`

		// A boolean value indicating whether the core repository should be cloned or not.
		Clone bool `yaml:"clone"`
	} `yaml:"core_repo"`

	Reader ReaderConfig `yaml:"reader"`
}

type ReaderConfig struct {
	// The maximum number of concurrent workers to use.
	MaxWorkers int `yaml:"max_workers"`

	// A boolean flag indicating whether the repo URL should be derived if no head is specified.
	DeriveRepo bool `yaml:"derive_repo"`

	// The license to use when no license is specified.
	FallbackLicense string `yaml:"fallback_license"`
}

// Print prints the configuration to the console.
func (c *Config) Print() {
	fmt.Printf("OutputDir: %s\n", c.OutputDir)
	fmt.Printf("CoreRepo.URL: %s\n", c.CoreRepo.URL)
	fmt.Printf("CoreRepo.Branch: %s\n", c.CoreRepo.Branch)
	fmt.Printf("CoreRepo.Dir: %s\n", c.CoreRepo.Dir)
	fmt.Printf("CoreRepo.Clone: %t\n", c.CoreRepo.Clone)
}

// Validate validates the configuration and creates directories if needed.
func (c *Config) Validate() error {
	// verify the output directory is not empty
	if c.OutputDir == "" {
		return ErrEmptyOutputDir
	}

	// check if the output directory exists
	s, err := os.Stat(c.OutputDir)
	if err != nil && os.IsNotExist(err) {
		// create the output directory
		err = os.MkdirAll(c.OutputDir, 0755)
		if err != nil {
			return err
		}
	} else if !s.IsDir() {
		return ErrNotADirectory(c.OutputDir)
	}

	// verify the output directory is empty
	if empty, err := isEmpty(c.OutputDir); err != nil {
		return err
	} else if !empty {
		return ErrDirectoryNotEmpty(c.OutputDir)
	}

	// verify the repository directory is not empty
	if c.CoreRepo.Dir == "" {
		return ErrEmptyCoreRepoDir
	}

	// check if the core repository directory exists
	s, err = os.Stat(c.CoreRepo.Dir)
	if err != nil && os.IsNotExist(err) && c.CoreRepo.Clone {
		// create the core repository directory
		err = os.MkdirAll(c.CoreRepo.Dir, 0755)
		if err != nil {
			return err
		}
	} else if !s.IsDir() {
		return ErrNotADirectory(c.CoreRepo.Dir)
	}

	// verify the core repository directory is not empty when the clone option is disabled
	if !c.CoreRepo.Clone {
		if empty, err := isEmpty(c.CoreRepo.Dir); err != nil {
			return err
		} else if empty {
			return ErrDirectoryIsEmpty(c.CoreRepo.Dir)
		}
	}

	// verify the core repository URL is not empty
	if c.CoreRepo.URL == "" {
		return ErrEmptyCoreRepoURL
	}

	// verify the core repository branch is not empty
	if c.CoreRepo.Branch == "" {
		return ErrEmptyCoreRepoBranch
	}

	// verify the number of workers is valid
	if c.Reader.MaxWorkers <= 0 {
		return ErrInvalidMaxWorkers
	}

	return nil
}

// NewConfig returns a new decoded Config struct from a given configPath.
func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

// isEmpty checks if the directory at the given path is empty.
func isEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
