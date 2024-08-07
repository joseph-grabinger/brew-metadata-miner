package main

import (
	"fmt"
	"log"
	"os"

	"main/config"
	"main/miner"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	config, err := config.NewConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully parsed the configuration file:")
	config.Print()

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully validated the configuration")

	if config.CoreRepo.Clone {
		// Clone the core repository.
		_, err = git.PlainClone(config.CoreRepo.Dir, false, &git.CloneOptions{
			URL:           config.CoreRepo.URL,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + config.CoreRepo.Branch),
			Progress:      os.Stdout,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Successfully cloned the core repository")
	}

	miner := miner.NewMiner(config)

	if err := miner.ReadFormulae(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully parsed all formulae from the core repository")

	if err := miner.WriteFormulae(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully piped all formulae to the output file")

}
