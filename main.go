package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/hueyjj/gero/nyaa"
	"github.com/jroimartin/gocui"
)

var (
	viewArr = []string{"search", "sidebar", "result"}
	active  = 0
)

func main() {
	// Get current user
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Setup logger
	logSaveDir := filepath.Join(user.HomeDir, "Downloads", "gero.log")
	f, err := os.OpenFile(logSaveDir, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("Program started")
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Highlight = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '1', gocui.ModAlt, focusSearch); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '2', gocui.ModAlt, focusResult); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '3', gocui.ModAlt, focusSidebar); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, cursorLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, cursorRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, cursorEndOfLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, cursorStartOfLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, cursorPageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, cursorPageDown); err != nil {
		return err
	}
	//if err := g.SetKeybinding("result", 'n', gocui.ModNone, nextPage); err != nil {
	//	return err
	//}
	//if err := g.SetKeybinding("result", 'N', gocui.ModNone, prevPage); err != nil {
	//	return err
	//}
	if err := g.SetKeybinding("result", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", 'h', gocui.ModNone, cursorLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", 'l', gocui.ModNone, cursorRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", '0', gocui.ModNone, cursorStartOfLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", '$', gocui.ModNone, cursorEndOfLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyCtrlU, gocui.ModNone, cursorPageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyCtrlD, gocui.ModNone, cursorPageDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyF1, gocui.ModNone, sortByComments); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyF2, gocui.ModNone, sortByDate); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyF3, gocui.ModNone, sortByDownloads); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyF4, gocui.ModNone, sortByLeechers); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyF5, gocui.ModNone, sortBySeeders); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyF6, gocui.ModNone, sortBySize); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeyEnter, gocui.ModNone, openTorrent); err != nil {
		return err
	}
	if err := g.SetKeybinding("result", gocui.KeySpace, gocui.ModNone, markTorrent); err != nil {
		return err
	}
	if err := g.SetKeybinding("sidebar", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("sidebar", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("sidebar", 'h', gocui.ModNone, cursorLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("sidebar", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("sidebar", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("sidebar", 'l', gocui.ModNone, cursorRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("search", gocui.KeyEnter, gocui.ModNone, submitQuery); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("search", maxX/2-20, 0, maxX/2+20, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		if _, err := g.SetCurrentView("search"); err != nil {
			return err
		}
		fmt.Fprintf(v, "naruto")
	}
	if v, err := g.SetView("sidebar", 0, 3, 10, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Search")
		fmt.Fprintln(v, "Bookmark")
		fmt.Fprintln(v, "Recent")
		fmt.Fprintln(v, "History")
		fmt.Fprintln(v, "Settings")
	}
	if v, err := g.SetView("result", 10, 3, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Result\nResult2\nResult3")
	}
	if v, err := g.SetView("helpbar", -1, maxY-3, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Highlight = true
		v.BgColor = gocui.ColorGreen
		v.FgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintf(v, " h (help)")
	}
	return nil
}

//func nextPage(g *gocui.Gui, v *gocui.View) error {
//	nyaa.NextPage(g)
//	g.Update(nyaa.UpdateTable)
//	return nil
//}
//
//func prevPage(g *gocui.Gui, v *gocui.View) error {
//	nyaa.PrevPage(g)
//	g.Update(nyaa.UpdateTable)
//	return nil
//}

func cursorEndOfLine(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	//ox, _ := v.Origin()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	cursorStartOfLine(g, v)
	v.MoveCursor(len(line), 0, false)
	return nil
}

func cursorStartOfLine(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	if err := v.SetCursor(0, cy); err != nil {
		return err
	}
	_, oy := v.Origin()
	if err := v.SetOrigin(0, oy); err != nil {
		return err
	}
	return nil
}

func cursorPageUp(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	ox, oy := v.Origin()
	_, maxY := v.Size()
	if cy+oy-maxY/2 < 0 {
		return nil
	} else if oy > 0 {
		nextOy := oy - maxY/2
		if nextOy < 0 {
			nextOy = 0
		}
		if err := v.SetOrigin(ox, nextOy); err != nil {
			geroError := fmt.Errorf("cy=%d ox=%d oy=%d maxY=%d", cy, ox, oy, maxY)
			return fmt.Errorf("cursorPageUp: %v %v", err, geroError)
		}
	}
	return nil
}

func cursorPageDown(g *gocui.Gui, v *gocui.View) error {
	bufLines := v.BufferLines()
	_, cy := v.Cursor()
	ox, oy := v.Origin()
	_, maxY := v.Size()
	if cy+oy+maxY/2 > len(bufLines)-2 {
		return nil
	} else if err := v.SetOrigin(ox, oy+maxY/2); err != nil {
		return err
	}
	return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		cx--
		if err := v.SetCursor(cx, cy); err != nil && ox > 0 {
			ox--
			if err := v.SetOrigin(ox, oy); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		bufLines := v.BufferLines()
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		// Checking if we were at the end of line
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		isEndOfLine := false
		if cx+ox >= len(line) {
			isEndOfLine = true
		}

		cy++
		// Don't do anything if there's nothing left after current point
		if cy+oy > len(bufLines)-2 {
			return nil
		} else if err := v.SetCursor(cx, cy); err != nil {
			oy++
			// Move cursor down
			if err := v.SetOrigin(ox, oy); err != nil {
				return err
			}
		}

		// Move cursor to end of line
		if isEndOfLine {
			cursorEndOfLine(g, v)
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		// Checking if we were at the end of line
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		isEndOfLine := false
		if cx+ox >= len(line) {
			isEndOfLine = true
		}

		cy--
		if err := v.SetCursor(cx, cy); err != nil && oy > 0 {
			oy--
			if err := v.SetOrigin(ox, oy); err != nil {
				return err
			}
		}

		// Move cursor to end of line
		if isEndOfLine {
			cursorEndOfLine(g, v)
		}
	}
	return nil
}

func cursorRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		cx++
		if err := v.SetCursor(cx, cy); err != nil {
			ox++
			if err := v.SetOrigin(ox, oy); err != nil {
				return err
			}
		}
	}
	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]
	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	log.Printf("Current active view: %s", name)
	active = nextIndex
	return nil
}

func focusSearch(g *gocui.Gui, v *gocui.View) error {
	var nextIndex int
	for i, name := range viewArr {
		if name == "search" {
			nextIndex = i
		}
	}
	if _, err := setCurrentViewOnTop(g, "search"); err != nil {
		return err
	}
	active = nextIndex
	return nil
}

func focusResult(g *gocui.Gui, v *gocui.View) error {
	var nextIndex int
	for i, name := range viewArr {
		if name == "result" {
			nextIndex = i
		}
	}
	if _, err := setCurrentViewOnTop(g, "result"); err != nil {
		return err
	}
	active = nextIndex
	return nil
}

func focusSidebar(g *gocui.Gui, v *gocui.View) error {
	var nextIndex int
	for i, name := range viewArr {
		if name == "sidebar" {
			nextIndex = i
		}
	}
	if _, err := setCurrentViewOnTop(g, "sidebar"); err != nil {
		return err
	}
	active = nextIndex
	return nil
}

//func previousView(g *gocui.Gui, v *gocui.View) error {
//	previousIndex := (active - 1)
//	name := viewArr[previousIndex]
//
//	// Reset
//	if previousIndex < 0 {
//		previousIndex = len(viewArr) - 1
//	}
//
//	previousIndex %= len(viewArr)
//	if _, err := setCurrentViewOnTop(g, name); err != nil {
//		return err
//	}
//
//	log.Printf("Current active view: %s", name)
//	active = previousIndex
//	return nil
//}

func submitQuery(g *gocui.Gui, v *gocui.View) error {
	userInput := v.Buffer()
	log.Printf("Search term: %s", userInput)
	nyaa.Query(userInput)
	g.Update(nyaa.UpdateTable)
	return nil
}

func openTorrent(g *gocui.Gui, v *gocui.View) error {
	if nyaa.TorrentsMarked() {
		// Open all marked torrents
		nyaa.OpenMarkedTorrents()
		nyaa.MarkedTorrentsRemoveAll()
		g.Update(nyaa.UpdateTable)
	} else {
		// Open torreunt under cursor
		_, oy := v.Origin()
		_, cy := v.Cursor()
		nyaa.OpenTorrent(oy + cy)
	}
	return nil
}

func markTorrent(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	index := oy + cy
	if nyaa.TorrentIsMarked(index) {
		nyaa.UnmarkTorrent(index)
	} else {
		nyaa.MarkTorrent(index)
	}
	g.Update(nyaa.UpdateTable)
	return nil
}

func sortByComments(g *gocui.Gui, v *gocui.View) error {
	nyaa.Sort(g, nyaa.Comments)
	return nil
}

func sortBySize(g *gocui.Gui, v *gocui.View) error {
	nyaa.Sort(g, nyaa.Size)
	return nil
}

func sortByDate(g *gocui.Gui, v *gocui.View) error {
	nyaa.Sort(g, nyaa.Date)
	return nil
}

func sortBySeeders(g *gocui.Gui, v *gocui.View) error {
	nyaa.Sort(g, nyaa.Seeders)
	return nil
}

func sortByLeechers(g *gocui.Gui, v *gocui.View) error {
	nyaa.Sort(g, nyaa.Leechers)
	return nil
}

func sortByDownloads(g *gocui.Gui, v *gocui.View) error {
	nyaa.Sort(g, nyaa.Downloads)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
