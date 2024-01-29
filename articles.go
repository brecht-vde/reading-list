// articles.go
package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Article struct {
	Tag            string
	Title          string
	Author         string
	Url            string
	PublishingDate time.Time
	Guid           string
}

type rssFeed struct {
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Guid    string `xml:"guid"`
	Author  string `xml:"creator"`
}

func GetArticles(blog *Blog) ([]Article, error) {
	response, err := http.Get(blog.Url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		subErr := parseFailure(response.Body)
		return nil, fmt.Errorf("could not fetch blog %v at %v because %v", blog.Tag, blog.Url, subErr)
	}

	feed, err := parseRssFeed(response.Body)

	if err != nil {
		return nil, err
	}

	return mapToArticle(feed, blog)
}

func parseFailure(reader io.Reader) error {
	data, err := io.ReadAll(reader)

	if err != nil {
		return err
	}

	content := string(data)

	return fmt.Errorf(content)
}

func parseRssFeed(reader io.Reader) (*rssFeed, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	var feed rssFeed
	err = xml.Unmarshal(data, &feed)

	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func mapToArticle(feed *rssFeed, blog *Blog) ([]Article, error) {
	var err error
	var articles []Article

	for _, item := range feed.Channel.Items {
		date, dateErr := time.Parse(time.RFC1123, item.PubDate)

		if dateErr != nil {
			err = errors.Join(fmt.Errorf("could not parse pub date for blog: %v, item: %v", blog.Tag, item.Guid))
			continue
		}

		article := Article{
			Tag:            blog.Tag,
			Title:          item.Title,
			Author:         item.Author,
			Url:            item.Link,
			PublishingDate: date,
			Guid:           item.Guid,
		}

		articles = append(articles, article)
	}

	return articles, err
}
