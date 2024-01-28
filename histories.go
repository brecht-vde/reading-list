package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type History struct {
	Tag string   `json:"tag"`
	Ids []string `json:"ids"`
}

func LoadHistories(path string) ([]History, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return loadHistories(file)
}

func loadHistories(r io.Reader) ([]History, error) {
	data, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var histories []History
	err = json.Unmarshal(data, &histories)

	if err != nil {
		return nil, err
	}

	return histories, nil
}

func SaveHistories(path string, histories []History) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	return saveHistories(file, histories)
}

func saveHistories(w io.Writer, histories []History) error {
	data, err := json.Marshal(histories)

	if err != nil {
		return err
	}

	n, err := w.Write(data)

	if err != nil {
		return err
	}

	if n <= 0 {
		return fmt.Errorf("could not write histories to file")
	}

	return nil
}
