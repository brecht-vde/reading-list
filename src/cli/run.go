package cli

import (
	"flag"
	"log"
	"os"

	"github.com/brecht-vde/reading-list/src/parsers"
)

func Run(args []string) {
	configuration := flag.String("c", "", "-c <path>/<to>/.config")

	error := flag.CommandLine.Parse(args)

	if error != nil {
		log.Fatalf(`Could not parse arguments %v`, args)
	}

	if _, error := os.Stat(*configuration); os.IsNotExist(error) {
		log.Fatalf(`The file does not exist '%v'`, *configuration)
	}

	parsers.Parse(configuration)
}
