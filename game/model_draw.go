package game

import (
	"sort"
	"strings"
	"unicode/utf8"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
	"github.com/bayou-brogrammer/go-rl/game/color"
)

// Draw implements gruid.Model.Draw. It draws a simple map that spans the whole
// grid.
func (m *model) Draw() gruid.Grid {
	mapgrid := m.grid.Slice(m.grid.Range().Shift(0, LogLines, 0, -1))

	switch m.mode {
	case modeGameMenu:
		return m.DrawGameMenu()
	case modeMessageViewer:
		m.grid.Copy(m.viewer.Draw())
		return m.grid
	case modeInventoryDrop, modeInventoryActivate:
		mapgrid.Copy(m.inventory.Draw())
		return m.grid
	}

	m.grid.Fill(gruid.Cell{Rune: ' '})
	g := m.game

	// We draw the map tiles.
	it := g.Map.Grid.Iterator()
	for it.Next() {
		if !g.Map.Explored[it.P()] {
			continue
		}
		c := gruid.Cell{Rune: g.Map.Rune(it.Cell())}
		if g.InFOV(it.P()) {
			c.Style.Bg = color.ColorFOV
		}
		mapgrid.Set(it.P(), c)
	}

	// We sort entity indexes using the render ordering.
	sortedEntities := make([]int, 0, len(g.ECS.Entities))
	for i := range g.ECS.Entities {
		sortedEntities = append(sortedEntities, i)
	}
	sort.Slice(sortedEntities, func(i, j int) bool {
		return g.ECS.RenderOrder(sortedEntities[i]) < g.ECS.RenderOrder(sortedEntities[j])
	})

	// We draw the sorted entities.
	for _, i := range sortedEntities {
		p := g.ECS.Positions[i]
		if !g.Map.Explored[p] || !g.InFOV(p) {
			continue
		}
		c := mapgrid.At(p)
		c.Rune, c.Style.Fg = g.ECS.GetStyle(i)
		mapgrid.Set(p, c)
		// NOTE: We retrieved current cell at e.Pos() to preserve
		// background (in FOV or not).
	}

	m.DrawNames(mapgrid)
	m.DrawLog(m.grid.Slice(m.grid.Range().Lines(0, LogLines)))
	m.DrawStatus(m.grid.Slice(m.grid.Range().Line(m.grid.Size().Y - 1)))

	return m.grid
}

var mainMenuAnchor = gruid.Point{10, 6}

// DrawGameMenu draws the game's main menu.
func (m *model) DrawGameMenu() gruid.Grid {
	m.grid.Fill(gruid.Cell{Rune: ' '})
	m.grid.Slice(m.gameMenu.Bounds().Add(mainMenuAnchor)).Copy(m.gameMenu.Draw())
	m.info.Draw(m.grid.Slice(m.grid.Range().Line(12).Shift(10, 0, 0, 0)))
	return m.grid
}

// DrawLog draws the last two lines of the log.
func (m *model) DrawLog(gd gruid.Grid) {
	j := 1
	for i := len(m.game.Log) - 1; i >= 0; i-- {
		if j < 0 {
			break
		}
		e := m.game.Log[i]
		st := gruid.Style{}
		st.Fg = e.Color
		m.log.Content = ui.NewStyledText(e.String(), st)
		m.log.Draw(gd.Slice(gd.Range().Line(j)))
		j--
	}
}

// DrawStatus draws the status line
func (m *model) DrawStatus(gd gruid.Grid) {
	st := gruid.Style{}
	st.Fg = color.ColorStatusHealthy
	g := m.game
	stats := g.ECS.Stats[g.ECS.PlayerID]
	if stats.HP < stats.MaxHP/2 {
		st.Fg = color.ColorStatusWounded
	}
	m.log.Content = ui.Textf("HP: %d/%d", stats.HP, stats.MaxHP).WithStyle(st)
	m.log.Draw(gd)
}

// DrawNames renders the names of the named entities at current mouse location
// if it is in the map.
func (m *model) DrawNames(gd gruid.Grid) {
	maprg := gruid.NewRange(0, LogLines, UIWidth, UIHeight-1)
	if !m.targ.pos.In(maprg) {
		return
	}

	p := m.targ.pos.Sub(maprg.Min)
	rad := m.targ.radius
	rg := gruid.Range{Min: p.Sub(gruid.Point{rad, rad}), Max: p.Add(gruid.Point{rad + 1, rad + 1})}
	rg = rg.Intersect(maprg.Sub(maprg.Min))
	rg.Iter(func(q gruid.Point) {
		c := gd.At(q)
		c.Style.Attrs |= color.AttrReverse
		gd.Set(q, c)
	})

	// We get the names of the entities at p.
	names := []string{}
	for i, q := range m.game.ECS.Positions {
		if q != p || !m.game.InFOV(q) {
			continue
		}
		name := m.game.ECS.GetName(i)
		if name != "" {
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		return
	}

	// We sort the names. This could be improved to sort by entity type
	// too, as well as to remove duplicates (for example showing “corpse
	// (3x)” if there are three corpses).
	sort.Strings(names)

	text := strings.Join(names, ", ")
	width := utf8.RuneCountInString(text) + 2
	rg = gruid.NewRange(p.X+1, p.Y-1, p.X+1+width, p.Y+2)
	// we adjust a bit the box's placement in case it's on a edge.
	if p.X+1+width >= UIWidth {
		rg = rg.Shift(-1-width, 0, -1-width, 0)
	}
	if p.Y+2 > MapHeight {
		rg = rg.Shift(0, -1, 0, -1)
	}
	if p.Y-1 < 0 {
		rg = rg.Shift(0, 1, 0, 1)
	}
	slice := gd.Slice(rg)
	m.desc.Content = ui.Text(text)
	m.desc.Draw(slice)
}
