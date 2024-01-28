package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestHistoriesMarshalling(t *testing.T) {
	expected := "[{\"url\":\"https://test1.local\",\"ids\":[\"123\",\"456\"]},{\"url\":\"https://test2.local\",\"ids\":[\"789\",\"012\"]}]"

	histories := []History{
		{
			Url: "https://test1.local",
			Ids: []string{
				"123",
				"456",
			},
		},
		{
			Url: "https://test2.local",
			Ids: []string{
				"789",
				"012",
			},
		},
	}

	bytes, err := json.Marshal(histories)

	if err != nil {
		t.Fail()
	}

	if string(bytes) != expected {
		t.Fail()
	}
}

func TestHistoriesUnmarshalling(t *testing.T) {
	expected := []History{
		{
			Url: "https://test1.local",
			Ids: []string{
				"123",
				"456",
			},
		},
		{
			Url: "https://test2.local",
			Ids: []string{
				"789",
				"012",
			},
		},
	}

	var histories []History
	data := "[{\"url\":\"https://test1.local\",\"ids\":[\"123\",\"456\"]},{\"url\":\"https://test2.local\",\"ids\":[\"789\",\"012\"]}]"

	err := json.Unmarshal([]byte(data), &histories)

	if err != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(histories, expected) {
		t.Fail()
	}
}
