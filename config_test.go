package main

import (
	"strings"
	"testing"
)

func Test_read(t *testing.T) {
	content := "http://link1.local"
	reader := strings.NewReader(content)
	urls := make(chan string, 1)
	errs := make(chan error, 1)
	close(errs)

	read(reader, urls, errs)

	url := <-urls

	if content != url {
		t.Errorf("expected: %v, but received: %v", content, url)
	}
}
