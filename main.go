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
	g, err := gocui.NewGui(gocui.OutputNormal)
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
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone, cursorUp); err != nil {
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
		fmt.Fprintln(v, "Search box")
		if _, err := g.SetCurrentView("search"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("sidebar", 1, 3, 20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Bookmark")
		fmt.Fprintln(v, "Recent")
		fmt.Fprintln(v, "History")
		fmt.Fprintln(v, "Settings")
	}

	if v, err := g.SetView("result", 20, 3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Result\nResult2\nResult3")
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if v.Name() == "sidebar" || v.Name() == "result" {
			ox, oy := v.Origin()
			cx, cy := v.Cursor()
			if err := v.SetCursor(cx, cy+1); err != nil {
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
			log.Printf("Cursor: cx=%d, cy=%d, ox=%d, oy=%d", cx, cy, ox, oy)
			text, _ := v.Line(cy + 1)
			if text == "" {
				text = "No text found"
			}
			log.Printf("%s\n", text)
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if v.Name() == "sidebar" || v.Name() == "result" {
			ox, oy := v.Origin()
			cx, cy := v.Cursor()
			if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
				if err := v.SetOrigin(ox, oy-1); err != nil {
					return err
				}
			}
			log.Printf("Cursor: cx=%d, cy=%d, ox=%d, oy=%d", cx, cy, ox, oy)
			text, _ := v.Line(cy - 1)
			if text == "" {
				text = "No text found"
			}
			log.Printf("%s\n", text)
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
	log.Printf("Current active view: %s", name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	active = nextIndex
	return nil
}

//func previousView(g *gocui.Gui, v *gocui.View) error {
//	previousIndex := (active - 1)
//	name := viewArr[previousIndex]
//	if previousIndex < 0 {
//		previousIndex = len(viewArr) - 1
//	}
//	previousIndex %= len(viewArr)
//
//	log.Printf("Current active view: %s", name)
//
//	if _, err := setCurrentViewOnTop(g, name); err != nil {
//		return err
//	}
//	active = previousIndex
//	return nil
//}

func submitQuery(g *gocui.Gui, v *gocui.View) error {
	userInput := v.Buffer()
	log.Printf("Input: %s", userInput)
	nyaa.Query(userInput)
	g.Update(nyaa.UpdateTable)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
