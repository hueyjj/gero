package nyaa

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jroimartin/gocui"
)

type Sort string
type Order string

const (
	Comments   = Sort("comments")
	Size       = Sort("size")
	Date       = Sort("id")
	Seeders    = Sort("seeders")
	Leechers   = Sort("leechers")
	Downloads  = Sort("downloads")
	Ascending  = Order("asc")
	Descending = Order("desc")
)

var (
	sort  = Comments
	order = Descending
	page  = 0
	table *Table
)

func UpdateTable(g *gocui.Gui) error {
	v, err := g.View("result")
	if err != nil {
		// handle error
	}
	v.Clear()

	for i := 0; i < len(table.Items); i++ {
		fmt.Fprintf(v, "%3d %s\n", i, table.Items[i].Title)
	}
	return nil
}

func Query(searchTerm string) error {
	// Nyaa.si http request
	url := fmt.Sprintf("https://nyaa.si/?f=0&c=0_0&q=%s&s=%s&o=%s&page=rss",
		strings.TrimSpace(searchTerm),
		sort,
		order,
	)
	log.Printf("GET request (url): %s", url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GET request complete.")

	// Read body into content (bytes)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal contents to a slice of items
	var rss Rss
	err = xml.Unmarshal(content, &rss)
	if err != nil {
		log.Fatal(err)
	}

	// Convert slice of items to a map
	var items []*Item
	for i := 0; i < len(rss.Items); i++ {
		items = append(items, &rss.Items[i])
	}

	table = &Table{
		Items:      items,
		SortMethod: Date,
	}

	return nil
}

func DownloadTorrent(index int) error {
	log.Printf("Downloading torrent (index=%d): %s", index, table.Items[index].Title)
	if 0 <= index && index < len(table.Items) {
		res, err := http.Get(table.Items[index].Link)
		if err != nil {
			log.Fatal(err)
		}
		content, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		dir, err := ioutil.TempDir("", "gero-torrents")
		if err != nil {
			log.Fatal(err)
		}

		splitLink := strings.SplitAfter(table.Items[index].Link, "/")
		name := splitLink[len(splitLink)-1]
		tmpTorrent := filepath.Join(dir, name)
		if err := ioutil.WriteFile(tmpTorrent, content, 0666); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

type Table struct {
	Items      []*Item
	SortMethod Sort
}

/**func (nyaa *Table) Sort() {
	switch nyaa.SortMethod {
	case Comments:
		nyaa.SortByComments()
	case Size:
		nyaa.SortBySize()
	case Date:
		nyaa.SortByDate()
	case Seeders:
		nyaa.SortBySeeders()
	case Leechers:
		nyaa.SortByLeechers()
	case Downloads:
		nyaa.SortByDownloads()
	default:
		nyaa.SortByDate()
	}
}
*/

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
