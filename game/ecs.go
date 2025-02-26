package main

import (
	"codeberg.org/anaseto/gruid"
	"github.com/bayou-brogrammer/go-rl/game/color"
)

// Entity represents an object or creature on the map.
type Entity interface{}

// ECS manages entities, as well as their positions. We don't go full “ECS”
// (Entity-Component-System) in this tutorial, opting for a simpler hybrid
// approach good enough for the tutorial purposes.
type ECS struct {
	Entities  map[int]Entity      // set of entities
	Positions map[int]gruid.Point // entity index: map position
	PlayerID  int                 // index of Player's entity (for convenience)
	NextID    int                 // next available id

	Stats     map[int]*Stats     // stats component
	AI        map[int]*AI        // AI component
	Name      map[int]string     // name component
	Style     map[int]Style      // default style component
	Inventory map[int]*Inventory // inventory component
	Statuses  map[int]Statuses   // statuses (confused, etc.)
}

// NewECS returns an initialized ECS structure.
func NewECS() *ECS {
	return &ECS{
		Entities:  map[int]Entity{},
		Positions: map[int]gruid.Point{},
		Stats:     map[int]*Stats{},
		AI:        map[int]*AI{},
		Name:      map[int]string{},
		Style:     map[int]Style{},
		Inventory: map[int]*Inventory{},
		Statuses:  map[int]Statuses{},
		NextID:    0,
	}
}

// Add adds a new entity at a given position and returns its index/id.
func (es *ECS) AddEntity(e Entity, p gruid.Point) int {
	id := es.NextID
	es.Entities[id] = e
	es.Positions[id] = p
	es.NextID++
	return id
}

// AddItem is a shorthand for adding item entities on the map.
func (es *ECS) AddItem(e Entity, p gruid.Point, name string, r rune) int {
	id := es.AddEntity(e, p)
	es.Name[id] = name
	es.Style[id] = Style{Rune: r, Color: color.ColorConsumable}
	return id
}

// RemoveEntity removes an entity, given its identifier.
func (es *ECS) RemoveEntity(i int) {
	delete(es.Entities, i)
	delete(es.Positions, i)
	delete(es.Stats, i)
	delete(es.AI, i)
	delete(es.Name, i)
	delete(es.Style, i)
	delete(es.Inventory, i)
	delete(es.Statuses, i)
}

// MoveEntity moves the i-th entity to p.
func (es *ECS) MoveEntity(i int, p gruid.Point) {
	es.Positions[i] = p
}

// MovePlayer moves the player entity to p.
func (es *ECS) MovePlayer(p gruid.Point) {
	es.MoveEntity(es.PlayerID, p)
}

// Player returns the Player entity. Just a shorthand for easily accessing the
// Player entity.
func (es *ECS) Player() *Player {
	return es.Entities[es.PlayerID].(*Player)
}

// PP returns the Player's position. Just a convenience shorthand.
func (es *ECS) PP() gruid.Point {
	return es.Positions[es.PlayerID]
}

// MonsterAt returns the id of the Monster at p, if any, or -1 if there is no
// monster at p.
func (es *ECS) MonsterAt(p gruid.Point) int {
	for i, q := range es.Positions {
		if p != q || !es.Alive(i) {
			continue
		}
		e := es.Entities[i]
		switch e.(type) {
		case *Monster:
			return i
		}
	}
	return -1
}

// NoBlockingEntityAt returns true if there is no blocking entity at p (no
// player nor monsters in this tutorial).
func (es *ECS) NoBlockingEntityAt(p gruid.Point) bool {
	i := es.MonsterAt(p)
	return es.PP() != p && !es.Alive(i)
}

// PlayerDied checks whether the player died.
func (es *ECS) PlayerDied() bool {
	return es.Dead(es.PlayerID)
}

// Alive checks whether an entity is alive.
func (es *ECS) Alive(i int) bool {
	st := es.Stats[i]
	return st != nil && st.HP > 0
}

// Dead checks whether an entity is dead (was alive).
func (es *ECS) Dead(i int) bool {
	st := es.Stats[i]
	return st != nil && st.HP <= 0
}

// GetStyle returns the graphical representation (rune and foreground color) of an
// entity.
func (es *ECS) GetStyle(i int) (r rune, c gruid.Color) {
	r = es.Style[i].Rune
	c = es.Style[i].Color
	if es.Dead(i) {
		// Alternate representation for corpses of dead monsters.
		r = '%'
		c = gruid.ColorDefault
	}
	return r, c
}

// GetName returns the name of an entity, which most often is name given by the
// Name component, except for corpses.
func (es *ECS) GetName(i int) (s string) {
	name := es.Name[i]
	if es.Dead(i) {
		name = name + " corpse"
	}
	return name
}

// StatusesNextTurn updates the remaining turns of entities' statuses.
func (es *ECS) StatusesNextTurn() {
	for _, sts := range es.Statuses {
		sts.NextTurn()
	}
}

// PutStatus puts on a particular status for a given entity for a certain
// number of turns.
func (es *ECS) PutStatus(i int, st status, turns int) {
	if es.Statuses[i] == nil {
		es.Statuses[i] = map[status]int{}
	}
	sts := es.Statuses[i]
	sts.Put(st, turns)
}

// Status checks whether an entity has a particular status effect.
func (es *ECS) Status(i int, st status) bool {
	_, ok := es.Statuses[i][st]
	return ok
}
