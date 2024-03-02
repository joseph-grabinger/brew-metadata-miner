package main

import (
	"fmt"
	"log"
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

func main() {
	config, err := NewConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully parsed the configuration file:")
	config.Print()

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	if config.CoreRepo.Clone {
		// err := os.MkdirAll(config.CoreRepo.Dir, 0755)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// clone the core repository
		_, err = git.PlainClone(config.CoreRepo.Dir, false, &git.CloneOptions{
			URL: config.CoreRepo.URL,
			//ReferenceName: plumbing.ReferenceName(config.CoreRepo.Branch),
			Progress: os.Stdout,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Successfully cloned the core repository")
}
