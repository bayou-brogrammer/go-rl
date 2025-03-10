// This file handles the player's log.

package game

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
)

// Add this type definition somewhere in your game package
type forcedTickError struct {
	msg string
}

func (e *forcedTickError) Error() string {
	return e.msg
}

// LogEntry contains information about a log entry.
type LogEntry struct {
	Text  string      // entry text
	Color gruid.Color // color
	Dups  int         // consecutive duplicates of same message
}

func (e LogEntry) String() string {
	if e.Dups == 0 {
		return e.Text
	}
	return fmt.Sprintf("%s (%d×)", e.Text, e.Dups)
}

// Log adds an entry to the player's log.
func (g *game) log(e LogEntry) {
	if len(g.Log) > 0 {
		if g.Log[len(g.Log)-1].Text == e.Text {
			g.Log[len(g.Log)-1].Dups++
			return
		}
	}
	g.Log = append(g.Log, e)
}

// Logf adds a formatted entry to the game log.
func (g *game) Logf(format string, color gruid.Color, a ...interface{}) {
	e := LogEntry{Text: fmt.Sprintf(format, a...), Color: color}
	g.log(e)
}

// InitializeHistoryViewer creates a new pager for viewing message's history.
func (m *model) InitializeMessageViewer() {
	m.viewer = ui.NewPager(ui.PagerConfig{
		Grid: gruid.NewGrid(UIWidth, UIHeight-1),
		Box:  &ui.Box{},
	})
}
