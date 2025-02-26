package main

import (
	"math/rand"
	"time"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
)

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	switch msg.(type) {
	case gruid.MsgInit:
		return m.init()
	}

	if _, ok := msg.(gruid.MsgQuit); ok {
		// md.mode = modeQuit
		return gruid.End()
	}

	m.action = action{} // reset last action information
	switch m.mode {
	case modeGameMenu:
		return m.updateGameMenu(msg)
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
	case modeMessageViewer:
		m.viewer.Update(msg)
		if m.viewer.Action() == ui.PagerQuit {
			m.mode = modeNormal
		}
		return nil
	case modeInventoryActivate, modeInventoryDrop:
		m.updateInventory(msg)
		return nil
	case modeTargeting, modeExamination:
		m.updateTargeting(msg)
		return nil
	}

	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		// Update action information on key down.
		m.updateMsgKeyDown(msg)
	case gruid.MsgMouse:
		if msg.Action == gruid.MouseMove {
			m.targ.pos = msg.P
		}
	}

	// Handle action (if any).
	return m.handleAction()
}

// updateGameMenu updates the Game Menu and switchs mode to normal after
// starting a new game or loading an old one.
func (m *model) updateGameMenu(msg gruid.Msg) gruid.Effect {
	rg := m.grid.Range().Intersect(m.grid.Range().Add(mainMenuAnchor))
	m.gameMenu.Update(rg.RelMsg(msg))

	switch m.gameMenu.Action() {
	case ui.MenuMove:
		m.info.SetText("")
	case ui.MenuInvoke:
		m.info.SetText("")
		switch m.gameMenu.Active() {
		case MenuNewGame:
			m.game = NewGame()
			m.mode = modeNormal
		case MenuContinue:
			data, err := LoadFile("save")
			if err != nil {
				m.info.SetText(err.Error())
				break
			}

			g, err := DecodeGame(data)
			if err != nil {
				m.info.SetText(err.Error())
				break
			}

			m.game = g
			m.mode = modeNormal
			// the random number generator is not saved
			m.game.Map.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
		case MenuQuit:
			return gruid.End()
		}
	case ui.MenuQuit:
		return gruid.End()
	}
	return nil
}

// updateTargeting updates targeting information in response to user input
// messages.
func (m *model) updateTargeting(msg gruid.Msg) {
	maprg := gruid.NewRange(0, LogLines, UIWidth, UIHeight-1)
	if !m.targ.pos.In(maprg) {
		m.targ.pos = m.game.ECS.PP().Add(maprg.Min)
	}

	p := m.targ.pos.Sub(maprg.Min)
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		switch msg.Key {
		case gruid.KeyArrowLeft, "h":
			p = p.Shift(-1, 0)
		case gruid.KeyArrowDown, "j":
			p = p.Shift(0, 1)
		case gruid.KeyArrowUp, "k":
			p = p.Shift(0, -1)
		case gruid.KeyArrowRight, "l":
			p = p.Shift(1, 0)
		case gruid.KeyEnter, ".":
			if m.mode == modeExamination {
				break
			}
			m.activateTarget(p)
			return
		case gruid.KeyEscape, "q":
			m.targ = targeting{}
			m.mode = modeNormal
			return
		}
		m.targ.pos = p.Add(maprg.Min)
	case gruid.MsgMouse:
		switch msg.Action {
		case gruid.MouseMove:
			m.targ.pos = msg.P
		case gruid.MouseMain:
			m.activateTarget(p)
		}
	}
}

// updateInventory handles input messages when the inventory window is open.
func (m *model) updateInventory(msg gruid.Msg) {
	// We call the Update function of the menu widget, so that we can
	// inspect information about user activity on the menu.
	m.inventory.Update(msg)

	switch m.inventory.Action() {
	case ui.MenuQuit:
		// The user requested to quit the menu.
		m.mode = modeNormal
		return
	case ui.MenuInvoke:
		// The user invoked a particular entry of the menu (either by
		// using enter or clicking on it).
		n := m.inventory.Active()

		var err error
		switch m.mode {
		case modeInventoryDrop:
			err = m.game.InventoryRemove(m.game.ECS.PlayerID, n)
		case modeInventoryActivate:
			if radius := m.game.TargetingRadius(n); radius >= 0 {
				m.targ = targeting{
					item:   n,
					pos:    m.game.ECS.PP().Shift(0, LogLines),
					radius: radius,
				}
				m.mode = modeTargeting
				return
			}
			err = m.game.InventoryActivate(m.game.ECS.PlayerID, n)
		}

		if err != nil {
			m.game.Logf("%v", ColorLogSpecial, err)
		} else {
			m.game.EndTurn()
		}
		m.mode = modeNormal
	}
}

func (m *model) updateMsgKeyDown(msg gruid.MsgKeyDown) {
	pdelta := gruid.Point{}
	m.targ.pos = gruid.Point{}
	switch msg.Key {
	case gruid.KeyArrowLeft, "h":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(-1, 0)}
	case gruid.KeyArrowDown, "j":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(0, 1)}
	case gruid.KeyArrowUp, "k":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(0, -1)}
	case gruid.KeyArrowRight, "l":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(1, 0)}
	case gruid.KeyEnter, ".":
		m.action = action{Type: ActionWait}
	case "Q":
		m.action = action{Type: ActionQuit}
	case "S":
		m.action = action{Type: ActionSave}
	case "m":
		m.action = action{Type: ActionViewMessages}
	case "i":
		m.action = action{Type: ActionInventory}
	case "d":
		m.action = action{Type: ActionDrop}
	case "g":
		m.action = action{Type: ActionPickup}
	case "x":
		m.action = action{Type: ActionExamine}
	}
}
