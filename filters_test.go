package main

import (
	"reflect"
	"testing"
)

func Test_FilterArticles(t *testing.T) {
	articles := []Article{
		{
			Guid: "123",
		},
		{
			Guid: "456",
		},
		{
			Guid: "789",
		},
	}

	history := History{
		Ids: []string{
			"123",
			"456",
		},
	}

	expected := []Article{
		{
			Guid: "789",
		},
	}

	FilterArticles(&articles, &history)

	if !reflect.DeepEqual(expected, articles) {
		t.Errorf("expected: %v, but received: %v", expected, articles)
	}
}
