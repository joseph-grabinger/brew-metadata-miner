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
	// check if the output directory exists
	s, err := os.Stat(c.OutputDir)
	if err != nil {
		if os.IsNotExist(err) {
			// create the output directory
			err = os.Mkdir(c.OutputDir, 0755)
			if err != nil {
				return err
			}
		}
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("%s is not a directory", c.OutputDir)
	}

	// verify the output directory is empty
	if empty, err := isEmpty(c.OutputDir); err != nil {
		return err
	} else if !empty {
		return fmt.Errorf("%s is not empty", c.OutputDir)
	}

	// check if the core repository directory exists
	s, err = os.Stat(c.CoreRepo.Dir)
	if err != nil {
		if os.IsNotExist(err) && c.CoreRepo.Clone {
			// create the core repository directory
			err = os.MkdirAll(c.CoreRepo.Dir, 0755)
			if err != nil {
				return err
			}
		}
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("%s is not a directory", c.CoreRepo.Dir)
	}

	// verify the core repository directory is not empty when the clone option is disabled
	if !c.CoreRepo.Clone {
		if empty, err := isEmpty(c.CoreRepo.Dir); err != nil {
			return err
		} else if empty {
			return fmt.Errorf("%s is empty", c.CoreRepo.Dir)
		}
	}

	// verify the core repository URL is not empty
	if c.CoreRepo.URL == "" {
		return fmt.Errorf("the core repository URL is empty")
	}

	// verify the core repository branch is not empty
	if c.CoreRepo.Branch == "" {
		return fmt.Errorf("the core repository branch is empty")
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