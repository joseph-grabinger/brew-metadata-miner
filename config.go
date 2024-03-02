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

func (c *Config) Print() {
	fmt.Printf("OutputDir: %s\n", c.OutputDir)
	fmt.Printf("CoreRepo.URL: %s\n", c.CoreRepo.URL)
	fmt.Printf("CoreRepo.Branch: %s\n", c.CoreRepo.Branch)
	fmt.Printf("CoreRepo.Dir: %s\n", c.CoreRepo.Dir)
	fmt.Printf("CoreRepo.Clone: %t\n", c.CoreRepo.Clone)
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
