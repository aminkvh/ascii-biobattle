package main

import (
	"math"
	"math/rand"
)

// UnitState describes what a unit is currently doing.
type UnitState int

const (
	StateWander UnitState = iota
	StateApproach
	StateAttack
	StateDefend
	StateHit
	StateDead
)

// Unit represents a battle entity (Knight or Robot).
type Unit struct {
	isRobot bool

	// Position of top-left corner of the sprite in terminal coords
	px, py float64

	// Movement
	vx          float64
	speed       float64
	lungeOffset float64 // temporary offset for attack lunge animations

	// Combat
	health      int
	maxHealth   int
	attackCD    int
	attackRange float64 // distance in columns

	// Animation
	state     UnitState
	frame     int
	frameTick int
	frames    []SpriteFrame
	hitTimer  int

	rng *rand.Rand
}

func NewKnightUnit(rng *rand.Rand, w, h int, t *Terrain) *Unit {
	sprites := KnightSprites
	spW := float64(sprites[0].Width())
	x := rng.Float64() * float64(w) * 0.35 // Spread Knights over 35% of left screen
	if x+spW > float64(w) {
		x = float64(w) - spW - 1
	}
	surfY := float64(t.GetSurfaceY(int(x)))
	y := surfY - float64(sprites[0].Height()) + 1

	return &Unit{
		isRobot:     false,
		px:          x,
		py:          y,
		speed:       0.06 + rng.Float64()*0.04, // slower, smoother movement
		health:      120,
		maxHealth:   120,
		attackRange: 13, // melee lightsaber range (wider than sprite width)
		frames:      sprites,
		rng:         rng,
	}
}

func NewRobotUnit(rng *rand.Rand, w, h int, t *Terrain) *Unit {
	sprites := RobotSprites
	spW := float64(sprites[0].Width())
	x := float64(w)*0.65 + rng.Float64()*float64(w)*0.35 // Spread Robots over 35% of right screen
	if x+spW > float64(w) {
		x = float64(w) - spW - 1
	}
	surfY := float64(t.GetSurfaceY(int(x)))
	y := surfY - float64(sprites[0].Height()) + 1

	return &Unit{
		isRobot:     true,
		px:          x,
		py:          y,
		speed:       0.05 + rng.Float64()*0.04, // slower, smoother movement
		health:      120,
		maxHealth:   120,
		attackRange: 13, // Robots now fight in melee with red lightsabers!
		frames:      sprites,
		rng:         rng,
	}
}

// cx / cy return the centre of this unit in terminal coordinates.
func (u *Unit) cx() float64 {
	return u.px + float64(u.frames[0].Width())/2
}
func (u *Unit) cy() float64 {
	return u.py + float64(u.frames[0].Height())/2
}

func (u *Unit) IsDead() bool { return u.state == StateDead && u.hitTimer <= 0 }

