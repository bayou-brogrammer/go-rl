package game

import (
	"fmt"
	"time"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
	"github.com/bayou-brogrammer/go-rl/game/color"
)

func (md *model) Update(msg gruid.Msg) gruid.Effect {
	switch msg.(type) {
	case gruid.MsgInit:
		return md.init()
	}

	if _, ok := msg.(gruid.MsgQuit); ok {
		// md.mode = modeQuit
		return gruid.End()
	}

	return md.update(msg)
}

func (md *model) update(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch md.mode {
	case modeQuit:
		return nil

	case modeEnd:
		switch msg := msg.(type) {
		case gruid.MsgKeyDown:
			switch msg.Key {
			case "q", gruid.KeyEscape:
				// You died: quit on "q" or "escape"
				return gruid.End()
			}
		}
		return nil

	case modeNormal:
		eff = md.updateNormal(msg)

	case modeGameMenu:
		return md.updateGameMenu(msg)

	case modeMessageViewer:
		md.viewer.Update(msg)
		if md.viewer.Action() == ui.PagerQuit {
			md.mode = modeNormal
		}
		return nil
	case modeInventoryActivate, modeInventoryDrop:
		md.updateInventory(msg)
		return nil
	case modeTargeting, modeExamination:
		md.updateTargeting(msg)
		return nil

	}

	return eff
}

func (md *model) updateNormal(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		eff = md.updateKeyDown(msg)
	case gruid.MsgMouse:
		eff = md.updateMouse(msg)
	}
	return eff
}

func (md *model) updateMouse(msg gruid.MsgMouse) gruid.Effect {
	if msg.Action == gruid.MouseMove {
		md.targ.pos = msg.P
	}

	return nil
}

func (md *model) updateKeyDown(msg gruid.MsgKeyDown) gruid.Effect {
	again, eff, err := md.normalModeKeyDown(msg.Key, msg.Mod&gruid.ModShift != 0)
	if err != nil {
		// md.game.Print(err.Error())
	}
	if again {
		return eff
	}

	return md.game.EndTurn()
}

func (md *model) normalModeKeyDown(key gruid.Key, shift bool) (again bool, eff gruid.Effect, err error) {
	action := md.keysNormal[key]
	again, eff, err = md.normalModeAction(action)
	if _, ok := err.(actionError); ok {
		// Create a special error type that will force a new tick
		err = &forcedTickError{
			msg: fmt.Sprintf("Key '%s' does nothing. Type ? for help.", key),
		}
	}
	return again, eff, err
}

// updateGameMenu updates the Game Menu and switchs mode to normal after
// starting a new game or loading an old one.
func (md *model) updateGameMenu(msg gruid.Msg) gruid.Effect {
	rg := md.grid.Range().Intersect(md.grid.Range().Add(mainMenuAnchor))
	md.gameMenu.Update(rg.RelMsg(msg))

	switch md.gameMenu.Action() {
	case ui.MenuMove:
		md.info.SetText("")
	case ui.MenuInvoke:
		md.info.SetText("")
		switch md.gameMenu.Active() {
		case MenuNewGame:
			md.game.initalizeFirstLevel()
			md.mode = modeNormal
		case MenuContinue:
			data, err := LoadFile("save")
			if err != nil {
				md.info.SetText(err.Error())
				break
			}

			g, err := DecodeGame(data)
			if err != nil {
				md.info.SetText(err.Error())
				break
			}

			md.game = g
			md.mode = modeNormal
			// the random number generator is not saved
			md.game.Map.SeedRand(time.Now().UnixNano())
		case MenuQuit:
			return gruid.End()
		}
	case ui.MenuQuit:
		return gruid.End()
	}
	return nil
}

// updateInventory handles input messages when the inventory window is open.
func (md *model) updateInventory(msg gruid.Msg) {
	// We call the Update function of the menu widget, so that we can
	// inspect information about user activity on the menu.
	md.inventory.Update(msg)

	switch md.inventory.Action() {
	case ui.MenuQuit:
		// The user requested to quit the menu.
		md.mode = modeNormal
		return
	case ui.MenuInvoke:
		// The user invoked a particular entry of the menu (either by
		// using enter or clicking on it).
		n := md.inventory.Active()

		var err error
		switch md.mode {
		case modeInventoryDrop:
			err = md.game.InventoryRemove(md.game.ECS.PlayerID, n)
		case modeInventoryActivate:
			if radius := md.game.TargetingRadius(n); radius >= 0 {
				md.targ = targeting{
					item:   n,
					pos:    md.game.ECS.PP().Shift(0, LogLines),
					radius: radius,
				}
				md.mode = modeTargeting
				return
			}
			err = md.game.InventoryActivate(md.game.ECS.PlayerID, n)
		}

		if err != nil {
			md.game.Logf("%v", color.ColorLogSpecial, err)
		} else {
			md.game.EndTurn()
		}
		md.mode = modeNormal
	}
}
