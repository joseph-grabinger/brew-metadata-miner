package main

import (
	"fmt"
	"log"
)

func main() {
	config, err := NewConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully parsed the configuration file:")
	config.Print()

}
