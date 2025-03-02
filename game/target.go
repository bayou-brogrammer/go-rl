package game

import (
	"codeberg.org/anaseto/gruid"
	"github.com/bayou-brogrammer/go-rl/game/color"
	"github.com/bayou-brogrammer/go-rl/game/logerror"
)

type targeting struct {
	pos    gruid.Point
	item   int // item to use after selecting target
	radius int
}

// InventoryActivateWithTarget uses a given item from the inventory, with
// an optional target.
func (g *game) InventoryActivateWithTarget(actor, n int, targ *gruid.Point) error {
	inv := g.ECS.Inventory[actor]
	if len(inv.Items) <= n {
		return logerror.New("Empty slot.")
	}

	i := inv.Items[n]
	switch e := g.ECS.Entities[i].(type) {
	case Consumable:
		err := e.Activate(g, itemAction{Actor: actor, Target: targ})
		if err != nil {
			return err
		}
	}
	// Put the last item on the previous one: this could be improved,
	// sorting elements in a certain way, or moving elements as necessary
	// to preserve current order.
	inv.Items[n] = inv.Items[len(inv.Items)-1]
	inv.Items = inv.Items[:len(inv.Items)-1]
	return nil
}

func (m *model) activateTarget(p gruid.Point) {
	err := m.game.InventoryActivateWithTarget(m.game.ECS.PlayerID, m.targ.item, &p)
	if err != nil {
		m.game.Logf("%v", color.ColorLogSpecial, err)
	} else {
		m.game.EndTurn()
	}

	m.mode = modeNormal
	m.targ = targeting{}
}

// updateTargeting updates targeting information in response to user input
// messages.
func (md *model) updateTargeting(msg gruid.Msg) {
	maprg := gruid.NewRange(0, LogLines, UIWidth, UIHeight-1)
	if !md.targ.pos.In(maprg) {
		md.targ.pos = md.game.ECS.PP().Add(maprg.Min)
	}

	p := md.targ.pos.Sub(maprg.Min)
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
			if md.mode == modeExamination {
				break
			}
			md.activateTarget(p)
			return
		case gruid.KeyEscape, "q":
			md.targ = targeting{}
			md.mode = modeNormal
			return
		}
		md.targ.pos = p.Add(maprg.Min)
	case gruid.MsgMouse:
		switch msg.Action {
		case gruid.MouseMove:
			md.targ.pos = msg.P
		case gruid.MouseMain:
			md.activateTarget(p)
		}
	}
}
