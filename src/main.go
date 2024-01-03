package main

import (
	"log"
	"os"

	"github.com/brecht-vde/reading-list/src/cli"
)

func main() {
	verb := os.Args[1]

	switch verb {
	case "run":
		cli.Run(os.Args[2:])
	default:
		log.Fatalf(`Invalid verb '%v'.`, verb)
	}
}
