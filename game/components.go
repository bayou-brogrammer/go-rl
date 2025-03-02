// This file describes entity components, for example for basic fighting or AI.

package game

import "codeberg.org/anaseto/gruid"

type Stats struct {
	HP      int // Health Points
	MaxHP   int // Maximum Health Points
	Power   int // attack power
	Defense int // defence
}

// Heal heals a fighter for a certain amount, if it does not exceed maximum HP.
// The final amount of healing is returned.
func (st *Stats) Heal(n int) int {
	st.HP += n
	if st.HP > st.MaxHP {
		n -= st.HP - st.MaxHP
		st.HP = st.MaxHP
	}
	return n
}

// AI holds simple AI data for monster's.
type AI struct {
	Path []gruid.Point // path to destination
}

// Style contains information relative to the default graphical representation
// of an entity.
type Style struct {
	Rune  rune
	Color gruid.Color
}

// Inventory holds items. For now, consumables.
type Inventory struct {
	Items []int
}

// status describes different kind of statuses.
type status int

const (
	StatusConfused status = iota
)

// Statuses maps ongoing statuses to their remaining turns.
type Statuses map[status]int

// NextTurn updates statuses for the next turn.
func (sts Statuses) NextTurn() {
	for st, turns := range sts {
		if turns == 0 {
			delete(sts, st)
			continue
		}
		sts[st]--
	}
}

// Put puts on a particular status for a given number of turns.
func (sts Statuses) Put(st status, turns int) {
	sts[st] = turns
}
