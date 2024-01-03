package parsers

import (
	"bufio"
	"log"
	"os"
)

func Parse(path *string) {
	file, error := os.Open(*path)

	if error != nil {
		log.Fatalf(`Could not read file '%v'. Error: '%v'`, *path, error)

	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
	}

	if error = scanner.Err(); error != nil {
		log.Fatalf(`Could not parse file '%v'. Error: '%v;`, *path, error)
	}
}
