package main

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func Test_parseResult(t *testing.T) {
	expected := &RssFeed{
		Channel: RssChannel{
			Title: "title",
			Items: []RssItem{
				{
					Blog:    "blog",
					Title:   "title",
					Link:    "http://link.local",
					PubDate: "",
					Categories: []string{
						"cat1",
						"cat2",
						"cat3",
					},
				},
			},
		},
	}

	data, err := xml.Marshal(expected)

	if err != nil {
		t.Error(err)
	}

	result, err := parseResult(data)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected: %v, but received: %v", expected, result)
	}
}

func Test_processFeed(t *testing.T) {
	feed := &RssFeed{
		Channel: RssChannel{
			Title: "channel.title",
			Items: []RssItem{
				{
					Blog:    "channel.title",
					Title:   "title",
					Link:    "http://link.local",
					PubDate: "",
					Categories: []string{
						"cat1",
						"cat2",
						"cat3",
					},
				},
			},
		},
	}

	expected := RssItem{
		Blog:    "channel.title",
		Title:   "title",
		Link:    "http://link.local",
		PubDate: "",
		Categories: []string{
			"cat1",
			"cat2",
			"cat3",
		},
	}

	var items = make(chan RssItem, 1)

	processFeed(true, feed, items)

	item := <-items

	if !reflect.DeepEqual(item, expected) {
		t.Errorf("expected: %v, but received: %v", expected, item)
	}
}
