package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
CLI parsing:
  - configuration file path
  - notion api key
  - ...

RSS Reading:

- Read configuration file
- Deserialize configuration
  - split on | character, put in name/url map or some struct

- Http GET urls concurrently
- Parse XML <item> elements
  - Parse XML <title> elements --> may be wrapped in CDATA
  - Parse XML <link> elements

Notion initialization:

- HTTP Post create database document IF not exists

Notion content:

- HTTP Post database entry per parsed XML element:
  - Blog name
  - Post title
  - Post url
*/
func main() {
	command, err := Validate(os.Args[1:])

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	err = Execute(command)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

type Command struct {
	config_path string
}

// pure
func Validate(args []string) (*Command, error) {
	var config_path_flag = flag.String("c", "", "-c <path>/<to>/.config")

	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, err
	}

	if *config_path_flag == "" {
		return nil, fmt.Errorf("argument -c is missing")
	}

	command := &Command{
		config_path: *config_path_flag,
	}

	return command, nil
}

// procedural
func Execute(command *Command) error {
	input, err := read(command)

	if err != nil {
		return err
	}

	items, err := parse(input)

	if err != nil {
		return err
	}

	for key, value := range items {
		log.Printf("name: '%v' url: '%v'", key, value)
	}

	return nil
}

// i/o
func read(command *Command) ([]string, error) {
	file, err := os.Open(command.config_path)

	if err != nil {
		return nil, err
	}

	var items []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		items = append(items, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return items, nil
}

// pure
func parse(items []string) (map[string]string, error) {
	output := make(map[string]string)

	for _, item := range items {
		values := strings.Split(item, "|")

		if len(values) == 0 {
			return nil, fmt.Errorf("invalid value: '%v'", item)
		}

		name := values[0]
		url := values[1]

		if _, ok := output[name]; ok {
			return nil, fmt.Errorf("duplicate entry found for '%v'", name)
		}

		output[name] = url
	}

	return output, nil
}
