package main

import (
	"math"
	"math/rand"
)

// Terrain generates and renders procedural hilly ground.
// The ground region is FILLED with textured characters, not just an outline.
type Terrain struct {
	width, height int
	theme         string
	rng           *rand.Rand

	// heights[x] = the y row of the SURFACE of the terrain at column x
	// (lower y value = higher on screen = taller hill)
	heights    []int
	baseY      int // the lowest the terrain can be (where water starts)
	offsetY    int // top of the terrain region (hills peak here)
	wavePhase  float64
}

func NewTerrain(w, h, baseY int, rng *rand.Rand, theme string) *Terrain {
	t := &Terrain{
		width:  w,
		height: h,
		theme:  theme,
		rng:    rng,
		baseY:  baseY,
		// Hills can peak at ~30% from baseY towards the top of screen
		offsetY: baseY - (h / 3),
	}
	t.generate()
	return t
}

func (t *Terrain) generate() {
	t.heights = make([]int, t.width)

	// Combine 3 sine waves at different frequencies for organic hills
	freq1 := 0.03 + t.rng.Float64()*0.02
	freq2 := 0.07 + t.rng.Float64()*0.03
	freq3 := 0.12 + t.rng.Float64()*0.05

	phase1 := t.rng.Float64() * math.Pi * 2
	phase2 := t.rng.Float64() * math.Pi * 2
	phase3 := t.rng.Float64() * math.Pi * 2

	amp := float64(t.baseY-t.offsetY) * 0.55
	mid := float64(t.baseY) + amp*0.2 // shift midpoint down so valleys dip well below water level (baseY)

	for x := 0; x < t.width; x++ {
		fx := float64(x)
		h := mid +
			math.Sin(fx*freq1+phase1)*amp*0.6 +
			math.Sin(fx*freq2+phase2)*amp*0.25 +
			math.Sin(fx*freq3+phase3)*amp*0.15

		t.heights[x] = clamp(int(h), t.offsetY+1, t.height-3)
	}
}

func (t *Terrain) Resize(w, h, baseY int) {
	t.width = w
	t.height = h
	t.baseY = baseY
	t.offsetY = baseY - (h / 3)
	t.generate()
}

func (t *Terrain) Update() {
	t.wavePhase += 0.05
}

// GetSurfaceY returns the surface row at column x (clamped).
func (t *Terrain) GetSurfaceY(x int) int {
	if x < 0 {
		return t.heights[0]
	}
	if x >= t.width {
		return t.heights[t.width-1]
	}
	return t.heights[x]
}

// MinSurfaceY returns the highest point the terrain reaches (smallest y).
func (t *Terrain) MinSurfaceY() int {
	m := t.height
	for _, h := range t.heights {
		if h < m {
			m = h
		}
	}
	return m
}

// IsUnderground returns true if the cell at (x,y) is inside the ground body.
func (t *Terrain) IsUnderground(x, y int) bool {
	return y >= t.GetSurfaceY(x)
}

// surfaceChar picks the character to draw at a ground cell based on position.
var terrainSurface = []rune{'_', '_', '~', '_'}
var terrainFill = []rune{'.', ',', '\'', '`', '.', ',', 'w', 'v', 'W', 'V', '~', ';', 'w', 'v', '.', ','}
var terrainSlope = []rune{'/', '\\', '|', '/'}

func (t *Terrain) Draw(s Screen) {

	for x := 0; x < t.width; x++ {
		surfY := t.heights[x]
		for y := surfY; y < t.height; y++ {
			// Slowly cycle character indexing and color brightness to prevent burn-in
			earthWave := math.Sin(float64(x)*0.08 + float64(y)*0.15 - t.wavePhase*0.2)
			brightnessMult := 1.0 + earthWave*0.12

			var cSurfFG, cFillFG, cBgCol Color
			if t.theme == "night" {
				cSurfFG = NewRGBColor(int32(30*brightnessMult), int32(100*brightnessMult), int32(30*brightnessMult))
				cFillFG = NewRGBColor(int32(20*brightnessMult), int32(60*brightnessMult), int32(20*brightnessMult))
				cBgCol = NewRGBColor(int32(10*brightnessMult), int32(35*brightnessMult), int32(10*brightnessMult))
			} else {
				cSurfFG = NewRGBColor(int32(60*brightnessMult), int32(160*brightnessMult), int32(40*brightnessMult))
				cFillFG = NewRGBColor(int32(45*brightnessMult), int32(110*brightnessMult), int32(30*brightnessMult))
				cBgCol = NewRGBColor(int32(30*brightnessMult), int32(80*brightnessMult), int32(20*brightnessMult))
			}

			var ch rune
			var fg Color
			if y == surfY {
				// surface line: pick between _ / \ based on slope
				left := surfY
				if x > 0 {
					left = t.heights[x-1]
				}
				right := surfY
				if x < t.width-1 {
					right = t.heights[x+1]
				}
				switch {
				case left > surfY && right > surfY:
					ch = '^'
					fg = cSurfFG
				case left > surfY:
					ch = '\\'
					fg = cFillFG
				case right > surfY:
					ch = '/'
					fg = cFillFG
				default:
					ch = '_'
					fg = cSurfFG
				}

				// Apply swaying grass to flat surfaces
				if ch == '_' {
					sway := math.Sin(float64(x)*0.25 + t.wavePhase)
					if sway > 0.7 {
						ch = 'w' // swaying grass clumps
						fg = cSurfFG
					} else if sway > 0.3 {
						ch = 'v' // single grass blade
						fg = cSurfFG
					} else if sway < -0.7 {
						ch = '"' // double blade
						fg = cSurfFG
					} else if sway < -0.3 {
						ch = ',' // tiny sprout
						fg = cFillFG
					}
				}
			} else if y == surfY+1 {
				// Animated roots / pebbles layer just below surface
				sway := math.Sin(float64(x)*0.4 - t.wavePhase*0.7)
				if sway > 0.6 {
					ch = '.'
				} else if sway < -0.6 {
					ch = ','
				} else {
					ch = '\''
				}
				fg = cFillFG
			} else {
				// Bedrock / deep soil fill below surface: slowly crawl the textures
				phaseShift := int(t.wavePhase * 0.3)
				idx := (x*3 + y*7 + phaseShift) % len(terrainFill)
				ch = terrainFill[idx]
				fg = cFillFG
			}
			style := StyleDefault.Foreground(fg).Background(cBgCol)
			s.SetContent(x, y, ch, nil, style)
		}
	}
}
