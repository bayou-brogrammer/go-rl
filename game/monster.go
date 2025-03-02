package game

import (
	"math/rand"

	"codeberg.org/anaseto/gruid/paths"
)

// Monster represents a monster.
type Monster struct{}

// HandleMonsterTurn handles a monster's turn. The function assumes the entity
// with the given index is indeed a monster initialized with fighter and AI
// components.
func (g *game) HandleMonsterTurn(i int) {
	if !g.ECS.Alive(i) {
		// Do nothing if the entity corresponds to a dead monster.
		return
	}
	if g.ECS.Status(i, StatusConfused) {
		g.HandleConfusedMonster(i)
		return
	}
	p := g.ECS.Positions[i]
	ai := g.ECS.AI[i]
	aip := &aiPath{g: g}
	pp := g.ECS.PP()
	if paths.DistanceManhattan(p, pp) == 1 {
		// If the monster is adjacent to the player, attack.
		g.BumpAttack(i, g.ECS.PlayerID)
		return
	}
	if !g.InFOV(p) {
		// The monster is not in player's FOV.
		if len(ai.Path) < 1 {
			// Pick new path to a random floor tile.
			ai.Path = g.PR.AstarPath(aip, p, g.Map.RandomFloor())
		}
		g.AIMove(i)
		// NOTE: this base AI can be improved for example to avoid
		// monster's getting stuck between them. It's enough to get
		// started, though.
		return
	}
	// The monster is in player's FOV, so we compute a suitable path to
	// reach the player.
	ai.Path = g.PR.AstarPath(aip, p, pp)
	g.AIMove(i)
}

// HandleConfusedMonster handles the behavior of a confused monster. It simply
// tries to bump into a random direction.
func (g *game) HandleConfusedMonster(i int) {
	p := g.ECS.Positions[i]
	p.X += -1 + 2*rand.Intn(2)
	p.Y += -1 + 2*rand.Intn(2)
	if !p.In(g.Map.Grid.Range()) {
		return
	}
	if p == g.ECS.PP() {
		g.BumpAttack(i, g.ECS.PlayerID)
		return
	}
	if g.Map.Walkable(p) && g.ECS.NoBlockingEntityAt(p) {
		g.ECS.MoveEntity(i, p)
	}
}

// AIMove moves a monster to the next position, if there is no blocking entity
// at the destination. It assumes the destination is walkable.
func (g *game) AIMove(i int) {
	ai := g.ECS.AI[i]
	if len(ai.Path) > 0 && ai.Path[0] == g.ECS.Positions[i] {
		ai.Path = ai.Path[1:]
	}
	if len(ai.Path) > 0 && g.ECS.NoBlockingEntityAt(ai.Path[0]) {
		// Only move if there is no blocking entity.
		g.ECS.MoveEntity(i, ai.Path[0])
		ai.Path = ai.Path[1:]
	}
}
