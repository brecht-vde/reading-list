package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func Read(path string, urls chan<- string, errs chan<- error) {
	file, err := os.Open(path)

	if err != nil {
		errs <- err
		return
	}

	defer file.Close()

	read(file, urls, errs)
}

func read(reader io.Reader, urls chan<- string, errs chan<- error) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		url := scanner.Text()
		urls <- url
		log.Printf("sent url: %v\n", url)
	}

	if err := scanner.Err(); err != nil {
		errs <- err
		return
	}
}
