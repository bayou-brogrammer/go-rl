// This file defines the main model of the game: the Update function that
// updates the model state in response to user input, and the Draw function,
// which draws the grid.

package game

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
	"github.com/bayou-brogrammer/go-rl/game/color"
)

const (
	MenuNewGame = iota
	MenuContinue
	MenuQuit
)

type mode int

const (
	modeNormal mode = iota // normal mode

	modeQuit             // quit mode
	modeQuitConfirmation // quit confirmation mode
	modeEnd              // win or death (currently only death)

	modeInventoryActivate // inventory activate mode
	modeInventoryDrop     // inventory drop mode
	modeGameMenu          // game menu mode
	modeMessageViewer     // message viewer mode
	modeTargeting         // targeting mode (item use)
	modeExamination       // keyboad map examination mode
)

type model struct {
	grid gruid.Grid // drawing grid
	game *game      // game state
	// action    action     // UI action
	mode      mode      // UI mode
	log       *ui.Label // label for log
	status    *ui.Label // label for status
	desc      *ui.Label // label for position description
	inventory *ui.Menu  // inventory menu
	viewer    *ui.Pager // message's history viewer
	targ      targeting // targeting information
	gameMenu  *ui.Menu  // game's main menu
	info      *ui.Label // info label in main menu (for errors)

	keysNormal map[gruid.Key]action // normal mode keys
	keysTarget map[gruid.Key]action // targeting mode keys
}

// init initializes the model: widgets' initialization, and starting mode.
func (md *model) init() gruid.Effect {
	if runtime.GOOS != "js" {
		// md.mode = modeWelcom
	}

	md.initKeys()
	md.initWidgets()

	g := md.game
	md.mode = modeGameMenu

	load, err := g.LoadGame()
	if !load {
		g.InitalizeLevel()
	} else {
		g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if err != nil {
		// g.PrintStyled("Warning: could not load old saved game… starting new game.", logError)
		log.Printf("Error: %v", err)
	}

	if runtime.GOOS == "js" {
		return nil
	}
	return gruid.Sub(subSig)
}

func (md *model) initWidgets() {
	md.log = &ui.Label{}
	md.status = &ui.Label{}
	md.info = &ui.Label{}
	md.desc = &ui.Label{Box: &ui.Box{}}
	md.InitializeMessageViewer()
	entries := []ui.MenuEntry{
		MenuNewGame:  {Text: ui.Text("(N)ew game"), Keys: []gruid.Key{"N", "n"}},
		MenuContinue: {Text: ui.Text("(C)ontinue last game"), Keys: []gruid.Key{"C", "c"}},
		MenuQuit:     {Text: ui.Text("(Q)uit")},
	}
	md.gameMenu = ui.NewMenu(ui.MenuConfig{
		Grid:    gruid.NewGrid(UIWidth/2, len(entries)+2),
		Box:     &ui.Box{Title: ui.Text("Gruid Roguelike Tutorial")},
		Entries: entries,
		Style:   ui.MenuStyle{Active: gruid.Style{}.WithFg(color.ColorMenuActive)},
	})
}

func subSig(ctx context.Context, msgs chan<- gruid.Msg) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)
	select {
	case <-ctx.Done():
	case <-sig:
		msgs <- gruid.MsgQuit{}
	}
}
