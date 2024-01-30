package main

import (
	"errors"
	"flag"
	"fmt"
)

type Arguments struct {
	Url       string
	Secret    string
	Version   string
	Database  string
	Histories string
	Blogs     string
	Year      int
}

func ParseArguments(args []string) (*Arguments, error) {
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)

	var url = fs.String("u", "https://api.notion.com/v1", "Notion API url e.g. '-u https://api.notion.com/v1', default value 'https://api.notion.com/v1'")
	var secret = fs.String("s", "", "Notion API integration secret e.g. '-s secret_...'")
	var version = fs.String("v", "2022-06-28", "Notion API version e.g. '-v 2022-06-28', default value '2022-06-28'")
	var database = fs.String("d", "", "Notion Database ID e.g. '-d 00000000000000000000000000000000'")
	var histories = fs.String("h", "./resources/histories.json", "histories file path e.g. '-h ./resources/histories.json', default value './resources/histories.json")
	var blogs = fs.String("b", "./resources/blogs.csv", "blogs file path e.g. '-b ./resources/blogs.csv', default value './resources/blogs.csv'")
	var year = fs.Int("y", 0, "inclusive cut off year for fetching articles e.g. '-y 2023', articles published in 2023 or later will be synchronized")

	err := fs.Parse(args)

	if err != nil {
		return nil, err
	}

	err = validateFlags(fs)

	if err != nil {
		return nil, err
	}

	arguments := Arguments{
		Url:       *url,
		Secret:    *secret,
		Version:   *version,
		Database:  *database,
		Histories: *histories,
		Blogs:     *blogs,
		Year:      *year,
	}

	return &arguments, nil
}

func validateFlags(fs *flag.FlagSet) error {
	var err error

	fs.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			err = errors.Join(fmt.Errorf("required argument '%v' is not defined", f.Name))
		}
	})

	return err
}
