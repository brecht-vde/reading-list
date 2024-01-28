package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_loadBlog(t *testing.T) {
	input := "tag,url\nblog,https://blog.test"

	expected := []Blog{
		{
			Tag: "blog",
			Url: "https://blog.test",
		},
	}

	reader := strings.NewReader(input)

	blogs, err := load(reader)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(blogs, expected) {
		t.Errorf("expected: %v, but received: %v", expected, blogs)
	}
}