func (u *Unit) Update(w *World) {
	if u.state == StateDead {
		u.hitTimer--
		return
	}

	// Hit flash timeout
	if u.hitTimer > 0 {
		u.hitTimer--
		if u.hitTimer == 0 && u.state == StateHit {
			u.state = StateApproach
		}
	}

	// Cooldown
	if u.attackCD > 0 {
		u.attackCD--
	}

	// Decay lunge offset back to 0
	u.lungeOffset *= 0.8
	if math.Abs(u.lungeOffset) < 0.05 {
		u.lungeOffset = 0
	}

	target := w.FindNearestEnemy(u)

	// ── State machine ──────────────────────────────────────────────
	if target == nil {
		u.state = StateWander
	} else {
		dx := target.cx() - u.cx()
		dist := math.Abs(dx)

		if dist <= u.attackRange {
			// In range
			if u.attackCD <= 0 {
				u.state = StateAttack
				var dmg int
				if u.isRobot {
					dmg = 15 + u.rng.Intn(10) // Buffed Robot melee damage (15-24) to counter Knight's shield
				} else {
					dmg = 9 + u.rng.Intn(8) // Knight melee damage (9-16)
				}
				target.TakeDamage(dmg, w)
				u.attackCD = 25 + u.rng.Intn(15)

				// Trigger visual attack lunge forward
				lungeMag := 2.2
				if dx > 0 {
					u.lungeOffset = lungeMag
				} else {
					u.lungeOffset = -lungeMag
				}

				// Spawn sparks corresponding to the attacker
				for k := 0; k < 4; k++ {
					w.AddParticle(NewSparkParticle(target.cx(), target.cy(), u.isRobot, u.rng))
				}
			} else {
				if !u.isRobot {
					u.state = StateDefend
				} else {
					u.state = StateApproach
				}
			}
		} else {
			u.state = StateApproach
		}
	}

	// ── Movement ───────────────────────────────────────────────────
	switch u.state {
	case StateWander:
		u.px += (u.rng.Float64()*2 - 1) * 0.1
	case StateApproach, StateDefend, StateAttack:
		if target != nil {
			dx := target.cx() - u.cx()
			dist := math.Abs(dx)

			// Both units are melee now, walk directly into contact range
			if dist > 2.0 {
				mult := 1.0
				if dist > 15.0 {
					mult = 1.5 // charge speed boost when approaching from afar
				} else if u.state == StateDefend || u.state == StateAttack {
					mult = 0.5 // move slower during sword swing or block
				}
				if dx > 0 {
					u.px += u.speed * mult
				} else {
					u.px -= u.speed * mult
				}
			}
		}
	}

	// ── Friendly Separation Force ──────────────────────────────────
	// Prevent friendly units from stacking on top of each other
	for _, other := range w.units {
		if other == u || other.isRobot != u.isRobot || other.state == StateDead {
			continue
		}
		distX := u.px - other.px
		absDistX := math.Abs(distX)
		if absDistX < 14.0 { // width of friendly separation zone (wider than sprite to avoid overlap)
			push := 0.25 // push force per tick
			if distX > 0 {
				u.px += push
			} else if distX < 0 {
				u.px -= push
			} else {
				// Exact overlap: break tie randomly
				if u.isRobot {
					u.px -= push
				} else {
					u.px += push
				}
			}
		}
	}

	// Clamp to screen edges
	maxX := float64(w.width) - float64(u.frames[0].Width()) - 1
	if u.px < 0 {
		u.px = 0
	}
	if u.px > maxX {
		u.px = maxX
	}

	// Snap to terrain surface (with smooth vertical interpolation to glide)
	surfY := float64(w.terrain.GetSurfaceY(int(u.cx())))
	targetY := surfY - float64(u.frames[0].Height()) + 1
	
	u.py += (targetY - u.py) * 0.4
	if u.py < 0 {
		u.py = 0
	}
	maxValY := float64(w.height) - float64(u.frames[0].Height())
	if u.py > maxValY {
		u.py = maxValY
	}

	// Advance animation frame
	u.frameTick++
	framesPerAnim := 8
	if u.frameTick >= framesPerAnim {
		u.frameTick = 0
		u.frame = (u.frame + 1) % 2 // walk A/B loop (frames 0 & 1)
	}
}

func (u *Unit) TakeDamage(dmg int, w *World) {
	// Block mechanic: 30% damage reduction when in Defend state (take 70% damage)
	if u.state == StateDefend {
		dmg = int(float64(dmg) * 0.7)
		if dmg < 1 {
			dmg = 1
		}
	}

	u.health -= dmg

	// Visual knockback: push unit slightly away from attacker
	knockback := 1.2 + u.rng.Float64()*1.2
	if u.isRobot {
		u.px += knockback // push right
	} else {
		u.px -= knockback // push left
	}

	if u.health <= 0 {
		u.health = 0
		u.state = StateDead
		u.hitTimer = 15 // show death sprite for 15 ticks
	} else {
		u.state = StateHit
		u.hitTimer = 6
	}
}

// spriteFrameIndex maps state to the correct sprite frame index.
func (u *Unit) spriteFrameIndex() int {
	switch u.state {
	case StateWander, StateApproach:
		return u.frame // 0 or 1 (walk cycle)
	case StateAttack:
		return 2
	case StateDefend:
		return 3
	case StateHit:
		return 4
	case StateDead:
		return 5
	}
	return 0
}

func (u *Unit) Draw(s Screen, t *Terrain) {
	fi := u.spriteFrameIndex()
	if fi >= len(u.frames) {
		fi = 0
	}
	frame := u.frames[fi]

	// Add temporary lungeOffset to draw position
	px, py := int(u.px+u.lungeOffset), int(u.py)
	sw, sh := s.Size()

	for row, line := range frame.Lines {
		y := py + row
		if y < 0 || y >= sh {
			continue
		}
		for col, ch := range line {
			x := px + col
			if x < 0 {
				continue
			}
			if x >= sw {
				break
			}
			if ch == ' ' {
				continue
			}
			// Don't draw over terrain ground body
			if t.IsUnderground(x, y) {
				continue
			}

			_, _, curStyle, _ := s.GetContent(x, y)
			_, bg, _ := curStyle.Decompose()
			style := u.charStyle(ch, row, frame.Height()).Background(bg)
			s.SetContent(x, y, ch, nil, style)
		}
	}

	// Health bar above sprite (Tiny, minimal bar)
	u.drawHealthBar(s, px, py-1, frame.Width())
}

