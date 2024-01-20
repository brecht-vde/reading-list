package internal

import (
	"fmt"
	"io"
	"net/http"
)

type NotionClient struct {
	BaseUrl string
	Version string
	Secret  string
}

func (c *NotionClient) Post(path string, body io.Reader) error {
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
		err = fmt.Errorf("could not patch database")
		return err
	}

	return err
}
