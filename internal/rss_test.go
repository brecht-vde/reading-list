package internal

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	reference := rssFeed{
		Channel: rssChannel{
			Items: []rssItem{
				{
					Title:   "Test 1",
					Link:    "https://test1.local",
					PubDate: "01/18/2024",
				},
				{
					Title:   "Test 2",
					Link:    "https://test2.local",
					PubDate: "01/18/2024",
				},
			},
		},
	}

	data := []byte(`
		<rss>
			<channel>
				<item>
					<title>Test 1</title>
					<link>https://test1.local</link>
					<pubDate>01/18/2024</pubDate>
				</item>
				<item>
					<title>Test 2</title>
					<link>https://test2.local</link>
					<pubDate>01/18/2024</pubDate>
				</item>
			</channel>
		</rss>
	`)

	feed, err := parse(data)

	if err != nil {
		t.Errorf("parsing did not succeed")
	}

	if !reflect.DeepEqual(feed, reference) {
		t.Errorf(`expected: '%v', but received: '%v'`, reference, feed)
	}
}
