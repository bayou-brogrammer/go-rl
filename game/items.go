// This file describes item entities.

package game

import (
	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
	"github.com/bayou-brogrammer/go-rl/game/color"
	"github.com/bayou-brogrammer/go-rl/game/logerror"
)

// Consumable describes a consumable item, like a potion.
type Consumable interface {
	// Activate makes use of an item using a specific action. It returns
	// an error if the consumable could not be activated.
	Activate(g *game, a itemAction) error
}

// itemAction describes information relative to usage of an item: which
// actor does the action, and whether the action has a particular target
// position.
type itemAction struct {
	Actor  int          // entity doing the action
	Target *gruid.Point // optional target
}

// HealingPotion describes a potion that heals of a given amount.
type HealingPotion struct {
	Amount int
}

func (pt *HealingPotion) Activate(g *game, a itemAction) error {
	st := g.ECS.Stats[a.Actor]
	if st == nil {
		// should not happen in practice
		return logerror.Errorf("%s cannot use healing potions.", g.ECS.Name[a.Actor])
	}
	hp := st.Heal(pt.Amount)
	if hp <= 0 {
		return logerror.New("Your health is already full.")
	}
	if a.Actor == g.ECS.PlayerID {
		g.Logf("You regained %d HP", color.ColorLogItemUse, hp)
	}
	return nil
}

// LightningScroll is an item that can be invoked to strike the closest enemy
// within a particular range.
type LightningScroll struct {
	Range  int
	Damage int
}

func (sc *LightningScroll) Activate(g *game, a itemAction) error {
	target := -1
	minDist := sc.Range + 1
	for i := range g.ECS.Stats {
		p := g.ECS.Positions[i]
		if i == a.Actor || g.ECS.Dead(i) || !g.InFOV(p) {
			continue
		}
		dist := paths.DistanceManhattan(p, g.ECS.Positions[a.Actor])
		if dist < minDist {
			target = i
			minDist = dist
		}
	}
	if target < 0 {
		return logerror.New("No enemy within range.")
	}
	g.Logf("A lightning bolt strikes %v.", color.ColorLogItemUse, g.ECS.GetName(target))
	g.ECS.Stats[target].HP -= sc.Damage
	return nil
}

// Targetter describes consumables (or other kind of activables) that need
// a target in order to be used.
type Targetter interface {
	// TargetingRadius returns the radius of the affected area around the
	// target.
	TargetingRadius() int
}

// ConfusionScroll is an item that can be invoked to confuse an enemy.
type ConfusionScroll struct {
	Turns int
}

func (sc *ConfusionScroll) Activate(g *game, a itemAction) error {
	if a.Target == nil {
		return logerror.New("You have to chose a target.")
	}
	p := *a.Target
	if !g.InFOV(p) {
		return logerror.New("You cannot target what you cannot see.")
	}
	if p == g.ECS.PP() {
		return logerror.New("You cannot confuse yourself.")
	}
	i := g.ECS.MonsterAt(p)
	if i <= 0 || !g.ECS.Alive(i) {
		return logerror.New("You have to target a monster.")
	}
	g.Logf("%s looks confused (scroll).", color.ColorLogPlayerAttack, g.ECS.GetName(i))
	g.ECS.PutStatus(i, StatusConfused, sc.Turns)
	return nil
}

func (sc *ConfusionScroll) TargetingRadius() int { return 0 }

// FireballScroll is an item that can be invoked to produce a flame explosion
// in an area around a target position.
type FireballScroll struct {
	Damage int
	Radius int
}

func (sc *FireballScroll) Activate(g *game, a itemAction) error {
	if a.Target == nil {
		return logerror.New("You have to chose a target.")
	}
	p := *a.Target
	if !g.InFOV(p) {
		return logerror.New("You cannot target what you cannot see.")
	}
	hits := 0
	// NOTE: this could be made more complicated by checking whether there
	// are monsters in the way. For now, it's a fireball that goes up and
	// then down and explodes on reaching the target!
	for i, st := range g.ECS.Stats {
		if g.ECS.Dead(i) {
			continue
		}
		q := g.ECS.Positions[i]
		dist := paths.DistanceManhattan(q, p)
		if dist > sc.Radius {
			continue
		}
		g.Logf("%v is engulfed in flames.", color.ColorLogPlayerAttack, g.ECS.GetName(i))
		st.HP -= sc.Damage
		hits++
	}
	if hits <= 0 {
		return logerror.New("There are no targets in the radius.")
	}
	return nil
}

func (sc *FireballScroll) TargetingRadius() int { return sc.Radius }
