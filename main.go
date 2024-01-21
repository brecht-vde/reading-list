package main

import (
	"log"
	"os"

	"github.com/brecht-vde/reading-list/cli"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("no verbs provided.")
		os.Exit(-1)
	}

	switch os.Args[1] {
	case "create-database":
		cmd, err := cli.NewCreateDateBaseCommand(os.Args[2:])

		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}

		err = cmd.Run()

		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	case "sync":

	default:
		log.Fatal("unknown verb provided")
		os.Exit(-1)
	}
}
