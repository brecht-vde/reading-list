package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type notionClient struct {
	BaseUrl string
	Version string
	Secret  string
}

func NewNotionClient(secret string) notionClient {
	return notionClient{
		BaseUrl: "https://api.notion.com/v1",
		Version: "2022-06-28",
		Secret:  secret,
	}
}

func (c *notionClient) Post(path string, body io.Reader) error {
	var err error

	url := fmt.Sprintf("%v/%v", c.BaseUrl, path)

	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		return err
	}

	request.Header.Add("Authorization", c.Secret)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Notion-Version", c.Version)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf("notion request did not complete successfully")
		return err
	}

	return err
}

type CreateDatabaseRequestParameters struct {
	ParentId string
	Title    string
	Url      string
}

func NewCreateDatabaseRequest(parameters CreateDatabaseRequestParameters) io.Reader {
	request := fmt.Sprintf(createDatabaseRequestTemplate, parameters.ParentId, parameters.Title, parameters.Url)
	return strings.NewReader(request)
}

const createDatabaseRequestTemplate = `
{
    "parent": {
        "type": "page_id",
        "page_id": "%v"
    },
    "title": [
        {
            "type": "text",
            "text": {
                "content": "%v",
                "link": {
                    "url": "%v"
                }
            }
        }
    ],
    "properties": {
        "Title": {
            "title": {}
        },
        "Author": {
            "rich_text": {}
        },
        "Url": {
            "url": {}
        },
        "Publishing Date": {
            "date": {}
        },
        "Read": {
            "checkbox": {}
        }
    }
}`

type InsertRowRequestParameters struct {
	DatabaseId     string
	Title          string
	Author         string
	Url            string
	PublishingDate string
}

func NewInsertRowRequest(parameters InsertRowRequestParameters) io.Reader {
	request := fmt.Sprintf(insertRowRequestParametersTemplate, parameters.DatabaseId, parameters.Title, parameters.Author, parameters.Url, parameters.PublishingDate)
	return strings.NewReader(request)
}

const insertRowRequestParametersTemplate = `
{
    "parent": {
        "type": "database_id",
        "database_id": "%v"
    },
    "properties": {
        "Title": {
            "id": "title",
            "type": "title",
            "title": [
                {
                    "type": "text",
                    "text": {
                        "content": "%v"
                    }
                }
            ]
        },
        "Author": {
            "type": "rich_text",
            "rich_text": [
                {
                    "type": "text",
                    "text": {
                        "content": "%v"
                    }
                }
            ]
        },
        "Url": {
            "type": "url",
            "url": "%v"
        },
        "Publishing Date": {
            "type": "date",
            "date": {
                "start": "%v"
            }
        }
    }
}`
