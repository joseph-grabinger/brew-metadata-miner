package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"main/config"
	"main/parser"

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
		// clone the core repository
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

	parser := parser.NewParser(config)

	start := time.Now()
	if err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Duration:", time.Since(start))

	if err := parser.Pipe(); err != nil {
		log.Fatal(err)
	}

	// parser.Analyze()
}
