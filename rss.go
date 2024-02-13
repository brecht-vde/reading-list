package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type RssFeed struct {
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Title string    `xml:"title"`
	Items []RssItem `xml:"item"`
}

type RssItem struct {
	Blog       string
	Title      string   `xml:"title"`
	Link       string   `xml:"link"`
	PubDate    string   `xml:"pubDate"`
	Categories []string `xml:"category"`
}

func Fetch(timed bool, urls <-chan string, items chan<- RssItem, errs chan<- error) {
	var wg sync.WaitGroup

	for url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			processUrl(timed, url, items, errs)
			log.Printf("processed url: %v\n", url)
		}(url)
	}

	wg.Wait()
}

func processUrl(timed bool, url string, items chan<- RssItem, errors chan<- error) {
	feed, err := fetchUrl(url)

	if err != nil {
		errors <- err
		return
	}

	processFeed(timed, feed, items)
}

func fetchUrl(url string) (*RssFeed, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("did not get a valid response, code: '%v', message: '%v'", response.StatusCode, string(data))
	}

	feed, err := parseResult(data)

	if err != nil {
		return nil, err
	}

	return feed, nil
}

func parseResult(data []byte) (*RssFeed, error) {
	var feed RssFeed
	err := xml.Unmarshal(data, &feed)

	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func processFeed(timed bool, feed *RssFeed, items chan<- RssItem) {
	for _, item := range feed.Channel.Items {
		item.Blog = feed.Channel.Title

		if timed && !shouldInclude(item.PubDate) {
			log.Printf("not including item: %v\n", item.Title)
			continue
		}

		items <- item
		log.Printf("sent item: %v\n", item.Title)
	}
}

func shouldInclude(pubDate string) bool {
	date, err := time.Parse(time.RFC1123, pubDate)

	if err != nil {
		return false
	}

	since := time.Since(date).Hours()

	if since > 24 || since < 0 {
		return false
	}

	return true
}
