package main

import (
	"fmt"
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
	if s, err := os.Stat(c.OutputDir); err != nil || s.IsDir() {
		// create the output directory
		err = os.Mkdir(c.OutputDir, 0755)
		if err != nil {
			return err
		}
	}

	// TODO: add validation of for the core repository directory dependent on clone flag

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
