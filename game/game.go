// This file handles game related affairs that are not specific to entities or
// the map.

package game

import (
	"fmt"
	"math/rand"
	"time"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
	"github.com/bayou-brogrammer/go-rl/game/color"
	"github.com/bayou-brogrammer/go-rl/game/constants"
	"github.com/bayou-brogrammer/go-rl/game/dungeon"
	"github.com/bayou-brogrammer/go-rl/game/logerror"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// game represents information relevant the current game's state.
type game struct {
	Version string
	Depth   int

	ECS *ECS             // entities present on the map
	Map *dungeon.Map     // the game map, made of tiles
	PR  *paths.PathRange // path range for the map
	Log []LogEntry       // log entries

	rand *rand.Rand // random number generator
}

func (g *game) initalizeFirstLevel() {
	g.Version = constants.Version
	g.Depth++ // start at 1

	size := gruid.Point{UIWidth, UIHeight}
	size.Y -= 3 // for log and status

	g.Map = dungeon.NewMap(size, dungeon.Town)
	g.PR = paths.NewPathRange(gruid.NewRange(0, 0, size.X, size.Y))

	// Initialize entities
	g.ECS = NewECS()

	// Initialization: create a player entity centered on the map.
	g.ECS.PlayerID = g.ECS.AddEntity(NewPlayer(), g.Map.RandomFloor())
	g.ECS.Stats[g.ECS.PlayerID] = &Stats{
		HP: 30, MaxHP: 30, Power: 5, Defense: 2,
	}
	g.ECS.Style[g.ECS.PlayerID] = Style{Rune: '@', Color: color.ColorPlayer}
	g.ECS.Name[g.ECS.PlayerID] = "player"
	g.ECS.Inventory[g.ECS.PlayerID] = &Inventory{}
	g.UpdateFOV()

	// Add some monsters
	// g.SpawnMonsters()

	// Add items
	g.PlaceItems()
}

func (g *game) InitalizeLevel() {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

}

// SpawnMonsters adds some monsters in the current map.
func (g *game) SpawnMonsters() {
	const numberOfMonsters = 12
	for i := 0; i < numberOfMonsters; i++ {
		m := &Monster{}
		// We generate either an orc or a troll with 0.8 and 0.2
		// probabilities respectively.
		const (
			orc = iota
			troll
		)
		kind := orc
		switch {
		case g.Map.GetRand().Intn(100) < 80:
		default:
			kind = troll
		}
		p := g.FreeFloorTile()
		i := g.ECS.AddEntity(m, p)
		switch kind {
		case orc:
			g.ECS.Stats[i] = &Stats{
				HP: 10, MaxHP: 10, Defense: 0, Power: 3,
			}
			g.ECS.Name[i] = "orc"
			g.ECS.Style[i] = Style{Rune: 'o', Color: color.ColorMonster}
		case troll:
			g.ECS.Stats[i] = &Stats{
				HP: 16, MaxHP: 16, Defense: 1, Power: 4,
			}
			g.ECS.Name[i] = "troll"
			g.ECS.Style[i] = Style{Rune: 'T', Color: color.ColorMonster}
		}
		g.ECS.AI[i] = &AI{}
	}
}

// FreeFloorTile returns a free floor tile in the map (it assumes it exists).
func (g *game) FreeFloorTile() gruid.Point {
	for {
		p := g.Map.RandomFloor()
		if g.ECS.NoBlockingEntityAt(p) {
			return p
		}
	}
}

// EndTurn is called when the player's turn ends. Currently, the player and
// monsters have all the same speed, so we make each monster act each time the
// player's does an action that ends a turn.
func (g *game) EndTurn() gruid.Effect {
	g.UpdateFOV()
	for i, e := range g.ECS.Entities {
		if g.ECS.PlayerDied() {
			return nil
		}
		switch e.(type) {
		case *Monster:
			g.HandleMonsterTurn(i)
		}
	}
	g.ECS.StatusesNextTurn()
	return nil
}

// UpdateFOV updates the field of view.
func (g *game) UpdateFOV() {
	player := g.ECS.Player()
	// player position
	pp := g.ECS.PP()

	// We shift the FOV's Range so that it will be centered on the new
	// player's position.
	rg := gruid.NewRange(-maxLOS, -maxLOS, maxLOS+1, maxLOS+1)
	player.FOV.SetRange(rg.Add(pp).Intersect(g.Map.Grid.Range()))

	// We mark cells in field of view as explored. We use the symmetric
	// shadow casting algorithm provided by the rl package.
	passable := func(p gruid.Point) bool {
		return g.Map.Cell(p) != dungeon.WallCell
	}

	for _, p := range player.FOV.SSCVisionMap(pp, maxLOS, passable, false) {
		if paths.DistanceManhattan(p, pp) > maxLOS {
			continue
		}
		if !g.Map.Explored[p] {
			g.Map.Explored[p] = true
		}
	}
}

// InFOV returns true if p is in the player's field of view. We only keep cells
// within maxLOS manhattan distance from the player, as natural given our
// current 4-way movement. With 8-way movement, the natural distance choice
// would be the Chebyshev one.
func (g *game) InFOV(p gruid.Point) bool {
	pp := g.ECS.PP()
	return g.ECS.Player().FOV.Visible(p) &&
		paths.DistanceManhattan(pp, p) <= maxLOS
}

// BumpAttack implements attack of a fighter entity on another.
func (g *game) BumpAttack(i, j int) {
	st := g.ECS.Stats[i]
	stj := g.ECS.Stats[j]
	damage := st.Power - stj.Defense
	attackDesc := fmt.Sprintf("%s attacks %s", cases.Title(language.English).String(g.ECS.Name[i]), g.ECS.Name[j])
	col := color.ColorLogMonsterAttack

	if i == g.ECS.PlayerID {
		col = color.ColorLogPlayerAttack
	}

	if damage > 0 {
		g.Logf("%s for %d damage", col, attackDesc, damage)
		stj.HP -= damage
	} else {
		g.Logf("%s but does no damage", col, attackDesc)
	}
}

// PlaceItems adds items in the current map.
func (g *game) PlaceItems() {
	const numberOfItems = 5
	for i := 0; i < numberOfItems; i++ {
		p := g.FreeFloorTile()
		r := g.Map.GetRand().Float64()
		switch {
		case r < 0.7:
			g.ECS.AddItem(&HealingPotion{Amount: 4}, p, "health potion", '!')
		case r < 0.8:
			g.ECS.AddItem(&ConfusionScroll{Turns: 10}, p, "confusion scroll", '?')
		case r < 0.9:
			g.ECS.AddItem(&FireballScroll{Damage: 12, Radius: 3}, p, "fireball scroll", '?')
		default:
			g.ECS.AddItem(&LightningScroll{Range: 5, Damage: 20},
				p, "lightning scroll", '?')
		}
	}
}

const ErrNoShow = "ErrNoShow"

// IventoryAdd adds an item to the player's inventory, if there is room. It
// returns an error if the item could not be added.
func (g *game) InventoryAdd(actor, i int) error {
	const maxSize = 26
	switch g.ECS.Entities[i].(type) {
	case Consumable:
		inv := g.ECS.Inventory[actor]
		if len(inv.Items) >= maxSize {
			return logerror.New("Inventory is full.")
		}
		inv.Items = append(inv.Items, i)
		delete(g.ECS.Positions, i)
		return nil
	}
	return logerror.New(ErrNoShow)
}

// Drop an item from the inventory.
func (g *game) InventoryRemove(actor, n int) error {
	inv := g.ECS.Inventory[actor]
	if len(inv.Items) <= n {
		return logerror.New("Empty slot.")
	}
	i := inv.Items[n]
	inv.Items[n] = inv.Items[len(inv.Items)-1]
	inv.Items = inv.Items[:len(inv.Items)-1]
	g.ECS.Positions[i] = g.ECS.PP()
	return nil
}

// InventoryActivate uses a given item from the inventory.
func (g *game) InventoryActivate(actor, n int) error {
	return g.InventoryActivateWithTarget(actor, n, nil)
}

// NeedsTargeting checks whether using the n-th item requires targeting,
// returning its radius (-1 if no targeting).
func (g *game) TargetingRadius(n int) int {
	inv := g.ECS.Inventory[g.ECS.PlayerID]
	if len(inv.Items) <= n {
		return -1
	}
	i := inv.Items[n]
	switch e := g.ECS.Entities[i].(type) {
	case Targetter:
		return e.TargetingRadius()
	}
	return -1
}
