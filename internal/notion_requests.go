package internal

import (
	"fmt"
	"io"
	"strings"
)

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
	request := fmt.Sprintf(InsertRowRequestParametersTemplate, parameters.DatabaseId, parameters.Title, parameters.Author, parameters.Url, parameters.PublishingDate)
	return strings.NewReader(request)
}

const InsertRowRequestParametersTemplate = `
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