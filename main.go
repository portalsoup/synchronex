package main

import (
	"log"
	"synchronex/command"
)

func main() {
	err := command.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