func (u *Unit) drawHealthBar(s Screen, x, y, width int) {
	if y < 0 {
		return
	}
	sw, _ := s.Size()
	
	// Minimal health bar: always 3 characters wide, centered above unit
	barWidth := 3
	xStart := x + (width-barWidth)/2

	pct := float64(u.health) / float64(u.maxHealth)
	filled := int(math.Round(pct * float64(barWidth)))
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}

	var hpFG Color
	switch {
	case pct > 0.6:
		hpFG = NewRGBColor(50, 220, 50) // green
	case pct > 0.3:
		hpFG = NewRGBColor(220, 180, 0) // yellow/orange
	default:
		hpFG = NewRGBColor(220, 30, 30) // red
	}

	for i := 0; i < barWidth; i++ {
		cx := xStart + i
		if cx < 0 || cx >= sw {
			continue
		}
		var ch rune
		var col Color
		if i < filled {
			ch = '■' // small solid block
			col = hpFG
		} else {
			ch = '·' // small dot
			col = NewRGBColor(60, 60, 60)
		}
		style := StyleDefault.Foreground(col).Background(ColorBlack)
		s.SetContent(cx, y, ch, nil, style)
	}
}

// charStyle assigns a neon colour to each character based on which unit type
// and where the character sits in the sprite (armour, face, legs, etc.)
func (u *Unit) charStyle(ch rune, row, totalRows int) Style {
	bg := ColorBlack

	if !u.isRobot {
		// ── Knight: silver armour with blue shield highlights ──────
		var fg Color
		switch {
		case row <= 2:
			// Visor and helmet
			if ch == '/' || ch == '\\' || ch == 'V' || ch == '_' {
				fg = NewRGBColor(210, 210, 225) // silver helm
			} else {
				fg = NewRGBColor(240, 240, 255)
			}
		case row <= 3:
			fg = NewRGBColor(170, 175, 190) // pauldrons / neck
		case row <= 6:
			// Chest armor: steel chest plate with blue shield highlights
			if ch == '#' || ch == '[' || ch == ']' {
				fg = NewRGBColor(65, 135, 245) // bright blue shield highlights
			} else {
				fg = NewRGBColor(180, 185, 200) // steel armor
			}
		case row <= 8:
			fg = NewRGBColor(160, 165, 185) // silver thighs / knees
		default:
			fg = NewRGBColor(130, 135, 150) // silver feet
		}

		// Hit flash (even ticks during hit)
		if u.state == StateHit && u.hitTimer%2 == 0 {
			fg = NewRGBColor(255, 100, 100)
		}

		// Blue lightsaber glow (neon blue)
		if ch == '=' || ch == '>' {
			fg = NewRGBColor(0, 160, 255) // bright blue saber blade
		}
		if ch == '/' && row >= 4 && row <= 7 {
			fg = NewRGBColor(0, 160, 255) // diagonal blue saber blade in walk/defend
		}

		return StyleDefault.Foreground(fg).Background(bg)
	}

	// ── Robot: red/orange core with metallic limbs ────────────────
	var fg Color
	switch {
	case row <= 2:
		// Head and visor
		if ch == 'o' {
			fg = NewRGBColor(255, 120, 0) // orange glowing sensors
		} else if ch == 'X' {
			fg = NewRGBColor(255, 60, 60) // red cooling core
		} else {
			fg = NewRGBColor(180, 45, 45) // red plating
		}
	case row <= 3:
		fg = NewRGBColor(120, 120, 130) // shoulder mounts (steel)
	case row <= 6:
		// Chest and core glow
		if ch == '#' {
			fg = NewRGBColor(255, 70, 0) // bright orange/red energy core
		} else if ch == '[' || ch == ']' {
			fg = NewRGBColor(110, 110, 120) // arm steel plates
		} else {
			fg = NewRGBColor(150, 35, 35) // crimson plating
		}
	case row <= 9:
		fg = NewRGBColor(120, 30, 30) // red/crimson legs
	default:
		fg = NewRGBColor(90, 90, 100) // metallic feet
	}

	// Hit flash
	if u.state == StateHit && u.hitTimer%2 == 0 {
		fg = NewRGBColor(255, 255, 50) // bright yellow flash
	}

	// Red lightsaber glow (neon red)
	if ch == '=' || ch == '<' {
		fg = NewRGBColor(255, 30, 30) // glowing red saber blade
	}
	if ch == '/' && row >= 4 && row <= 7 {
		fg = NewRGBColor(255, 30, 30) // diagonal red saber blade in walk/defend
	}

	return StyleDefault.Foreground(fg).Background(bg)
}
