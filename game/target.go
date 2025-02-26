package main

import (
	"errors"

	"codeberg.org/anaseto/gruid"
	"github.com/bayou-brogrammer/go-rl/game/color"
)

// InventoryActivateWithTarget uses a given item from the inventory, with
// an optional target.
func (g *game) InventoryActivateWithTarget(actor, n int, targ *gruid.Point) error {
	inv := g.ECS.Inventory[actor]
	if len(inv.Items) <= n {
		return errors.New("Empty slot.")
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
