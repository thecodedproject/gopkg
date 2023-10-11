package main

import (
	internal "github.com/thecodedproject/gopkg/example_generators/lint/internal"
	log "log"
)

func main() {

	err := internal.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}
}

