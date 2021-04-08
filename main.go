package main

import (
	"github.com/satmaelstorm/gost/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Println(err)
	}
}
