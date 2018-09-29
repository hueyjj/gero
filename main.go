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
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, toggleHelpPage); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, focusSearch); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, focusResult); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, focusSidebar); err != nil {
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
	}
	if v, err := g.SetView("sidebar", 0, 3, 10, maxY-4); err != nil {
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
	if v, err := g.SetView("result", 10, 3, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
	}
	if v, err := g.SetView("helpbar", -1, maxY-4, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Highlight = true
		v.Wrap = true
		v.BgColor = gocui.ColorGreen
		v.FgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintf(v, " c-c (quit) c-h (help) tab (cycle) F1-F6 (sort) c-f (search) c-r (result) c-s (menu) space (mark)")
	}
	return nil
}

func toggleHelpPage(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if g.CurrentView().Name() == "help" {
		delHelpPage(g, v)
	} else {
		if v, err := g.SetView("help", 0, 0, maxX, maxY); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Wrap = true

			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "ctrl+c", "quit program")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "ctrl+h", "toggle help page")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "tab", "cycle through views")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "ctrl+f", "focus search")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "ctrl+r", "focus result")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "ctrl+s", "focus menu")
			fmt.Fprintf(v, "\n")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "hjkl or arrow keys", "move left, down, up, right")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "End or $", "jump to end of line")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "Home or 0 (zero)", "jump to start of line")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "Page up or ctrl+u", "jump half page up")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "Page down or ctrl+d", "jump half page down")
			fmt.Fprintf(v, "\n")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "F1", "sort by comments")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "F2", "sort by date")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "F3", "sort by downloads")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "F4", "sort by leechers")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "F5", "sort by seedeers")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "F6", "sort by size")
			fmt.Fprintf(v, "\n")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "Enter (search)", "submits query")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "Enter (results)", "download highlighted torrent OR downloads marked torrent (this has precedence)")
			fmt.Fprintf(v, "\x1b[38;5;226m%20s\x1b[0m %-s\n", "Spacebar", "mark torrent to download (press enter after to download marked torrents)")
			if _, err := g.SetCurrentView("help"); err != nil {
				return err
			}
		}
	}
	return nil
}

func delHelpPage(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("help"); err != nil {
		return err
	}
	focusSearch(g, v)
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
	_, oy := v.Origin()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	length := len(line)

	maxX, _ := v.Size()
	// Shift origin
	n := length / maxX
	if n > 0 {
		if err := v.SetOrigin(maxX*n, oy); err != nil {
			return err
		}
	}
	if err := v.SetCursor(0, cy); err != nil {
		return err
	}

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
		ox--
		if err := v.SetOrigin(ox, oy); err != nil && ox > 0 {
			return err
		}
	}
	return nil
}

func cursorRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		ox++
		if err := v.SetOrigin(ox, oy); err != nil {
			return err
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

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	currentView := g.CurrentView()

	if currentView.Name() == "search" {
		// Move cursor to end of line
		line, err := currentView.Line(0)
		if err != nil {
			return nil, err
		}
		if err := currentView.SetCursor(len(line), 0); err != nil {
			return nil, err
		}
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
