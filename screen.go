package main

import (
	"image/color"
	"sync"
)

type Color struct {
	r, g, b int32
}

func NewRGBColor(r, g, b int32) Color {
	return Color{r: r, g: g, b: b}
}

func (c Color) RGBA() (r, g, b, a uint32) {
	// Multiply by 257 to scale 0-255 to 0-65535
	return uint32(c.r) * 257, uint32(c.g) * 257, uint32(c.b) * 257, 0xffff
}

func (c Color) ToRGBA() color.RGBA {
	return color.RGBA{R: uint8(c.r), G: uint8(c.g), B: uint8(c.b), A: 255}
}

var ColorBlack = NewRGBColor(0, 0, 0)
var ColorDefault = ColorBlack

type Style struct {
	fg Color
	bg Color
}

func (s Style) Foreground(c Color) Style { s.fg = c; return s }
func (s Style) Background(c Color) Style { s.bg = c; return s }
func (s Style) Decompose() (fg Color, bg Color, attrs int32) { return s.fg, s.bg, 0 }

var StyleDefault = Style{fg: NewRGBColor(255, 255, 255), bg: ColorBlack}

type Cell struct {
	Primary rune
	Style   Style
}

type Screen interface {
	SetContent(x, y int, primary rune, comb []rune, style Style)
	GetContent(x, y int) (primary rune, comb []rune, style Style, width int)
	Clear()
	Size() (width, height int)
	Show()
	Sync()
}

type GridScreen struct {
	width  int
	height int
	cells  []Cell
	mu     sync.RWMutex
}

func NewGridScreen(width, height int) *GridScreen {
	return &GridScreen{
		width:  width,
		height: height,
		cells:  make([]Cell, width*height),
	}
}

func (g *GridScreen) Resize(width, height int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.width = width
	g.height = height
	g.cells = make([]Cell, width*height)
}

func (g *GridScreen) Size() (int, int) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.width, g.height
}

func (g *GridScreen) SetContent(x, y int, primary rune, comb []rune, style Style) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if x < 0 || y < 0 || x >= g.width || y >= g.height {
		return
	}
	g.cells[y*g.width+x] = Cell{
		Primary: primary,
		Style:   style,
	}
}

func (g *GridScreen) Clear() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i := range g.cells {
		g.cells[i] = Cell{Primary: ' ', Style: StyleDefault}
	}
}

func (g *GridScreen) Show() {}
func (g *GridScreen) Sync() {}

func (g *GridScreen) GetCell(x, y int) Cell {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if x < 0 || y < 0 || x >= g.width || y >= g.height {
		return Cell{Primary: ' ', Style: StyleDefault}
	}
	return g.cells[y*g.width+x]
}

func (g *GridScreen) GetContent(x, y int) (primary rune, comb []rune, style Style, width int) {
	c := g.GetCell(x, y)
	return c.Primary, nil, c.Style, 1
}
