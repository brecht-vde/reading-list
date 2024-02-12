package main

import (
	"log"
	"os"
	"sync"
)

func main() {
	arg, err := ParseArguments(os.Args[1:])
	handleErr(err)

	var wg sync.WaitGroup

	urls := make(chan string)
	errs1 := make(chan error)
	wg.Add(1)
	go func(path string) {
		defer wg.Done()
		defer close(urls)
		defer close(errs1)
		Read(path, urls, errs1)
	}(arg.Config)

	wg.Add(1)
	go printErrs(errs1, &wg)

	items := make(chan RssItem)
	errs2 := make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(items)
		defer close(errs2)
		Fetch(arg.Timed, urls, items, errs2)
	}()

	wg.Add(1)
	go printErrs(errs2, &wg)

	errs3 := make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(errs3)
		client := NewNotionClient(arg.Url, arg.Secret, arg.Version, arg.Database)
		client.Save(items, errs3)
	}()

	wg.Add(1)
	go printErrs(errs3, &wg)

	wg.Wait()
}

func printErrs(errs <-chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for err := range errs {
		log.Printf("error occurred: %v\n", err)
	}

	os.Exit(1)
}

func handleErr(err error) {
	if err != nil {
		log.Printf("error occurred: %v\n", err)
		os.Exit(1)
	}
}
