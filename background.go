package main

import (
	"math"
	"math/rand"
)

// Star represents a single twinkling star in the sky.
type Star struct {
	X, Y  int
	Char  rune
	Color Color
	Phase int
	Speed int // ticks per state change
}

// ShootingStar represents a fast-moving diagonal star with a trail.
type ShootingStar struct {
	X, Y    float64
	Vx, Vy  float64
	Length  int
	Life    int
	MaxLife int
}

// BioFlash is a rare biological/whimsical animation that scrolls across the sky.
type BioFlash struct {
	X      float64 // left edge of the art
	Y      int     // top row in the sky
	speed  float64 // columns per tick (negative = move left)
	kind   int     // 0=DNA helix, 1=amino acid chain, 2=banner plane
	ticker int     // internal animation tick
	label  string  // for banner planes: slogan chosen at spawn time
}

// ── Bio art generators ────────────────────────────────────────────────────────

// makeDNAHelix generates a proper animated ASCII double helix.
// The two strands are positioned using cosine waves that are π apart,
// so they cross in the middle, widen at the top/bottom.
func makeDNAHelix(tick int) ([]string, Color, Color) {
	const height = 14
	const width = 15

	// Base pair labels cycling A-T, T-A, G-C, C-G
	basePairs := [][2]byte{
		{'A', 'T'}, {'T', 'A'}, {'G', 'C'}, {'C', 'G'},
		{'C', 'G'}, {'A', 'T'}, {'T', 'A'}, {'G', 'C'},
	}

	// Phase advances slowly each tick for smooth rotation
	phase := float64(tick) * 0.18

	rows := make([]string, height)
	center := float64(width-1) / 2.0
	amplitude := center * 0.88

	for row := 0; row < height; row++ {
		line := make([]byte, width)
		for i := range line {
			line[i] = ' '
		}

		angle := float64(row)*0.65 + phase

		// Strand 1 and strand 2 are π apart
		p1 := int(math.Round(center + math.Cos(angle)*amplitude))
		p2 := int(math.Round(center + math.Cos(angle+math.Pi)*amplitude))

		// Clamp
		if p1 < 0 {
			p1 = 0
		}
		if p1 >= width {
			p1 = width - 1
		}
		if p2 < 0 {
			p2 = 0
		}
		if p2 >= width {
			p2 = width - 1
		}

		bp := basePairs[(row+tick/8)%len(basePairs)]
		lo, hi := p1, p2
		b0, b1 := bp[0], bp[1]
		if lo > hi {
			lo, hi = hi, lo
			b0, b1 = b1, b0
		}

		if lo == hi {
			// Strands cross — draw an X
			line[lo] = 'X'
		} else {
			line[lo] = b0
			line[hi] = b1
			dist := hi - lo - 1
			for i := lo + 1; i < hi; i++ {
				if dist <= 1 {
					line[i] = '-'
				} else if dist <= 4 {
					line[i] = '-'
				} else {
					line[i] = '='
				}
			}
		}
		rows[row] = string(line)
	}
	return rows, NewRGBColor(0, 220, 200), NewRGBColor(0, 120, 255)
}

// makeAminoChain generates a flowing amino acid / protein chain.
func makeAminoChain(tick int) ([]string, Color, Color) {
	// Amino acid one-letter codes cycling through
	aminos := "ACDEFGHIKLMNPQRSTVWY"
	phase := tick / 5

	rows := []string{
		"     ╭─────╮     ",
		"  ╭──┤ " + string(aminos[(phase+0)%len(aminos)]) + string(aminos[(phase+1)%len(aminos)]) + " ├──╮  ",
		"  │  ╰──┬──╯  │  ",
		"╭─┴─╮   │  ╭──┴─╮",
		"│ " + string(aminos[(phase+2)%len(aminos)]) + string(aminos[(phase+3)%len(aminos)]) + " │   │  │ " + string(aminos[(phase+4)%len(aminos)]) + string(aminos[(phase+5)%len(aminos)]) + " │",
		"╰─┬─╯   │  ╰──┬─╯",
		"  │  ╭──┴──╮  │  ",
		"  ╰──┤ " + string(aminos[(phase+6)%len(aminos)]) + string(aminos[(phase+7)%len(aminos)]) + " ├──╯  ",
		"     ╰─────╯     ",
	}
	return rows, NewRGBColor(255, 160, 50), NewRGBColor(180, 80, 0)
}

