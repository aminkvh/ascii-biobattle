package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// World is the root state container. It holds all subsystems and orchestrates
// the update and draw order each tick.
type World struct {
	width, height int
	density       int
	intensity     int
	theme         string

	rng *rand.Rand

	// Subsystems
	bg      *Background
	terrain *Terrain
	water   *Water
	hud     *HUD

	// Entities
	units     []*Unit
	particles []*Particle

	// Rain particles (atmospheric)
	rainTimer int

	// Stats
	tick        int
	knightKills int
	robotKills  int
}

// Layout constants (as fractions of screen height)
const (
	waterFrac  = 0.15 // bottom 15% = water
	groundFrac = 0.35 // next 20% = hill zone (terrain can peak here)
)

func NewWorld(w, h, density, intensity int, theme string) *World {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	world := &World{
		width:     w,
		height:    h,
		density:   density,
		intensity: intensity,
		theme:     theme,
		rng:       rng,
	}
	world.init()
	return world
}

func (w *World) init() {
	waterStartY := w.height - maxInt(3, int(float64(w.height)*waterFrac))
	baseY := waterStartY // terrain surface can be at most here

	w.bg = NewBackground(w.width, w.height, w.rng, w.theme)
	w.terrain = NewTerrain(w.width, w.height, baseY, w.rng, w.theme)
	w.water = NewWater(w.width, w.height, waterStartY, w.rng, w.theme)
	w.hud = NewHUD(w.theme)

	w.units = nil
	for i := 0; i < w.density; i++ {
		w.units = append(w.units, NewKnightUnit(w.rng, w.width, w.height, w.terrain))
		w.units = append(w.units, NewRobotUnit(w.rng, w.width, w.height, w.terrain))
	}
}

func (w *World) Resize(width, height int) {
	w.width = width
	w.height = height
	// Re-init subsystems but keep kill counts
	kk, rk := w.knightKills, w.robotKills
	tk := w.tick
	w.init()
	w.knightKills = kk
	w.robotKills = rk
	w.tick = tk
}

func (w *World) Update() {
	w.tick++

	// Subsystem updates
	w.bg.Update()
	w.terrain.Update()
	w.water.Update()

	// Spawn atmospheric rain drops occasionally
	w.rainTimer++
	if w.rainTimer >= 3 {
		w.rainTimer = 0
		if w.rng.Intn(4) == 0 {
			x := float64(w.rng.Intn(w.width))
			w.particles = append(w.particles, NewRainParticle(x, w.rng))
		}
	}

	// Update units
	for i := len(w.units) - 1; i >= 0; i-- {
		u := w.units[i]
		u.Update(w)

		if u.IsDead() {
			// Kill tracking
			if u.isRobot {
				w.knightKills++
			} else {
				w.robotKills++
			}
			// Big death explosion
			for k := 0; k < 12; k++ {
				w.particles = append(w.particles, NewExplosionParticle(u.cx(), u.cy(), w.rng))
			}
			// Remove and respawn to keep battle going
			w.units = append(w.units[:i], w.units[i+1:]...)
			if u.isRobot {
				w.units = append(w.units, NewRobotUnit(w.rng, w.width, w.height, w.terrain))
			} else {
				w.units = append(w.units, NewKnightUnit(w.rng, w.width, w.height, w.terrain))
			}
		}
	}

	// Update particles
	for i := len(w.particles) - 1; i >= 0; i-- {
		w.particles[i].Update()
		if w.particles[i].IsDead() {
			w.particles = append(w.particles[:i], w.particles[i+1:]...)
		}
	}
}

func (w *World) Draw(s tcell.Screen) {
	// Layer order: background → terrain → water → units → particles → HUD
	w.bg.Draw(s, w.terrain)
	w.terrain.Draw(s)
	w.water.Draw(s, w.terrain)

	for _, u := range w.units {
		u.Draw(s, w.terrain)
	}
	for _, p := range w.particles {
		p.Draw(s)
	}

	w.hud.Draw(s, w.width, w.height, w.knightKills, w.robotKills, w.tick)
}

// FindNearestEnemy returns the closest alive enemy to u, or nil if none exist.
func (w *World) FindNearestEnemy(u *Unit) *Unit {
	var best *Unit
	bestDist := -1.0

	for _, other := range w.units {
		if other == u || other.isRobot == u.isRobot || other.IsDead() {
			continue
		}
		dx := u.cx() - other.cx()
		dy := (u.cy() - other.cy()) * 0.5 // weight horizontal more
		dist := dx*dx + dy*dy
		if bestDist < 0 || dist < bestDist {
			bestDist = dist
			best = other
		}
	}

	return best
}

func (w *World) AddParticle(p *Particle) {
	w.particles = append(w.particles, p)
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
