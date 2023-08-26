package main

import (
	"juicer/cmd"
	"log"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run chess server")
	}
}
