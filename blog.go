// blog.go
package main

import (
	"bufio"
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

func LoadBlogs(path string) ([]*Blog, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return load(file)
}

func load(r io.Reader) ([]*Blog, error) {
	var err error
	scanner := bufio.NewScanner(r)

	if !scanner.Scan() {
		return nil, fmt.Errorf("file does not contain content")
	}

	headers := strings.Split(scanner.Text(), ",")
	err = validateHeaders(headers)

	if err != nil {
		return nil, err
	}

	var blogs []*Blog

	for scanner.Scan() {
		blog, parseErr := parseValue(scanner.Text())

		if parseErr != nil {
			err = errors.Join(err, parseErr)
			continue
		}

		blogs = append(blogs, blog)
	}

	if err != nil {
		return nil, err
	}

	return blogs, nil
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

func parseValue(line string) (*Blog, error) {
	values := strings.Split(line, ",")
	err := validateValues(values)

	if err != nil {
		return nil, err
	}

	blog := Blog{
		Tag: values[0],
		Url: values[1],
	}

	return &blog, nil
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
