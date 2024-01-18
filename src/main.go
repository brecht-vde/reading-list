package main

import (
	"bufio"
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

/*
CLI parsing:
  - configuration file path
  - notion api key
  - ...

RSS Reading:

- Read configuration file
- Deserialize configuration
  - split on | character, put in name/url map or some struct

- Http GET urls concurrently
- Parse XML <item> elements
  - Parse XML <title> elements --> may be wrapped in CDATA
  - Parse XML <link> elements

Notion initialization:

- HTTP Post create database document IF not exists

Notion content:

- HTTP Post database entry per parsed XML element:
  - Blog name
  - Post title
  - Post url
*/
func main() {

	updateDatabase()

	saveArticles(FakeArticles)

	FetchFeeds(urls)

	blob := []byte(test_rss)

	ParseRss(blob)

	command, err := Validate(os.Args[1:])

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	err = Execute(command)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

type Command struct {
	config_path string
}

// pure
func Validate(args []string) (*Command, error) {
	var config_path_flag = flag.String("c", "", "-c <path>/<to>/.config")

	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, err
	}

	if *config_path_flag == "" {
		return nil, fmt.Errorf("argument -c is missing")
	}

	command := &Command{
		config_path: *config_path_flag,
	}

	return command, nil
}

// procedural
func Execute(command *Command) error {
	input, err := read(command)

	if err != nil {
		return err
	}

	items, err := parse(input)

	if err != nil {
		return err
	}

	for key, value := range items {
		log.Printf("name: '%v' url: '%v'", key, value)
	}

	return nil
}

// i/o
func read(command *Command) ([]string, error) {
	file, err := os.Open(command.config_path)

	if err != nil {
		return nil, err
	}

	var items []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		items = append(items, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return items, nil
}

// pure
func parse(items []string) (map[string]string, error) {
	output := make(map[string]string)

	for _, item := range items {
		values := strings.Split(item, "|")

		if len(values) == 0 {
			return nil, fmt.Errorf("invalid value: '%v'", item)
		}

		name := values[0]
		url := values[1]

		if _, ok := output[name]; ok {
			return nil, fmt.Errorf("duplicate entry found for '%v'", name)
		}

		output[name] = url
	}

	return output, nil
}

/*
	RSS, necessary items:

	<rss>
		<channel>
			<item>
				<title></title>
				<link></link>
				<pubDate></pubDate>
			</item>
			...
		</channel>
	</rss>
*/

var test_rss string = `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:dc="http://purl.org/dc/elements/1.1/" version="2.0">
  <channel>
    <title>HubSpot Product Blog (Live)</title>
    <link>https://product.hubspot.com/blog</link>
    <description>The HubSpot Product Blog is where we riff about topics of interest to us from engineering, to UX, to product management and beyond.</description>
    <language>en-us</language>
    <pubDate>Tue, 09 Jan 2024 15:29:37 GMT</pubDate>
    <dc:date>2024-01-09T15:29:37Z</dc:date>
    <dc:language>en-us</dc:language>
    <item>
      <title>Looking Back: Our Most Popular Posts of 2023</title>
      <link>https://product.hubspot.com/blog/2023-retro</link>
      <description>&lt;div class="hs-featured-image-wrapper"&gt; 
 &lt;a href="https://product.hubspot.com/blog/2023-retro" title="" class="hs-featured-image-link"&gt; &lt;img src="https://product.hubspot.com/hubfs/Preventing%20Serial%20Processing%20%281%29.png" alt="Looking Back: Our Most Popular Posts of 2023" class="hs-featured-image" style="width:auto !important; max-width:50%; float:left; margin:0 15px 15px 0;"&gt; &lt;/a&gt; 
&lt;/div&gt; 
&lt;p&gt;2023 was an incredible year of growth for us at HubSpot directly driven by our Product and Engineering teams who tackled hard problems that helped our customers grow better. Some of our HubSpotters shared impactful case studies, impressive growth stories and insights to highlight their work and just what it means to work at HubSpot on our Product Blog. We've gathered some of our most popular blogs of 2023 below in case you missed them, and stayed tuned this year for more behind-the-scenes content from our Product and Engineering teams!&lt;/p&gt; 
&lt;p&gt;Looking for more stories about our culture and what it's like to work at HubSpot? Check out our &lt;a href="https://www.hubspot.com/careers-blog"&gt;Careers Blog&lt;/a&gt;!&lt;/p&gt; 
&lt;p&gt;_________&lt;/p&gt;</description>
      <content:encoded>&lt;div class="hs-featured-image-wrapper"&gt; 
 &lt;a href="https://product.hubspot.com/blog/2023-retro" title="" class="hs-featured-image-link"&gt; &lt;img src="https://product.hubspot.com/hubfs/Preventing%20Serial%20Processing%20%281%29.png" alt="Looking Back: Our Most Popular Posts of 2023" class="hs-featured-image" style="width:auto !important; max-width:50%; float:left; margin:0 15px 15px 0;"&gt; &lt;/a&gt; 
&lt;/div&gt; 
&lt;p&gt;2023 was an incredible year of growth for us at HubSpot directly driven by our Product and Engineering teams who tackled hard problems that helped our customers grow better. Some of our HubSpotters shared impactful case studies, impressive growth stories and insights to highlight their work and just what it means to work at HubSpot on our Product Blog. We've gathered some of our most popular blogs of 2023 below in case you missed them, and stayed tuned this year for more behind-the-scenes content from our Product and Engineering teams!&lt;/p&gt; 
&lt;p&gt;Looking for more stories about our culture and what it's like to work at HubSpot? Check out our &lt;a href="https://www.hubspot.com/careers-blog"&gt;Careers Blog&lt;/a&gt;!&lt;/p&gt; 
&lt;p&gt;_________&lt;/p&gt;  
&lt;img src="https://track.hubspot.com/__ptq.gif?a=51294&amp;amp;k=14&amp;amp;r=https%3A%2F%2Fproduct.hubspot.com%2Fblog%2F2023-retro&amp;amp;bu=https%253A%252F%252Fproduct.hubspot.com%252Fblog&amp;amp;bvt=rss" alt="" width="1" height="1" style="min-height:1px!important;width:1px!important;border-width:0!important;margin-top:0!important;margin-bottom:0!important;margin-right:0!important;margin-left:0!important;padding-top:0!important;padding-bottom:0!important;padding-right:0!important;padding-left:0!important; "&gt;</content:encoded>
      <category>Engineering--Infrastructure</category>
      <category>Culture</category>
      <category>UX</category>
      <category>Engineering</category>
      <category>Product</category>
      <category>Engineering--Backend</category>
      <category>Engineering--Frontend</category>
      <category>Product Management</category>
      <pubDate>Tue, 09 Jan 2024 15:29:37 GMT</pubDate>
      <author>agerow@hubspot.com (Ashlee Gerow)</author>
      <guid>https://product.hubspot.com/blog/2023-retro</guid>
      <dc:date>2024-01-09T15:29:37Z</dc:date>
    </item>
    <item>
      <title>Preventing Serial Processing on the Import Pipeline</title>
      <link>https://product.hubspot.com/blog/improving-performance-import-pipeline</link>
      <description>&lt;div class="hs-featured-image-wrapper"&gt; 
 &lt;a href="https://product.hubspot.com/blog/improving-performance-import-pipeline" title="" class="hs-featured-image-link"&gt; &lt;img src="https://product.hubspot.com/hubfs/Preventing%20Serial%20Processing.png" alt="Preventing Serial Processing on the Import Pipeline" class="hs-featured-image" style="width:auto !important; max-width:50%; float:left; margin:0 15px 15px 0;"&gt; &lt;/a&gt; 
&lt;/div&gt; 
&lt;p&gt;&lt;em&gt;Written by Yash Tulsiani, Technical Lead @ HubSpot.&lt;/em&gt;&lt;br&gt;&lt;br&gt;&lt;em&gt;The HubSpot import system is responsible for ingesting hundreds of millions of spreadsheet rows per day. We translate and write this data into the HubSpot CRM. In this post, we will look at how we solved an edge case in our Kafka consumer that was leading to poor import performance for all HubSpot import users.&lt;br&gt;&lt;/em&gt;&lt;/p&gt; 
&lt;p&gt;_________________&lt;/p&gt;</description>
      <content:encoded>&lt;div class="hs-featured-image-wrapper"&gt; 
 &lt;a href="https://product.hubspot.com/blog/improving-performance-import-pipeline" title="" class="hs-featured-image-link"&gt; &lt;img src="https://product.hubspot.com/hubfs/Preventing%20Serial%20Processing.png" alt="Preventing Serial Processing on the Import Pipeline" class="hs-featured-image" style="width:auto !important; max-width:50%; float:left; margin:0 15px 15px 0;"&gt; &lt;/a&gt; 
&lt;/div&gt; 
&lt;p&gt;&lt;em&gt;Written by Yash Tulsiani, Technical Lead @ HubSpot.&lt;/em&gt;&lt;br&gt;&lt;br&gt;&lt;em&gt;The HubSpot import system is responsible for ingesting hundreds of millions of spreadsheet rows per day. We translate and write this data into the HubSpot CRM. In this post, we will look at how we solved an edge case in our Kafka consumer that was leading to poor import performance for all HubSpot import users.&lt;br&gt;&lt;/em&gt;&lt;/p&gt; 
&lt;p&gt;_________________&lt;/p&gt;  
&lt;img src="https://track.hubspot.com/__ptq.gif?a=51294&amp;amp;k=14&amp;amp;r=https%3A%2F%2Fproduct.hubspot.com%2Fblog%2Fimproving-performance-import-pipeline&amp;amp;bu=https%253A%252F%252Fproduct.hubspot.com%252Fblog&amp;amp;bvt=rss" alt="" width="1" height="1" style="min-height:1px!important;width:1px!important;border-width:0!important;margin-top:0!important;margin-bottom:0!important;margin-right:0!important;margin-left:0!important;padding-top:0!important;padding-bottom:0!important;padding-right:0!important;padding-left:0!important; "&gt;</content:encoded>
      <category>Engineering--Infrastructure</category>
      <category>Engineering</category>
      <category>Engineering--Backend</category>
      <pubDate>Thu, 07 Dec 2023 13:00:00 GMT</pubDate>
      <guid>https://product.hubspot.com/blog/improving-performance-import-pipeline</guid>
      <dc:date>2023-12-07T13:00:00Z</dc:date>
      <dc:creator>Yash Tulsiani</dc:creator>
    </item>
  </channel>
</rss>
`

type RssRoot struct {
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Items []RssItem `xml:"item"`
}

type RssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func ParseRss(blob []byte) (RssRoot, error) {
	var rss RssRoot

	if err := xml.Unmarshal(blob, &rss); err != nil {
		return rss, err
	}

	return rss, nil
}

/*
	https://medium.com/feed/airbnb-engineering
	https://aws.amazon.com/blogs/aws/feed/
	https://www.elastic.co/blog/feed
	https://codeascraft.com/feed/
	https://githubengineering.com/atom.xml
	https://tech.grammarly.com/feed.xml
*/

var urls []string = []string{
	"https://medium.com/feed/airbnb-engineering",
	"https://aws.amazon.com/blogs/aws/feed/",
	"https://codeascraft.com/feed/",
	"https://githubengineering.com/atom.xml",
	"https://tech.grammarly.com/feed.xml",
}

func FetchFeeds(urls []string) {
	ch := make(chan RssRoot)

	for _, url := range urls {
		go FetchUrl(url, ch)
	}

	feeds := []RssRoot{}

	for i := 0; i < len(urls); i++ {
		feeds = append(feeds, <-ch)
	}

	for _, feed := range feeds {
		for _, item := range feed.Channel.Items {
			fmt.Println(item.Title)
		}
	}
}

func FetchUrl(url string, c chan<- RssRoot) {
	resp, err := http.Get(url)

	if err != nil {
		c <- RssRoot{}
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		c <- RssRoot{}
		return
	}

	feed, err := ParseRss(body)

	if err != nil {
		c <- RssRoot{}
		return
	}

	c <- feed
}

/*
	- Create a notion database
	- Add items to the database (author, title, link)
	- Use a specific field to store the last sync date in

*/

type Article struct {
	Title     string
	Author    string
	Link      string
	Published string
}

var FakeArticles []Article = []Article{
	{Title: "Fake 1", Author: "Fake 4", Link: "https://fake7.local", Published: time.Now().Format("YYYY/MM/DD")},
	{Title: "Fake 2", Author: "Fake 5", Link: "https://fake8.local", Published: time.Now().Format("YYYY/MM/DD")},
	{Title: "Fake 3", Author: "Fake 6", Link: "https://fake9.local", Published: time.Now().Format("YYYY/MM/DD")},
}

func saveArticles(articles []Article) {
	ch := make(chan bool)

	for _, article := range articles {
		go PostUrl(article, ch)
	}

	for v := range ch {
		fmt.Println(v)
	}
}

func PostUrl(article Article, c chan<- bool) {
	template := fmt.Sprintf(`
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
			"Link": {
				"type": "url",
				"url": "%v"
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
			"Published": {
				"type": "date",
				"date": {
					"start": "%v"
				}
			}
		}
	}
	`, "ced11cae-8f2c-4448-ab19-5c8ad32a5426", article.Title, article.Link, article.Author, article.Published)

	body := strings.NewReader(template)
	request, err := http.NewRequest("POST", "https://api.notion.com/v1/pages/", body)

	if err != nil {
		c <- false
		return
	}

	request.Header.Add("Authorization", "hidden")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Notion-Version", "2022-06-28")

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	if err != nil {
		c <- false
		return
	}

	c <- true
}

/*
	Create a database
*/

func updateDatabase() {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%v", "5387b1915a694800ade3f68b7ce36eb6")

	template := `
	{
		"title": [
			{
				"text": {
					"content": "Reading List"
				}
			}
		],
		"description": [
			{
				"text":{
					"content": "Read all about it: github.com/brecht_vde/reading-list 😊"
				}
			}
		],
		"properties": {
			"Author": {
				"name": "Author",
				"type": "rich_text",
				"rich_text": {}
			},
			"Url": {
				"name": "Url",
				"type": "url",
				"url": {}
			},
			"Publishing Date": {
				"name": "Publishing/Release Date",
				"type": "date",
				"date": {}
			},
			"Read": {
				"name": "Read",
				"type": "checkbox",
				"checkbox": {}
			},
			"Name": {
				"name": "Title"
			},
			"Tags": null
		}
	}
	`

	body := strings.NewReader(template)
	request, err := http.NewRequest("PATCH", url, body)

	if err != nil {
		return
	}

	request.Header.Add("Authorization", "hidden")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Notion-Version", "2022-06-28")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return
	}

	defer response.Body.Close()
}
