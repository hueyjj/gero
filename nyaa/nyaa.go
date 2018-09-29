package nyaa

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jroimartin/gocui"
)

type SortMethod string
type Order string

const (
	Comments   = SortMethod("comments")
	Size       = SortMethod("size")
	Date       = SortMethod("id")
	Seeders    = SortMethod("seeders")
	Leechers   = SortMethod("leechers")
	Downloads  = SortMethod("downloads")
	Ascending  = Order("asc")
	Descending = Order("desc")
)

var (
	table *Table
	sort  = Date
	order = Descending
	//page           = 0
	lastSearchTerm = ""
	markedItems    []int
)

func UpdateTable(g *gocui.Gui) error {
	v, err := g.View("result")
	if err != nil {
		// handle error
	}
	v.Clear()

	for i, item := range table.Items {
		str := fmt.Sprintf("\x1b[38;5;7m%3d\x1b[0m ", i)

		// Check if marked
		if TorrentIsMarked(i) {
			str = fmt.Sprintf("\x1b[48;5;226m\x1b[30m%3d\x1b[0m ", i)
		}

		spaces := maxLength("seeders")
		str += fmt.Sprintf("\x1b[38;5;34m%*s\x1b[0m/", spaces, item.Seeders)

		spaces = maxLength("leechers")
		str += fmt.Sprintf("\x1b[38;5;196m%-*s\x1b[0m ", spaces, item.Leechers)

		spaces = maxLength("downloads")
		str += fmt.Sprintf("\x1b[38;5;6m%*s\x1b[0m ", spaces, item.Downloads)

		pubDate := formatDate(item.PubDate)
		spaces = maxLength("pubdate")
		str += fmt.Sprintf("\x1b[48;5;66m\x1b[8m%*s\x1b[0m ", spaces, pubDate)

		spaces = maxLength("size")
		str += fmt.Sprintf("\x1b[38;5;6m%*s\x1b[0m ", spaces, item.Size)

		str += fmt.Sprintf("\x1b[38;5;7m%s\n\x1b[0m", item.Title)

		fmt.Fprintf(v, str)
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

	// Unmarshal contents
	var rss Rss
	err = xml.Unmarshal(content, &rss)
	if err != nil {
		log.Fatal(err)
	}

	var items []*Item
	for i := 0; i < len(rss.Items); i++ {
		items = append(items, &rss.Items[i])
	}

	table = &Table{
		Items: items,
	}
	lastSearchTerm = searchTerm

	return nil
}

//func NextPage(g *gocui.Gui) {
//	page++
//	Query(lastSearchTerm)
//	UpdateTable(g)
//}
//
//func PrevPage(g *gocui.Gui) {
//	if page--; page < 0 {
//		return
//	}
//	Query(lastSearchTerm)
//	UpdateTable(g)
//}

func OpenTorrent(index int) error {
	torrent, err := downloadTorrent(index)
	if err != nil {
		return err
	}
	log.Printf("Firing xdg-open to open %s\n", torrent)
	if runtime.GOOS == "windows" {
		exec.Command("start", torrent).Start()
	} else if runtime.GOOS == "linux" {
		exec.Command("xdg-open", torrent).Start()
	} else if runtime.GOOS == "darwin" {
		exec.Command("open", torrent).Start()
	}
	return nil
}

func OpenMarkedTorrents() error {
	for _, index := range markedItems {
		OpenTorrent(index)
	}
	return nil
}

func TorrentsMarked() bool {
	if len(markedItems) > 0 {
		return true
	}
	return false
}

func MarkTorrent(index int) error {
	markedItems = append(markedItems, index)
	return nil
}

func UnmarkTorrent(index int) error {
	i := 0
	for n, item := range markedItems {
		if index == item {
			i = n
		}
	}
	markedItems = append(markedItems[:i], markedItems[i+1:]...)
	return nil
}

func MarkedTorrentsRemoveAll() {
	markedItems = []int{}
}

func TorrentIsMarked(index int) bool {
	for _, i := range markedItems {
		if index == i {
			return true
		}
	}
	return false
}

func Sort(g *gocui.Gui, sortMethod SortMethod) {
	sort = sortMethod
	Query(lastSearchTerm)
	UpdateTable(g)
}

func downloadTorrent(index int) (string, error) {
	if 0 <= index && index < len(table.Items) {
		log.Printf("Downloading torrent (index=%d): %s", index, table.Items[index].Title)
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
		torrent := filepath.Join(dir, name)
		if err := ioutil.WriteFile(torrent, content, 0666); err != nil {
			log.Fatal(err)
		}
		log.Printf("Download complete for torrent (index=%d): %s", index, table.Items[index].Title)
		return torrent, nil
	}
	return "", fmt.Errorf("Index %d out of bounds [0-%d]", index, len(table.Items))
}

func maxLength(str string) int {
	max := 0
	for _, item := range table.Items {
		if str == "seeders" {
			if len(item.Seeders) > max {
				max = len(item.Seeders)
			}
		} else if str == "leechers" {
			if len(item.Leechers) > max {
				max = len(item.Leechers)
			}
		} else if str == "downloads" {
			if len(item.Downloads) > max {
				max = len(item.Downloads)
			}
		} else if str == "pubdate" {
			if len(item.PubDate) > max {
				max = len(formatDate(item.PubDate))
			}
		} else if str == "size" {
			if len(item.Size) > max {
				max = len(item.Size)
			}
		}
	}
	return max
}

func formatDate(pubDate string) string {
	s := strings.Split(pubDate, " ")
	if len(s) != 6 {
		return ""
	} else {
		time := strings.Split(s[4], ":")
		var hour, min string
		if len(time) != 3 {
			hour = "0"
			min = "0"
		} else {
			hour = time[0]
			min = time[1]
		}
		return fmt.Sprintf("%s-%s-%s %s:%s", s[3], s[2], s[1], hour, min)
	}
}

type Table struct {
	Items []*Item
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
