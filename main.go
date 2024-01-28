package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// Step 0: Parse CLI arguments
	config, err := TryParseConfiguration(os.Args)

	if err != nil {
		log.Fatalf("configuration could not be parsed: %v", err)
		os.Exit(-1)
	}

	// Step 1: load history file
	histories, err := LoadHistoriesV1(config.History)

	if err != nil {
		log.Fatalf("history could not be loaded: %v", err)
		os.Exit(-1)
	}

	// Step 2: Fetch database & read description field
	client := NewNotionClient(*config)
	database, err := client.GetDatabase(config.NotionDatabaseId)

	if err != nil {
		log.Fatalf("database could not be loaded: %v", err)
		os.Exit(-1)
	}

	// Step 3: Parse description field & create a slice of URLs to fetch feeds from
	urls, err := database.ExtractUrls()

	if err != nil {
		log.Fatalf("urls could not be extracted from the database: %v", err)
		os.Exit(-1)
	}

	// Step 4: Fetch each feed
	feeds, err := GetRssFeeds(urls)

	if err != nil {
		log.Fatalf("feeds could not be fetched: %v", err)
		os.Exit(-1)
	}

	// Step 5: map each item into a dao
	// TODO: get updated histories back
	articles, err := GetUniqueArticleV1s(feeds, histories)

	if err != nil {
		log.Fatalf("articles could not be mapped: %v", err)
		os.Exit(-1)
	}

	// Step 6: add entries to the notion database
	err = client.PostArticleV1s(config.NotionDatabaseId, articles)

	if err != nil {
		log.Fatalf("articles could not be posted: %v", err)
		os.Exit(-1)
	}

	// Step 7: modify history.yml and commit back to github
}

// Step 0
// TODO: add an option for setting a date limit (e.g. only import posts from year 2023 & later)
// TODO: add a flag for the blogposts csv (using database description is too complicated)
type Configuration struct {
	NotionApiUrl     string
	NotionApiSecret  string
	NotionApiVersion string
	NotionDatabaseId string
	History          string
}

func TryParseConfiguration(args []string) (*Configuration, error) {
	var config *Configuration
	var err error

	var secret = flag.String("s", "", "your notion api integration secret")
	var url = flag.String("u", "", "the notion api base url")
	var version = flag.String("v", "2022-06-28", "the notion api version")
	var database = flag.String("d", "", "the notion database id")
	var history = flag.String("h", "", "path to the history file")

	flag.Parse()

	var messages string
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			messages += fmt.Sprintf(`flag '-%v', cannot be empty. `, f.Name)
		}
	})

	if messages != "" {
		err = fmt.Errorf(messages)
	}

	if err != nil {
		return nil, err
	}

	config = &Configuration{
		NotionApiUrl:     *url,
		NotionApiSecret:  *secret,
		NotionApiVersion: *version,
		NotionDatabaseId: *database,
		History:          *history,
	}

	return config, err
}

// Step 1 load & read history file
type HistoryV1 struct {
	Url string   `json:"url"`
	Ids []string `json:"ids"`
}

func LoadHistoriesV1(path string) ([]HistoryV1, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return loadHistoriesInternal(file)
}

func loadHistoriesInternal(reader io.Reader) ([]HistoryV1, error) {
	bytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	var histories []HistoryV1
	err = json.Unmarshal(bytes, &histories)

	if err != nil {
		return nil, err
	}

	return histories, nil
}

// Step 2 fetch notion db
// TODO: discard this, we won't use the database description for storing the feeds
// TODO: simplify notion client, allow adding body
// TODO: build in a retry mechanism? Notion API is rate limited at around 3 requests per second
// TODO: maybe useful to check if the database exists
type Database struct {
	Description []RichText `json:"description"`
}

