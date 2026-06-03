package main

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// Water renders the animated water/river section at the bottom of the screen.
type Water struct {
	width, height int
	startY        int // first row of water
	theme         string
	rng           *rand.Rand
	phase         float64
}

func NewWater(w, h, startY int, rng *rand.Rand, theme string) *Water {
	return &Water{
		width:  w,
		height: h,
		startY: startY,
		theme:  theme,
		rng:    rng,
	}
}

func (wa *Water) Resize(w, h, startY int) {
	wa.width = w
	wa.height = h
	wa.startY = startY
}

func (wa *Water) Update() {
	wa.phase += 0.15
}

// Draw renders water only in columns where the terrain surface dips below startY (valleys).
func (wa *Water) Draw(s tcell.Screen, t *Terrain) {
	for x := 0; x < wa.width; x++ {
		surfY := t.GetSurfaceY(x)
		if surfY <= wa.startY {
			continue // Hill is above water level, no water here
		}

		// Draw water from wa.startY to surfY - 1
		for y := wa.startY; y < surfY; y++ {
			depth := y - wa.startY
			totalWaterDepth := surfY - wa.startY
			if totalWaterDepth <= 0 {
				continue
			}

			// Pick colors based on depth relative to totalWaterDepth
			var (
				fgR, fgG, fgB int32
				bgR, bgG, bgB int32
			)
			tVal := float64(depth) / float64(totalWaterDepth)
			if wa.theme == "night" {
				fgR = int32(lerp(0, 0, tVal))
				fgG = int32(lerp(220, 60, tVal))
				fgB = int32(lerp(240, 120, tVal))
				bgR = 0
				bgG = int32(lerp(40, 15, tVal))
				bgB = int32(lerp(140, 50, tVal))
			} else {
				fgR = int32(lerp(100, 30, tVal))
				fgG = int32(lerp(200, 80, tVal))
				fgB = int32(lerp(255, 180, tVal))
				bgR = 0
				bgG = int32(lerp(40, 10, tVal))
				bgB = int32(lerp(120, 40, tVal))
			}
			fg := tcell.NewRGBColor(fgR, fgG, fgB)
			bg := tcell.NewRGBColor(bgR, bgG, bgB)

			wave := math.Sin(float64(x)*0.18+wa.phase) +
				math.Sin(float64(x)*0.07-wa.phase*0.7) +
				math.Sin(float64(y)*0.5+wa.phase*0.4)

			var ch rune
			switch {
			case depth == 0:
				// Surface row: visible wave crests
				if wave > 0.8 {
					ch = '≈'
				} else if wave > 0.1 {
					ch = '~'
				} else if wave > -0.5 {
					ch = '-'
				} else {
					ch = '.'
				}
			case depth == 1:
				if wave > 0.5 {
					ch = '~'
				} else if wave > -0.2 {
					ch = '-'
				} else {
					ch = '.'
				}
			default:
				// Deeper rows: mostly dots and dashes, occasional tilde
				if wave > 1.0 {
					ch = '~'
				} else if wave > 0.3 {
					ch = '.'
				} else if wave > -0.3 {
					ch = ' '
				} else {
					ch = '.'
				}
			}

			style := tcell.StyleDefault.Foreground(fg).Background(bg)
			s.SetContent(x, y, ch, nil, style)
		}
	}
}

func lerp(a, b int, t float64) float64 {
	return float64(a) + (float64(b)-float64(a))*t
}
