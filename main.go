package main

import (
	"embed"
	"log"

	"github.com/dankobg/juicer/cmd"
)

//go:embed public/*
var publicFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

func main() {
	if err := cmd.Run(publicFiles, templateFiles); err != nil {
		log.Fatalf("failed to run juicer chess server, %v", err)
	}
}
