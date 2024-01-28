package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func Test_LoadHistories(t *testing.T) {
	input := `
	[
		{
			"tag": "test",
			"ids": [
				"123",
				"456"
			]
		}
	]`

	expected := []History{
		{
			Tag: "test",
			Ids: []string{
				"123",
				"456",
			},
		},
	}

	reader := strings.NewReader(input)
	histories, err := loadHistories(reader)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(histories, expected) {
		t.Errorf("expected: %v, but received: %v", expected, histories)
	}
}

func Test_WriteHistories(t *testing.T) {
	histories := []History{
		{
			Tag: "test",
			Ids: []string{
				"123",
				"456",
			},
		},
	}

	var w bytes.Buffer
	err := saveHistories(&w, histories)

	if err != nil {
		t.Error(err)
	}
}
