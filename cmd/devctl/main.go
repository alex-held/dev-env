package main

import (
	"context"
	"log"
	"os"
	"os/exec"
)

func main() {
	ctx := context.Background()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := exec.Run(ctx, pwd, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
