package nyaa

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
	for i := 0; i < len(rss.Items); i++ {
		log.Print(rss.Items[i].String())
	}

	return nil, nil
}

type Nyaa struct {
	results map[int]*result
}

type result struct {
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

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Items   []Item   `xml:"channel>item"`
}

type Item struct {
	XMLName    xml.Name `xml:"item"`
	Title      string   `xml:"title"`
	Link       string   `xml:"link"`
	GUID       string   `xml:"guid"`
	PubDate    string   `xml:"pubDate"`
	Seeders    string   `xml:"https://nyaa.si/xmlns/nyaa seeders"`
	Leechers   string   `xml:"https://nyaa.si/xmlns/nyaa leechers"`
	Downloads  string   `xml:"https://nyaa.si/xmlns/nyaa downloads"`
	InfoHash   string   `xml:"https://nyaa.si/xmlns/nyaa infoHash"`
	CategoryID string   `xml:"https://nyaa.si/xmlns/nyaa categoryId"`
	Category   string   `xml:"https://nyaa.si/xmlns/nyaa category"`
	Size       string   `xml:"https://nyaa.si/xmlns/nyaa size"`
}

func (item *Item) String() string {
	fmt.Println(item.Seeders)
	str := fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s\n",
		item.Title,
		item.Link,
		item.GUID,
		item.PubDate,
		item.Seeders,
		item.Leechers,
		item.Downloads,
		item.InfoHash,
		item.CategoryID,
		item.Category,
		item.Size,
	)
	return str
}