// Slogan tier sizes — explicit so the weighted sampler below is correct.
const (
	sloganCommonN = 10 // 60% chance
	sloganNerdyN  = 14 // 30% chance  (biophysics + obscure facts)
	sloganRareN   = 6  // 10% chance  (philosophical)
)

// allSlogans is the full pool ordered: common → nerdy → rare.
var allSlogans = []string{
	// ── Common (biology / code crossover) ────────────────────────────────
	"~~ LIFE IS CODE ~~~~~~~~~~~",    // 0
	"~~ DNA OR DIE! ~~~~~~~~~~~~",    // 1
	"~~ EVOLVE OR PERISH ~~~~~~~~",   // 2
	"~~ ACGT FOR EVER ~~~~~~~~~~~",   // 3
	"~~ HELLO, CARBON! ~~~~~~~~~~",   // 4
	"~~ 3.8 BILLION YEARS ~~~~~~~~",  // 5
	"~~ MITOSIS > DIVISION ~~~~~~~",  // 6
	"~~ ATP: NATURE'S BATTERY ~~~~",  // 7
	"~~ YOU ARE 60% WATER ~~~~~~~~",  // 8
	"~~ CELLS ALL THE WAY DOWN ~~~",  // 9

	// ── Nerdy: biophysics + obscure molecular bio ─────────────────────────
	"~~ kT IS EVERYTHING ~~~~~~~~~",  // 10
	"~~ BOLTZMANN NEVER LIES ~~~~~~",  // 11
	"~~ BROWNIAN MOTION IS LIFE ~~",  // 12
	"~~ LIPID BILAYERS RULE ~~~~~~",   // 13
	"~~ PROTEIN FOLDING IS NP-HARD",  // 14
	"~~ VAN DER WAALS SAYS HI ~~~~",  // 15
	"~~ DIFFUSION RUNS THE SHOW ~~",  // 16
	"~~ RNA WORLD HYPOTHESIS ~~~~~",  // 17
	"~~ CRISPR WAS HERE ~~~~~~~~~~",  // 18
	"~~ 98.7% CHIMP DNA ~~~~~~~~~~",  // 19
	"~~ ENTROPY ALWAYS WINS ~~~~~~",  // 20
	"~~ TELOMERES = AGING ~~~~~~~~",  // 21
	"~~ MAXWELL'S DEMON DISAGREES ",  // 22
	"~~ ATP SYNTHASE SPINS @ 9000 ",  // 23

	// ── Rare (philosophical) ──────────────────────────────────────────────
	"~~ WE ARE STAR STUFF ~~~~~~~~~",  // 24
	"~~ LIFE FINDS A WAY ~~~~~~~~~~",  // 25
	"~~ E = mc^2 ... AND DNA ~~~~~~",  // 26
	"~~ THE SELFISH GENE SAYS HI ~~",  // 27
	"~~ DARWIN WAS RIGHT ~~~~~~~~~~",  // 28
	"~~ ON THE ORIGIN OF CODE ~~~~~",  // 29
}

// makeBannerPlane generates a detailed biplane towing a slogan banner.
// label is chosen once at spawn time so it stays consistent for the full pass.
func makeBannerPlane(label string, tick int) ([]string, Color, Color) {
	// Spinning propeller: 4-frame cycle  | / - \
	propFrames := []string{" | ", " / ", " - ", " \\ "}
	prop := propFrames[(tick/3)%len(propFrames)]

	// Tow rope length pulses slightly for a flapping look
	ropeLen := 8 + (tick/6)%3
	rope := ""
	for i := 0; i < ropeLen; i++ {
		rope += "-"
	}

	// 6-row biplane: top-wing, cockpit row, fuselage, bottom-wing, wheels
	rows := []string{
		`        ___________         `,
		`       /___________\        `,
		 prop + `===( >  [o]   )` + rope + label,
		`       \_____==____/        `,
		`       /___________\        `,
		`           |   |            `,
		`          _|_ _|_           `,
	}
	return rows, NewRGBColor(255, 215, 50), NewRGBColor(200, 120, 0)
}

// bioArtLines dispatches to the right generator.
func bioArtLines(kind, tick int, label string) ([]string, Color, Color) {
	switch kind {
	case 0:
		return makeDNAHelix(tick)
	case 1:
		return makeAminoChain(tick)
	default:
		return makeBannerPlane(label, tick)
	}
}

// ── Background ────────────────────────────────────────────────────────────────

