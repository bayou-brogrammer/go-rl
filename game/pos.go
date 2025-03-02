package game

import "codeberg.org/anaseto/gruid"

func (a actionType) actionToXY() (int,int) {
	switch a {
	case ActionN:
		return 0, -1 // north - up
	case ActionS:
		return 0, 1  // south - down
	case ActionW:
		return -1, 0 // west - left
	case ActionE:
		return 1, 0  // east - right
	default:
		return 0, 0 // no movement for other actions
	}
}

func keyToDir(k actionType) (p gruid.Point) {
	switch k {
	case ActionW:
		p = gruid.Point{-1, 0}
	case ActionE:
		p = gruid.Point{1, 0}
	case ActionS:
		p = gruid.Point{0, 1}
	case ActionN:
		p = gruid.Point{0, -1}
	}
	return p
}
