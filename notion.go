package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type NotionClient struct {
	client   *http.Client
	url      string
	secret   string
	version  string
	database string
}

type Page struct {
	Parent     Parent     `json:"parent"`
	Properties Properties `json:"properties"`
}

type Parent struct {
	DatabaseId string `json:"database_id"`
}

type Properties struct {
	Blog           SelectProp      `json:"Blog"`
	Title          TitleProp       `json:"Title"`
	PublishingDate DateProp        `json:"Publishing Date"`
	Categories     MultiSelectProp `json:"Categories"`
	Url            UrlProp         `json:"Url"`
}

type Text struct {
	Content string `json:"content"`
}

type TitleProp struct {
	Title []PageRichText `json:"title"`
}

type TextProp struct {
	RichText []PageRichText `json:"rich_text"`
}

type PageRichText struct {
	Text Text `json:"text"`
}

type DateProp struct {
	Date DatePropType `json:"date"`
}

type DatePropType struct {
	Start string `json:"start"`
}

type UrlProp struct {
	Url string `json:"url"`
}

type SelectProp struct {
	Select SelectPropType `json:"select"`
}

type SelectPropType struct {
	Name string `json:"name"`
}

type MultiSelectProp struct {
	MultiSelect []SelectPropType `json:"multi_select"`
}

func NewNotionClient(url, secret, version, database string) *NotionClient {
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &NotionClient{
		client:   client,
		url:      url,
		secret:   secret,
		version:  version,
		database: database,
	}
}

func (n *NotionClient) Save(items <-chan RssItem, errs chan<- error) {
	var wg sync.WaitGroup

	for item := range items {
		wg.Add(1)
		go func(item RssItem) {
			defer wg.Done()
			n.processItem(item, errs)
			log.Printf("saved item: %v\n", item.Title)
		}(item)
	}

	wg.Wait()
}

func (n *NotionClient) newNotionRequest(method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", n.secret)
	request.Header.Add("Notion-Version", n.version)
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func (n *NotionClient) processItem(item RssItem, errors chan<- error) {
	url := fmt.Sprintf("%v/pages", n.url)
	page := mapItem(item, n.database)

	data, err := json.Marshal(page)

	if err != nil {
		errors <- err
		return
	}

	request, err := n.newNotionRequest("POST", url, bytes.NewReader(data))

	if err != nil {
		errors <- err
		return
	}

	response, err := n.withRetry(request)

	if err != nil {
		errors <- fmt.Errorf("%v: ", err)
		return
	}

	defer response.Body.Close()

	message, err := io.ReadAll(response.Body)

	if err != nil {
		errors <- err
		return
	}

	if response.StatusCode != 200 {
		errors <- fmt.Errorf("%v - %v: %v\n%v", response.StatusCode, item.Blog, string(message), string(data))
		return
	}
}

func (n *NotionClient) withRetry(request *http.Request) (*http.Response, error) {
	for i := 1; i <= 6; i++ {
		response, err := n.client.Do(request)

		if err != nil {
			return nil, err
		}

		if i == 6 {
			return response, nil
		}

		if response.StatusCode >= 400 {
			backoff := math.Pow(2, float64(i))
			jitter := rand.Intn(10)
			time.Sleep(time.Second * time.Duration(int(backoff)+jitter))
			continue
		}

		return response, nil
	}

	return nil, fmt.Errorf("failed completely")
}

func mapItem(item RssItem, db string) Page {
	return Page{
		Parent: Parent{
			DatabaseId: db,
		},
		Properties: Properties{
			Title: TitleProp{
				Title: []PageRichText{
					{
						Text: Text{
							Content: sanitize(item.Title),
						},
					},
				},
			},
			Blog: SelectProp{
				Select: SelectPropType{
					Name: sanitize(item.Blog),
				},
			},
			PublishingDate: DateProp{
				Date: DatePropType{
					Start: mapPubDate(item.PubDate),
				},
			},
			Categories: MultiSelectProp{
				MultiSelect: mapCategories(item.Categories),
			},
			Url: UrlProp{
				Url: item.Link,
			},
		},
	}
}

func mapCategories(categories []string) []SelectPropType {
	multiselect := make([]SelectPropType, 0)

	for _, c := range categories {
		prop := SelectPropType{
			Name: sanitize(c),
		}

		multiselect = append(multiselect, prop)
	}

	return multiselect
}

func sanitize(item string) string {
	return strings.ReplaceAll(item, ",", " ")
}

func mapPubDate(pubDate string) string {
	date, err := time.Parse(time.RFC1123, pubDate)

	if err != nil {
		return ""
	}

	return date.Format("2006-01-02T15:04:05Z07:00")
}
