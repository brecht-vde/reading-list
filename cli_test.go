package main

import (
	"reflect"
	"testing"
)

func Test_ParseArguments(t *testing.T) {
	input := []string{
		"-u",
		"url_test",
		"-s",
		"secret_test",
		"-v",
		"version_test",
		"-d",
		"database_test",
		"-h",
		"histories_path_test",
		"-b",
		"blogs_path_test",
		"-y",
		"2023",
	}

	expected := &Arguments{
		Url:       "url_test",
		Secret:    "secret_test",
		Version:   "version_test",
		Database:  "database_test",
		Histories: "histories_path_test",
		Blogs:     "blogs_path_test",
		Year:      2023,
	}

	args, err := ParseArguments(input)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("expected: %v, but received: %v", expected, args)
	}
}

func Test_ParseArguments_Minimal(t *testing.T) {
	input := []string{
		"-s",
		"secret_test",
		"-d",
		"database_test",
	}

	expected := &Arguments{
		Url:       "https://api.notion.com/v1",
		Secret:    "secret_test",
		Version:   "2022-06-28",
		Database:  "database_test",
		Histories: "./resources/histories.json",
		Blogs:     "./resources/blogs.csv",
		Year:      0,
	}

	args, err := ParseArguments(input)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("expected: %v, but received: %v", expected, args)
	}
}