type RichText struct {
	Type string `json:"type"`
	Text Text   `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type NotionClient struct {
	client  *http.Client
	url     string
	secret  string
	version string
}

func NewNotionClient(configuration Configuration) *NotionClient {
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &NotionClient{
		client:  client,
		url:     configuration.NotionApiUrl,
		secret:  configuration.NotionApiSecret,
		version: configuration.NotionApiVersion,
	}
}

func (n *NotionClient) newNotionRequest(method, url string) (*http.Request, error) {
	request, error := http.NewRequest(method, url, nil)

	if error != nil {
		return nil, error
	}

	request.Header.Add("Authorization", n.secret)
	request.Header.Add("Notion-Version", n.version)

	return request, nil
}

func (n *NotionClient) GetDatabase(id string) (*Database, error) {
	url := fmt.Sprintf(`%v/databases/%v`, n.url, id)

	request, err := n.newNotionRequest("GET", url)

	if err != nil {
		return nil, err
	}

	response, err := n.client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("could not get database: %v", response.Status)
	}

	return getDatabaseInternal(response.Body)
}

func getDatabaseInternal(reader io.Reader) (*Database, error) {
	bytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	var database Database

	err = json.Unmarshal(bytes, &database)

	if err != nil {
		return nil, err
	}

	return &database, nil
}

// Step 3: parse database for feeds
// TODO: disregard this, won't be needed
func (d *Database) ExtractUrls() ([]string, error) {
	content := d.Description[0].Text.Content

	if content == "" {
		return nil, fmt.Errorf("the description field does not contain any content")
	}

	urls := strings.Split(content, "\n")

	return urls, nil
}

// Step 4: fetch content from all feeds
type RssFeed struct {
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Items []RssItem `xml:"item"`
}

type RssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Guid    string `xml:"guid"`
	Author  string `xml:"creator"`
}

// TODO: parallel? better error handling, app shouldn't break if 1 feed is down...
// TODO: check if channels / waitgroups offer a better solution for this
// TODO: maybe an iterator pattern?
func GetRssFeeds(urls []string) ([]RssFeed, error) {
	var messages string
	var feeds []RssFeed

	for i := 0; i < len(urls); i++ {
		response, err := http.Get(urls[i])

		if err != nil {
			messages += fmt.Sprintf("could not fetch feed: '%v'. reason: %v", urls[i], err)
			break
		}

		if response.StatusCode != 200 {
			messages += fmt.Sprintf("could not fetch feed: '%v'. reason: %v", urls[i], response.Status)
			break
		}

		defer response.Body.Close()

		feed, err := getRssFeedInternal(response.Body)

		if err != nil {
			messages += fmt.Sprintf("could not parse feed: '%v'. reason: %v", urls[i], err)
			break
		}

		feeds = append(feeds, *feed)
	}

	if messages != "" {
		return nil, fmt.Errorf(messages)
	}

	return feeds, nil
}

func getRssFeedInternal(reader io.Reader) (*RssFeed, error) {
	bytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	var feed RssFeed
	err = xml.Unmarshal(bytes, &feed)

	if err != nil {
		return nil, err
	}

	return &feed, nil
}

// step 5: make a model to represent articles
// TODO: do we need ArticleV1 as "domain" model?
// TODO: add a tag so we know which blog the article belongs to (some blogs have different authors)
type ArticleV1 struct {
	Title          string
	Author         string
	Guid           string
	PublishingDate time.Time
	Url            string
}

// TODO: not liking the uniqueness check, perhaps the initial map should be passed as an argument instead of building it here, not sure yet
// TODO: the RssFeed object does not indicate which blog it belongs to, so need some sort of property to use as identifier
func GetUniqueArticleV1s(feeds []RssFeed, histories []HistoryV1) ([]ArticleV1, error) {
	unique := make(map[string]struct{})
	var articles []ArticleV1

	// TODO: move this elsewhere
	for _, history := range histories {
		for i := 0; i < len(history.Ids); i++ {
			_, exists := unique[history.Ids[i]]

			if !exists {
				unique[history.Ids[i]] = struct{}{}
			}
		}
	}

	for i := 0; i < len(feeds); i++ {
		for j := 0; j < len(feeds[i].Channel.Items); j++ {
			item := feeds[i].Channel.Items[j]

			_, exists := unique[item.Guid]

			if exists {
				continue
			}

			article := ArticleV1{}
			article.Author = item.Author
			article.Guid = item.Guid
			article.Title = item.Title
			article.Url = item.Link

			date, err := time.Parse(time.RFC1123, item.PubDate)

			if err != nil {
				return nil, err
			}

			article.PublishingDate = date

			articles = append(articles, article)
			unique[article.Guid] = struct{}{}
		}
	}

	return articles, nil
}

// Step 6 map to notion & save
// TODO: find some good naming for all of this. The structure is kind of hardcoded to specific properties. Then again, is it worth abstracting more? Probably not since I'm not writing a full API wrapper
// PS: notion api is annoying :)
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

// TODO: better erorr handling, parallel requests?
func (n *NotionClient) PostArticleV1s(id string, articles []ArticleV1) error {
	for i := 0; i < len(articles); i++ {
		err := n.postArticleV1Internal(id, articles[i])

		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NotionClient) postArticleV1Internal(id string, article ArticleV1) error {
	url := fmt.Sprintf(`%v/pages/`, n.url)

	request, err := n.newNotionRequest("POST", url)

	if err != nil {
		return err
	}

	notionArticleV1 := article.toNotion(id)
	data, err := json.Marshal(notionArticleV1)

	if err != nil {
		return err
	}

	request.Body = io.NopCloser(bytes.NewBuffer(data))
	request.Header.Add("Content-Type", "application/json")

	response, err := n.client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	test, _ := io.ReadAll(response.Body)
	_ = test

	if response.StatusCode != 200 {
		return fmt.Errorf("could not post article: %v", response.Status)
	}

	return nil
}

// TODO: not sure if this is a good convention in golang?
func (a *ArticleV1) toNotion(id string) Page {
	return Page{
		Parent: Parent{
			DatabaseId: id,
		},
		Properties: Properties{
			Title: TitleProp{
				Title: []PageRichText{
					PageRichText{
						Text: Text{
							Content: a.Title,
						},
					},
				},
			},
			Author: TextProp{
				RichText: []PageRichText{
					PageRichText{
						Text: Text{
							Content: a.Author,
						},
					},
				},
			},
			Guid: TextProp{
				RichText: []PageRichText{
					PageRichText{
						Text: Text{
							Content: a.Guid,
						},
					},
				},
			},
			PublishingDate: DateProp{
				Date: DatePropType{
					Start: a.PublishingDate.Format("2006-01-02"),
				},
			},
			Url: UrlProp{
				Url: a.Url,
			},
		},
	}
}

// TODO: write back history file
