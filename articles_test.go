package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_mapToArticle(t *testing.T) {
	blog := &Blog{
		Tag: "test",
		Url: "https://test.local",
	}

	date, err := time.Parse("2006-01-02", "2023-12-31")

	if err != nil {
		t.Errorf("invalid date specified for article in test case")
	}

	expected := &[]Article{
		{
			Tag:            "test",
			Title:          "title",
			Author:         "author",
			Url:            "url",
			Guid:           "guid",
			PublishingDate: date,
		},
	}

	feed := &rssFeed{
		Channel: rssChannel{
			Items: []rssItem{
				{
					Title:   "title",
					Link:    "url",
					Author:  "author",
					PubDate: date.Format(time.RFC1123),
					Guid:    "guid",
				},
			},
		},
	}

	article, err := mapToArticle(feed, blog)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(article, expected) {
		t.Errorf("expected: %v, but received %v", expected, article)
	}
}

func Test_parseRssFeed(t *testing.T) {
	input := `
		<rss>
			<channel>
				<item>
					<title>title</title>
					<link>url</link>
					<guid>guid</guid>
					<dc:creator>author</dc:creator>
					<pubDate>Sun, 31 Dec 2023 00:00:00 GMT</pubDate>
				</item>
			</channel>
		</rss>
	`
	expected := &rssFeed{
		Channel: rssChannel{
			Items: []rssItem{
				{
					Title:   "title",
					Link:    "url",
					PubDate: "Sun, 31 Dec 2023 00:00:00 GMT",
					Guid:    "guid",
					Author:  "author",
				},
			},
		},
	}

	reader := strings.NewReader(input)
	feed, err := parseRssFeed(reader)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(feed, expected) {
		t.Errorf("expected: %v, but received: %v", expected, feed)
	}
}

func Test_parseFailure(t *testing.T) {
	input := "asjdnas;lkdas;dnaskdaslkdjasdkasjdkasd"
	reader := strings.NewReader(input)
	expected := fmt.Errorf(input)

	err := parseFailure(reader)

	if expected.Error() != err.Error() {
		t.Errorf("expected: %v, but received: %v", expected, err)
	}
}
