package miner

import (
	"main/config"
	"main/miner/reader"
	"main/miner/types"
	"main/miner/writer"
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
	return writer.WriteFormulae(m.config.OutputDir, m.formulae)
}
