package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Blog struct {
	Tag string
	Url string
}

func Load(path string) ([]Blog, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return load(file)
}

func load(r io.Reader) ([]Blog, error) {
	data, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	if lines == nil || len(lines) <= 1 {
		return nil, fmt.Errorf("file does not contain content")
	}

	headers := strings.Split(lines[0], ",")
	err = validateHeaders(headers)

	if err != nil {
		return nil, err
	}

	return parseValues(lines[1:])
}

func validateHeaders(headers []string) error {
	if headers == nil || len(headers) != 2 {
		return fmt.Errorf("csv file does not contain required amount of headers")
	}

	if !strings.EqualFold(headers[0], "tag") {
		return fmt.Errorf("csv does not contain header 'tag' at position 0")
	}

	if !strings.EqualFold(headers[1], "url") {
		return fmt.Errorf("csv does not contain header 'url' at position 1")
	}

	return nil
}

func parseValues(lines []string) ([]Blog, error) {
	var err error
	var blogs []Blog

	for i, line := range lines {
		values := strings.Split(line, ",")
		validationErr := validateValues(values)

		if validationErr != nil {
			err = errors.Join(fmt.Errorf("invalid format at index '%v', %v", i, validationErr))
			continue
		}

		blog := Blog{
			Tag: values[0],
			Url: values[1],
		}

		blogs = append(blogs, blog)
	}

	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func validateValues(values []string) error {
	if values == nil || len(values) != 2 {
		return fmt.Errorf("incorrect amount of values")
	}

	if values[0] == "" {
		return fmt.Errorf("tag cannot be empty")
	}

	if values[1] == "" {
		return fmt.Errorf("url cannot be empty")
	}

	return nil
}