// Background fills the sky area with a calm, twinkling night starfield
// and occasional shooting stars.
type Background struct {
	width, height int
	theme         string
	rng           *rand.Rand
	stars         []*Star
	shootingStars []*ShootingStar
	driftTimer    int

	// Bio flash system
	bioTimer int // counts down to next spawn
	bioFlash *BioFlash
}

var starChars = []rune{'.', '·', '*', '+'}

func NewBackground(w, h int, rng *rand.Rand, theme string) *Background {
	bg := &Background{
		width:    w,
		height:   h,
		theme:    theme,
		rng:      rng,
		bioTimer: 300 + rng.Intn(100), // first flash after ~10-13s at 30fps
	}
	bg.buildStars()
	return bg
}

func (bg *Background) buildStars() {
	bg.stars = nil
	bg.shootingStars = nil

	// Density: around 3% of the sky area
	numStars := (bg.width * bg.height) * 3 / 100
	if numStars < 10 {
		numStars = 10
	}

	occupied := make(map[int]bool)

	for i := 0; i < numStars; i++ {
		x := bg.rng.Intn(bg.width)
		y := bg.rng.Intn(bg.height)
		key := y*bg.width + x

		if occupied[key] {
			continue
		}
		occupied[key] = true

		char := starChars[bg.rng.Intn(len(starChars))]
		var col Color
		switch char {
		case '.':
			col = NewRGBColor(80, 80, 100)
		case '·':
			col = NewRGBColor(130, 130, 160)
		case '*':
			col = NewRGBColor(210, 210, 240)
		case '+':
			col = NewRGBColor(190, 210, 255)
		}

		bg.stars = append(bg.stars, &Star{
			X:     x,
			Y:     y,
			Char:  char,
			Color: col,
			Phase: bg.rng.Intn(100),
			Speed: 10 + bg.rng.Intn(20),
		})
	}
}

func (bg *Background) Resize(w, h int) {
	bg.width = w
	bg.height = h
	bg.buildStars()
}

func (bg *Background) Update() {
	// Drift stars slowly to prevent screen burn-in
	bg.driftTimer++
	if bg.driftTimer >= 60 {
		bg.driftTimer = 0
		for _, s := range bg.stars {
			s.X--
			if s.X < 0 {
				s.X = bg.width - 1
				s.Y = bg.rng.Intn(bg.height)
			}
		}
	}

	// 1. Twinkle stars
	for _, s := range bg.stars {
		s.Phase++
		if s.Phase%s.Speed == 0 {
			switch bg.rng.Intn(4) {
			case 0:
				s.Char = '.'
				s.Color = NewRGBColor(70, 70, 90)
			case 1:
				s.Char = '·'
				s.Color = NewRGBColor(120, 120, 150)
			case 2:
				s.Char = '*'
				s.Color = NewRGBColor(210, 210, 230)
			case 3:
				s.Char = '+'
				s.Color = NewRGBColor(180, 200, 240)
			}
		}
	}

	// 2. Update shooting stars
	for i := len(bg.shootingStars) - 1; i >= 0; i-- {
		ss := bg.shootingStars[i]
		ss.X += ss.Vx
		ss.Y += ss.Vy
		ss.Life++
		if ss.Life >= ss.MaxLife || int(ss.X) < 0 || int(ss.X) >= bg.width || int(ss.Y) < 0 || int(ss.Y) >= bg.height {
			bg.shootingStars = append(bg.shootingStars[:i], bg.shootingStars[i+1:]...)
		}
	}

	// 3. Spawn shooting stars
	if len(bg.shootingStars) < 2 && bg.rng.Float64() < 0.006 {
		startX := float64(bg.width/3 + bg.rng.Intn(bg.width*2/3))
		startY := float64(bg.rng.Intn(bg.height / 4))
		vx := -1.0 - bg.rng.Float64()*1.2
		vy := 0.3 + bg.rng.Float64()*0.4
		length := 3 + bg.rng.Intn(4)
		maxLife := 15 + bg.rng.Intn(15)
		bg.shootingStars = append(bg.shootingStars, &ShootingStar{
			X:       startX,
			Y:       startY,
			Vx:      vx,
			Vy:      vy,
			Length:  length,
			MaxLife: maxLife,
		})
	}

	// 4. Bio flash: advance or spawn
	if bg.bioFlash != nil {
		bf := bg.bioFlash
		bf.ticker++
		bf.X += bf.speed

		// Retire when fully off the left edge
		lines, _, _ := bioArtLines(bf.kind, bf.ticker, bf.label)
		maxW := 0
		for _, l := range lines {
			if len(l) > maxW {
				maxW = len(l)
			}
		}
		if int(bf.X)+maxW < 0 {
			bg.bioFlash = nil
			bg.bioTimer = 270 + bg.rng.Intn(200) // 9–16 s until next
		}
	} else {
		bg.bioTimer--
		if bg.bioTimer <= 0 {
			kind := bg.rng.Intn(3)

			// Three-tier weighted slogan selection:
			//   0-5  (60%) → common
			//   6-8  (30%) → nerdy / biophysics
			//   9    (10%) → rare / philosophical
			var chosenLabel string
			switch roll := bg.rng.Intn(10); {
			case roll < 6:
				chosenLabel = allSlogans[bg.rng.Intn(sloganCommonN)]
			case roll < 9:
				chosenLabel = allSlogans[sloganCommonN+bg.rng.Intn(sloganNerdyN)]
			default:
				chosenLabel = allSlogans[sloganCommonN+sloganNerdyN+bg.rng.Intn(sloganRareN)]
			}

			// Y position: spread across top 40% of sky, at least 1 row down
			maxY := maxInt(2, bg.height*2/5)
			spawnY := 1 + bg.rng.Intn(maxY)

			bg.bioFlash = &BioFlash{
				X:     float64(bg.width + 5),
				Y:     spawnY,
				speed: -(0.32 + bg.rng.Float64()*0.28),
				kind:  kind,
				label: chosenLabel,
			}
		}
	}
}

