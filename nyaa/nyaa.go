package nyaa

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type Nyaa struct {
	results map[int]*result
}

type result struct {
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName    xml.Name `xml:"item"`
	Title      string   `xml:"title"`
	Link       string   `xml:"link"`
	GUID       string   `xml:"guid"`
	PubDate    string   `xml:"pubDate"`
	Seeders    string   `xml:"nyaa:seeders"`
	Leechers   string   `xml:"nyaa:leechers"`
	Downloads  string   `xml:"nyaa:downloads"`
	InfoHash   string   `xml:"nyaa:infoHash"`
	CategoryID string   `xml:"nyaa:categoryId"`
	Category   string   `xml:"nyaa:category"`
	Size       string   `xml:"nyaa:size"`
}

func Query(url string) (*Nyaa, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var rss Rss
	err = xml.Unmarshal(bytes, &rss)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(rss.Channel.Items); i++ {
		log.Println(rss.Channel.Items[i].Title)
	}

	return nil, nil
}

func (nyaa *Nyaa) FormattedString() string {
	return ""
}

func (nyaa *Nyaa) SortByCategory() error {
	return nil
}

func (nyaa *Nyaa) SortByTitle() error {
	return nil
}

func (nyaa *Nyaa) SortBySize() error {
	return nil
}

func (nyaa *Nyaa) SortByDate() error {
	return nil
}

func (nyaa *Nyaa) SortBySeeders() error {
	return nil
}

func (nyaa *Nyaa) SortByLeechers() error {
	return nil
}

func (nyaa *Nyaa) SortByDownloads() error {
	return nil
}
