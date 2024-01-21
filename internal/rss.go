package internal

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

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
}

func FetchRss(url string) (rssFeed, error) {
	var feed rssFeed
	var err error

	data, err := fetchRss(url)

	if err != nil {
		return feed, err
	}

	feed, err = parseRss(data)

	if err != nil {
		return feed, err
	}

	return feed, err
}

func fetchRss(url string) ([]byte, error) {
	var data []byte
	var err error

	response, err := http.Get(url)

	if err != nil {
		return data, err
	}

	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)

	if err != nil {
		return data, err
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("fetching url: '%v' has an invalid status: '%v', with message: '%v'", url, response.StatusCode, data)
		return data, err
	}

	return data, err
}

func parseRss(blob []byte) (rssFeed, error) {
	var feed rssFeed

	err := xml.Unmarshal(blob, &feed)

	if err != nil {
		return feed, err
	}

	return feed, nil
}

// func fetchMany(urls []string, wg *sync.WaitGroup) ([]feed, error) {

// 	for i := range urls {
// 		wg.Add(i)
// 	}

// 	defer wg.Done()
// }

// func FetchFeeds(urls []string) {
// 	ch := make(chan RssRoot)

// 	for _, url := range urls {
// 		go FetchUrl(url, ch)
// 	}

// 	feeds := []RssRoot{}

// 	for i := 0; i < len(urls); i++ {
// 		feeds = append(feeds, <-ch)
// 	}

// 	for _, feed := range feeds {
// 		for _, item := range feed.Channel.Items {
// 			fmt.Println(item.Title)
// 		}
// 	}
// }

// func FetchUrl(url string, c chan<- RssRoot) {
// 	resp, err := http.Get(url)

// 	if err != nil {
// 		c <- RssRoot{}
// 		return
// 	}

// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)

// 	if err != nil {
// 		c <- RssRoot{}
// 		return
// 	}

// 	feed, err := ParseRss(body)

// 	if err != nil {
// 		c <- RssRoot{}
// 		return
// 	}

// 	c <- feed
// }
