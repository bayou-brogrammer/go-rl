// This file defines the main model of the game: the Update function that
// updates the model state in response to user input, and the Draw function,
// which draws the grid.

package main

import (
	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
	"github.com/bayou-brogrammer/go-rl/game/color"
)

type mode int

const (
	modeNormal            mode = iota // normal mode
	modeEnd                           // win or death (currently only death)
	modeInventoryActivate             // inventory activate mode
	modeInventoryDrop                 // inventory drop mode
	modeGameMenu                      // game menu mode
	modeMessageViewer                 // message viewer mode
	modeTargeting                     // targeting mode (item use)
	modeExamination                   // keyboad map examination mode
)

type model struct {
	grid      gruid.Grid // drawing grid
	game      *game      // game state
	action    action     // UI action
	mode      mode       // UI mode
	log       *ui.Label  // label for log
	status    *ui.Label  // label for status
	desc      *ui.Label  // label for position description
	inventory *ui.Menu   // inventory menu
	viewer    *ui.Pager  // message's history viewer
	targ      targeting  // targeting information
	gameMenu  *ui.Menu   // game's main menu
	info      *ui.Label  // info label in main menu (for errors)
}

// init initializes the model: widgets' initialization, and starting mode.
func (m *model) init() gruid.Effect {
	m.log = &ui.Label{}
	m.status = &ui.Label{}
	m.info = &ui.Label{}
	m.desc = &ui.Label{Box: &ui.Box{}}
	m.InitializeMessageViewer()
	m.mode = modeGameMenu
	entries := []ui.MenuEntry{
		MenuNewGame:  {Text: ui.Text("(N)ew game"), Keys: []gruid.Key{"N", "n"}},
		MenuContinue: {Text: ui.Text("(C)ontinue last game"), Keys: []gruid.Key{"C", "c"}},
		MenuQuit:     {Text: ui.Text("(Q)uit")},
	}
	m.gameMenu = ui.NewMenu(ui.MenuConfig{
		Grid:    gruid.NewGrid(UIWidth/2, len(entries)+2),
		Box:     &ui.Box{Title: ui.Text("Gruid Roguelike Tutorial")},
		Entries: entries,
		Style:   ui.MenuStyle{Active: gruid.Style{}.WithFg(color.ColorMenuActive)},
	})
	return nil
}

type targeting struct {
	pos    gruid.Point
	item   int // item to use after selecting target
	radius int
}

const (
	MenuNewGame = iota
	MenuContinue
	MenuQuit
)
