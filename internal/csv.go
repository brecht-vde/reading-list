package internal

import (
	"fmt"
	"io"
	"strings"
)

type Csv struct {
	Rows []CsvRow
}

type CsvRow struct {
	Name string
	Url  string
}

func LoadCsv(data *io.Reader) (Csv, error) {
	var csv Csv
	var err error

	blob, err := loadCsv(data)

	if err != nil {
		return csv, err
	}

	csv, err = parseCsv(blob)

	if err != nil {
		return csv, err
	}

	return csv, err
}

func loadCsv(data *io.Reader) ([]byte, error) {
	var blob []byte
	var err error

	blob, err = io.ReadAll(*data)

	if err != nil {
		return blob, err
	}

	return blob, err
}

func parseCsv(blob []byte) (Csv, error) {
	var csv Csv
	var err error

	text := string(blob)

	if text == "" {
		err = fmt.Errorf("cannot parse csv files without content")
		return csv, err
	}

	text = strings.ReplaceAll(text, "\r", "")

	rows, err := parseCsvRows(text)

	if err != nil {
		return csv, err
	}

	csv, err = parseCsvValues(rows)

	if err != nil {
		return csv, err
	}

	return csv, err
}

func parseCsvRows(text string) ([]string, error) {
	var rows []string
	var err error

	rows = strings.Split(text, "\n")

	if len(rows) < 2 {
		err = fmt.Errorf("cannot parse csv files with less than 2 rows")
		return rows, err
	}

	return rows, err
}

func parseCsvValues(rows []string) (Csv, error) {
	var csv Csv = Csv{}
	var err error

	for i := 1; i < len(rows); i++ {
		values := strings.Split(rows[i], ",")

		if len(values) != 2 {
			err = fmt.Errorf("invalid row found at index %v: %v", i, err)
			return csv, err
		}

		csvRow := CsvRow{
			Name: values[0],
			Url:  values[1],
		}

		csv.Rows = append(csv.Rows, csvRow)
	}

	return csv, err
}
