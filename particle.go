package main

import (
	"math"
	"math/rand"
)

// ParticleType enumerates the visual effect types.
type ParticleType int

const (
	PExplosion ParticleType = iota
	PSpark
	PLaser
	PRain
)

// Particle is a short-lived visual effect.
type Particle struct {
	ptype   ParticleType
	x, y    float64
	vx, vy  float64
	life    int
	maxLife int
	color   Color
	ch      rune

	// For lasers: start + end point
	x2, y2 float64
}

func (p *Particle) IsDead() bool { return p.life <= 0 }

func (p *Particle) Update() {
	p.life--
	if p.ptype == PSpark {
		p.x += p.vx
		p.y += p.vy
		p.vy += 0.12 // gravity
		p.vx *= 0.95 // drag
	}
	if p.ptype == PRain {
		p.y += p.vy
		p.x += p.vx
	}
}

func (p *Particle) Draw(s Screen) {
	if p.life <= 0 {
		return
	}
	sw, sh := s.Size()

	switch p.ptype {
	case PExplosion:
		p.drawExplosion(s, sw, sh)
	case PSpark:
		p.drawPoint(s, sw, sh)
	case PLaser:
		p.drawLaser(s, sw, sh)
	case PRain:
		p.drawPoint(s, sw, sh)
	}
}

func (p *Particle) drawPoint(s Screen, sw, sh int) {
	x, y := int(p.x), int(p.y)
	if x < 0 || x >= sw || y < 0 || y >= sh {
		return
	}
	_, _, curStyle, _ := s.GetContent(x, y)
	_, bg, _ := curStyle.Decompose()
	style := StyleDefault.Foreground(p.color).Background(bg)
	s.SetContent(x, y, p.ch, nil, style)
}

// drawExplosion renders a radial burst that expands then fades.
func (p *Particle) drawExplosion(s Screen, sw, sh int) {
	age := float64(p.maxLife-p.life) / float64(p.maxLife)
	radius := age * 5.0

	// Pick character and colour based on age
	var ch rune
	var r, g, b int32
	switch {
	case age < 0.2:
		ch = '@'
		r, g, b = 255, 255, 200
	case age < 0.4:
		ch = '#'
		r, g, b = 255, 160, 0
	case age < 0.6:
		ch = '*'
		r, g, b = 220, 80, 0
	case age < 0.8:
		ch = '+'
		r, g, b = 150, 40, 0
	default:
		ch = '.'
		r, g, b = 80, 20, 0
	}
	col := NewRGBColor(r, g, b)

	// Draw a rough circle of characters
	steps := 20
	for i := 0; i < steps; i++ {
		angle := float64(i) * math.Pi * 2 / float64(steps)
		cx := int(p.x + math.Cos(angle)*radius)
		cy := int(p.y + math.Sin(angle)*radius*0.5) // squash vertically for terminal aspect
		if cx >= 0 && cx < sw && cy >= 0 && cy < sh {
			_, _, curStyle, _ := s.GetContent(cx, cy)
			_, bg, _ := curStyle.Decompose()
			style := StyleDefault.Foreground(col).Background(bg)
			s.SetContent(cx, cy, ch, nil, style)
		}
	}
	// Centre glyph
	cx, cy := int(p.x), int(p.y)
	if cx >= 0 && cx < sw && cy >= 0 && cy < sh {
		_, _, curStyle, _ := s.GetContent(cx, cy)
		_, bg, _ := curStyle.Decompose()
		style := StyleDefault.Foreground(col).Background(bg)
		s.SetContent(cx, cy, ch, nil, style)
	}
}

// drawLaser draws a Bresenham line between (x,y) and (x2,y2).
func (p *Particle) drawLaser(s Screen, sw, sh int) {
	x0, y0 := int(p.x), int(p.y)
	x1, y1 := int(p.x2), int(p.y2)

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	// Colour fades based on life
	frac := float64(p.life) / float64(p.maxLife)
	r := int32(255)
	g := int32(30 + frac*100)
	b := int32(0)
	laserCol := NewRGBColor(r, g, b)

	x, y := x0, y0
	for {
		if x >= 0 && x < sw && y >= 0 && y < sh {
			ch := '-'
			if dx < dy {
				ch = '|'
			}
			if math.Abs(float64(x-x0)) < 2 || math.Abs(float64(x-x1)) < 2 {
				ch = '='
			}
			_, _, curStyle, _ := s.GetContent(x, y)
			_, bg, _ := curStyle.Decompose()
			style := StyleDefault.Foreground(laserCol).Background(bg)
			s.SetContent(x, y, rune(ch), nil, style)
		}
		if x == x1 && y == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ── Factory functions ────────────────────────────────────────────────────────

func NewExplosionParticle(x, y float64, rng *rand.Rand) *Particle {
	return &Particle{
		ptype:   PExplosion,
		x:       x + float64(rng.Intn(3)-1),
		y:       y + float64(rng.Intn(2)-1),
		life:    12 + rng.Intn(8),
		maxLife: 20,
	}
}

func NewSparkParticle(x, y float64, isRobot bool, rng *rand.Rand) *Particle {
	angle := rng.Float64() * math.Pi * 2
	speed := 0.5 + rng.Float64()*1.5
	var ch rune
	switch rng.Intn(4) {
	case 0:
		ch = '*'
	case 1:
		ch = '+'
	case 2:
		ch = 'x'
	default:
		ch = '·'
	}

	var col Color
	if isRobot {
		// Red sparks (Robot lightsaber)
		r := int32(220 + rng.Intn(35))
		g := int32(rng.Intn(60))
		b := int32(rng.Intn(30))
		col = NewRGBColor(r, g, b)
	} else {
		// Blue sparks (Knight lightsaber)
		r := int32(rng.Intn(30))
		g := int32(120 + rng.Intn(80))
		b := int32(220 + rng.Intn(35))
		col = NewRGBColor(r, g, b)
	}

	return &Particle{
		ptype:   PSpark,
		x:       x,
		y:       y,
		vx:      math.Cos(angle) * speed,
		vy:      math.Sin(angle) * speed * 0.5,
		life:    8 + rng.Intn(10),
		maxLife: 18,
		color:   col,
		ch:      ch,
	}
}

func NewLaserParticle(x1, y1, x2, y2 float64) *Particle {
	return &Particle{
		ptype:   PLaser,
		x:       x1,
		y:       y1,
		x2:      x2,
		y2:      y2,
		life:    5,
		maxLife: 5,
	}
}

func NewRainParticle(x float64, rng *rand.Rand) *Particle {
	length := 2 + rng.Intn(4)
	return &Particle{
		ptype:   PRain,
		x:       x,
		y:       float64(-length),
		vx:      0.1,
		vy:      1.2 + rng.Float64()*0.8,
		life:    20 + rng.Intn(20),
		maxLife: 40,
		color:   NewRGBColor(60, 100, 200),
		ch:      '|',
	}
}
