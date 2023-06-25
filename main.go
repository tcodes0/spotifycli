package main

import (
	"log"

	"github.com/masroorhasan/spotifycli/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatal("execute: ", err)
	}
}
