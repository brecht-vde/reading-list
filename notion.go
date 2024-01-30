package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Title          TitleProp `json:"Title"`
	Author         TextProp  `json:"Author"`
	Guid           TextProp  `json:"Guid"`
	PublishingDate DateProp  `json:"Publishing Date"`
	Url            UrlProp   `json:"Url"`
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

func (n *NotionClient) SaveArticle(article *Article) error {
	url := fmt.Sprintf("%v/pages", n.url)
	page := mapToNotion(article, n.database)

	data, err := json.Marshal(page)

	if err != nil {
		return err
	}

	request, err := n.newNotionRequest("POST", url, bytes.NewReader(data))

	if err != nil {
		return err
	}

	response, err := n.client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	message, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("invalid status code for saving article: %v", string(message))
	}

	return nil
}

func mapToNotion(article *Article, id string) Page {
	return Page{
		Parent: Parent{
			DatabaseId: id,
		},
		Properties: Properties{
			Title: TitleProp{
				Title: []PageRichText{
					{
						Text: Text{
							Content: article.Title,
						},
					},
				},
			},
			Author: TextProp{
				RichText: []PageRichText{
					{
						Text: Text{
							Content: article.Author,
						},
					},
				},
			},
			Guid: TextProp{
				RichText: []PageRichText{
					{
						Text: Text{
							Content: article.Guid,
						},
					},
				},
			},
			PublishingDate: DateProp{
				Date: DatePropType{
					Start: article.PublishingDate.Format("2006-01-02"),
				},
			},
			Url: UrlProp{
				Url: article.Url,
			},
		},
	}
}