// Draw renders the starry night sky background.
func (bg *Background) Draw(s Screen, t *Terrain) {
	var skyBG Color
	if bg.theme == "night" {
		skyBG = NewRGBColor(0, 0, 8)
	} else {
		skyBG = NewRGBColor(5, 5, 20)
	}

	// 1. Clear sky down to terrain surface
	for x := 0; x < bg.width; x++ {
		surfY := t.GetSurfaceY(x)
		for y := 0; y < surfY; y++ {
			s.SetContent(x, y, ' ', nil, StyleDefault.Background(skyBG))
		}
	}

	// 2. Stars
	for _, star := range bg.stars {
		if star.X >= 0 && star.X < bg.width {
			surfY := t.GetSurfaceY(star.X)
			if star.Y < surfY {
				style := StyleDefault.Foreground(star.Color).Background(skyBG)
				s.SetContent(star.X, star.Y, star.Char, nil, style)
			}
		}
	}

	// 3. Shooting stars
	for _, ss := range bg.shootingStars {
		for j := 0; j < ss.Length; j++ {
			tx := int(ss.X - ss.Vx*float64(j)*0.7)
			ty := int(ss.Y - ss.Vy*float64(j)*0.7)
			if tx >= 0 && tx < bg.width {
				surfY := t.GetSurfaceY(tx)
				if ty >= 0 && ty < surfY {
					var ch rune
					var col Color
					if j == 0 {
						ch = '*'
						col = NewRGBColor(255, 255, 255)
					} else if j < 2 {
						ch = '+'
						col = NewRGBColor(200, 220, 255)
					} else {
						ch = '·'
						col = NewRGBColor(100, 120, 160)
					}
					style := StyleDefault.Foreground(col).Background(skyBG)
					s.SetContent(tx, ty, ch, nil, style)
				}
			}
		}
	}

	// 4. Bio flash
	if bg.bioFlash != nil {
		bf := bg.bioFlash
		lines, fgCol, dimCol := bioArtLines(bf.kind, bf.ticker, bf.label)
		sw, sh := s.Size()
		startX := int(bf.X)

		for row, line := range lines {
			y := bf.Y + row
			if y < 0 || y >= sh {
				continue
			}
			for col, ch := range line {
				x := startX + col
				if x < 0 || x >= sw {
					continue
				}
				if ch == ' ' {
					continue
				}
				surfY := t.GetSurfaceY(x)
				if y >= surfY {
					continue
				}
				// Alternate between bright fg and dim accent for a glow pulse
				fg := fgCol
				if (col+row+bf.ticker/3)%4 == 0 {
					fg = dimCol
				}
				style := StyleDefault.Foreground(fg).Background(skyBG)
				s.SetContent(x, y, ch, nil, style)
			}
		}
	}
}
