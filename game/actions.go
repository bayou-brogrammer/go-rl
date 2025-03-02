// This file manages actions resulting from user input.

package game

import (
	"log"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
	"github.com/bayou-brogrammer/go-rl/game/color"
)

////////////////////
// action errors
////////////////////

type actionError int

const (
	actionErrorUnknown actionError = iota
)

func (e actionError) Error() string {
	switch e {
	case actionErrorUnknown:
		return "unknown action"
	}
	return ""
}

////////////////////
// action
////////////////////

type actionType int

// action represents information relevant to the last UI action performed.
type action struct {
	Type  actionType  // kind of action (movement, quitting, ...)
	Delta gruid.Point // direction for ActionBump
}

// These constants represent the possible UI actions.
const (
	ActionNone actionType = iota
	ActionW
	ActionS
	ActionN
	ActionE
	ActionBump         // bump request (attack or movement)
	ActionDrop         // menu to drop an inventory item
	ActionInventory    // inventory menu to use an item
	ActionPickup       // pickup an item on the ground
	ActionWait         // wait a turn
	ActionQuit         // quit the game (without saving)
	ActionSave         // save the game
	ActionViewMessages // view history messages
	ActionExamine      // examine map
)

func (md *model) normalModeAction(action action) (again bool, eff gruid.Effect, err error) {
	switch action.Type {
	case ActionNone:
		again = true
		err = actionErrorUnknown

	case ActionBump:
		np := md.game.ECS.PP().Add(action.Delta)
		md.game.Bump(np)

	case ActionDrop:
		md.OpenInventory("Drop item")
		md.mode = modeInventoryDrop

	case ActionInventory:
		md.OpenInventory("Use item")
		md.mode = modeInventoryActivate

	case ActionPickup:
		md.game.PickupItem()

	case ActionWait:
		// md.game.EndTurn()

	case ActionSave:
		data, err := EncodeGame(md.game)
		if err == nil {
			err = SaveFile("save", data)
		}
		if err != nil {
			md.game.Logf("Could not save game.", color.ColorLogSpecial)
			log.Printf("could not save game: %v", err)
			break
		}

	case ActionQuit:
		again = true
		// // Remove any previously saved files (if any).
		// RemoveDataFile("save")
		md.Quit()

	case ActionViewMessages:
		md.mode = modeMessageViewer
		lines := []ui.StyledText{}
		for _, e := range md.game.Log {
			st := gruid.Style{}
			st.Fg = e.Color
			lines = append(lines, ui.NewStyledText(e.String(), st))
		}
		md.viewer.SetLines(lines)

	case ActionExamine:
		md.mode = modeExamination
		md.targ.pos = md.game.ECS.PP().Shift(0, LogLines)
	}

	if md.game.ECS.PlayerDied() {
		md.game.Logf("You died -- press “q” or escape to quit", color.ColorLogSpecial)
		md.mode = modeEnd
		return false, nil, nil
	}

	if err != nil {
		again = true
	}
	return again, eff, err
}

// Bump moves the player to a given position and updates FOV information,
// or attacks if there is a monster.
func (g *game) Bump(to gruid.Point) {
	if !g.Map.Walkable(to) {
		return
	}
	if i := g.ECS.MonsterAt(to); g.ECS.Alive(i) {
		// We show a message to standard error. Later in the tutorial,
		// we'll put a message in the UI instead.
		g.BumpAttack(g.ECS.PlayerID, i)
		g.EndTurn()
		return
	}
	// We move the player to the new destination.
	g.ECS.MovePlayer(to)
	g.EndTurn()
}

// PickupItem takes an item on the floor.
func (g *game) PickupItem() {
	pp := g.ECS.PP()
	for i, p := range g.ECS.Positions {
		if p != pp {
			// Skip entities whose position is diffferent than the
			// player's.
			continue
		}
		err := g.InventoryAdd(g.ECS.PlayerID, i)
		if err != nil {
			if err.Error() == ErrNoShow {
				// Happens for example if the current entity is
				// not a consumable.
				continue
			}
			g.Logf("Could not pickup: %v", color.ColorLogSpecial, err)
			return
		}
		g.Logf("You pickup %v", color.ColorLogItemUse, g.ECS.Name[i])
		g.EndTurn()
		return
	}
}

// OpenInventory opens the inventory and allows the player to select an item.
func (m *model) OpenInventory(title string) {
	inv := m.game.ECS.Inventory[m.game.ECS.PlayerID]
	// We build a list of entries.
	entries := []ui.MenuEntry{}
	r := 'a'
	for _, it := range inv.Items {
		name := m.game.ECS.Name[it]
		entries = append(entries, ui.MenuEntry{
			Text: ui.Text(string(r) + " - " + name),
			// allow to use the character r to select the entry
			Keys: []gruid.Key{gruid.Key(r)},
		})
		r++
	}
	// We create a new menu widget for the inventory window.
	m.inventory = ui.NewMenu(ui.MenuConfig{
		Grid:    gruid.NewGrid(40, MapHeight),
		Box:     &ui.Box{Title: ui.Text(title)},
		Entries: entries,
	})
}

func (md *model) Quit() {
	// md.g.PrintStyled("Do you really want to quit without saving? [y/N]", logConfirm)
	md.mode = modeQuitConfirmation
}
