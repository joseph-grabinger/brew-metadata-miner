package miner

import (
	"bufio"
	"os"
	"path/filepath"

	"main/config"
	"main/miner/reader"
	"main/miner/types"
)

type miner struct {
	config *config.Config

	// A map of formulae, where the key is the name of the formula.
	formulae map[string]*types.Formula
}

// NewMiner creates a new parser.
func NewMiner(config *config.Config) *miner {
	return &miner{
		config:   config,
		formulae: make(map[string]*types.Formula),
	}
}

// ReadFormaulas reads all formulae from the core repository into the formulas map.
func (m *miner) ReadFormulae() error {
	f, err := reader.ReadFormulae(m.config.CoreRepo.Dir, m.config.MaxWorkers)
	m.formulae = f
	return err
}

// WriteFormulae writes the formulae to the output file.
func (m *miner) WriteFormulae() error {
	path := filepath.Join(m.config.OutputDir, "deps-brew.tsv")
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	// Close file on function exit and check its' returned error.
	defer func() error {
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, formula := range m.formulae {
		// Write package line.
		line := formula.FormatPackageLine()
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}

		// Write dependency lines.
		for _, dep := range formula.Dependencies {
			f := m.formulae[dep.Name]
			if f == nil {
				panic("Dependency " + dep.Name + " not found in formula " + formula.Name)
			}

			line := f.FormatDependencyLine(dep)
			_, err := writer.WriteString(line)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
