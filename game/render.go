package game

// renderOrder is a type representing the priority of an entity rendering.
type renderOrder int

// Those constants represent distinct kinds of rendering priorities. In case
// two entities are at a given position, only the one with the highest priority
// gets displayed.
const (
	RONone renderOrder = iota
	ROCorpse
	ROItem
	ROActor
)

// RenderOrder returns the rendering priority of an entity.
func (es *ECS) RenderOrder(i int) (ro renderOrder) {
	switch es.Entities[i].(type) {
	case *Player:
		ro = ROActor
	case *Monster:
		if es.Dead(i) {
			ro = ROCorpse
		} else {
			ro = ROActor
		}
	case *Consumable:
		ro = ROItem
	}
	return ro
}
