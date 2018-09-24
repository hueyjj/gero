package nyaa

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jroimartin/gocui"
)

var ResultTable *Table

func UpdateTable(g *gocui.Gui) error {
	v, err := g.View("result")
	if err != nil {
		// handle error
	}
	v.Clear()

	for i := 0; i < len(ResultTable.Items); i++ {
		fmt.Fprintln(v, ResultTable.Items[i].Title)
	}
	return nil
}

func Query(searchTerm string) error {
	// Nyaa.si http request
	url := fmt.Sprintf("https://nyaa.si/?q=%s&f=0&c=0_0&page=rss", strings.TrimSpace(searchTerm))
	log.Printf("GET request (url): %s", url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// Read body into bytes
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal contents to a slice of items
	var rss Rss
	err = xml.Unmarshal(bytes, &rss)
	if err != nil {
		log.Fatal(err)
	}

	// Convert slice of items to a map
	items := make(map[int]*Item)
	for i := 0; i < len(rss.Items); i++ {
		items[i] = &rss.Items[i]
	}

	table := Table{
		Items: items,
		Sort:  Seeders,
	}

	ResultTable = &table

	return nil
}

type Sort string

const (
	Category  = Sort("Category")
	Title     = Sort("Title")
	Size      = Sort("Size")
	Date      = Sort("Date")
	Seeders   = Sort("Seeders")
	Leechers  = Sort("Leechers")
	Downloads = Sort("Downloads")
)

type Table struct {
	Items map[int]*Item
	Sort  Sort
}

func (nyaa *Table) SortByCategory() error {
	// if nyaa.Sort is Category already, then just reverse
	return nil
}

func (nyaa *Table) SortByTitle() error {
	return nil
}

func (nyaa *Table) SortBySize() error {
	return nil
}

func (nyaa *Table) SortByDate() error {
	return nil
}

func (nyaa *Table) SortBySeeders() error {
	return nil
}

func (nyaa *Table) SortByLeechers() error {
	return nil
}

func (nyaa *Table) SortByDownloads() error {
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
