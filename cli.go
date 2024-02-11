package main

import (
	"errors"
	"flag"
	"fmt"
)

type Arguments struct {
	Url      string
	Secret   string
	Version  string
	Database string
	Config   string
	Timed    bool
}

func ParseArguments(args []string) (*Arguments, error) {
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)

	var url = fs.String("u", "https://api.notion.com/v1", "Notion API url e.g. '-u https://api.notion.com/v1', default value 'https://api.notion.com/v1'")
	var secret = fs.String("s", "", "Notion API integration secret e.g. '-s secret_...'")
	var version = fs.String("v", "2022-06-28", "Notion API version e.g. '-v 2022-06-28', default value '2022-06-28'")
	var database = fs.String("d", "", "Notion Database ID e.g. '-d 00000000000000000000000000000000'")
	var config = fs.String("c", "./resources/config.csv", "config file path e.g. '-b ./resources/config.csv', default value './resources/config.csv'")
	var timed = fs.Bool("t", true, "Limits the articles fetched to today's date, default value true")

	err := fs.Parse(args)

	if err != nil {
		return nil, err
	}

	err = validateFlags(fs)

	if err != nil {
		return nil, err
	}

	arguments := Arguments{
		Url:      *url,
		Secret:   *secret,
		Version:  *version,
		Database: *database,
		Config:   *config,
		Timed:    *timed,
	}

	return &arguments, nil
}

func validateFlags(fs *flag.FlagSet) error {
	var err error

	fs.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			err = errors.Join(err, fmt.Errorf("required argument '%v' is not defined", f.Name))
		}
	})

	return err
}
