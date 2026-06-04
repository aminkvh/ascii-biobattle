package main

import (
	"fmt"
)

// HUD renders the battle statistics overlay.
type HUD struct {
	theme string
}

func NewHUD(theme string) *HUD {
	return &HUD{theme: theme}
}

func (h *HUD) Draw(s Screen, w, height, knightKills, robotKills, tick int) {
	// ── Top bar ─────────────────────────────────────────────────────
	topBg := NewRGBColor(0, 0, 0)
	knightCol := NewRGBColor(150, 180, 255) // blue-silver
	robotCol := NewRGBColor(255, 80, 30)    // orange-red
	sepCol := NewRGBColor(80, 80, 80)

	// Fill top bar background
	for x := 0; x < w; x++ {
		s.SetContent(x, 0, ' ', nil, StyleDefault.Background(topBg))
	}

	// Left: Knight score
	knightStr := fmt.Sprintf(" KNIGHTS  K:%05d ", knightKills)
	h.drawStr(s, 1, 0, knightStr, knightCol, topBg)

	// Centre: title
	title := "[ ASCII BATTLE SCREENSAVER ]"
	tx := (w - len(title)) / 2
	if tx > 0 {
		h.drawStr(s, tx, 0, title, NewRGBColor(0, 220, 180), topBg)
	}

	// Right: Robot score
	robotStr := fmt.Sprintf(" K:%05d  AI ROBOTS ", robotKills)
	rx := w - len(robotStr) - 1
	if rx > 0 {
		h.drawStr(s, rx, 0, robotStr, robotCol, topBg)
	}

	// ── Bottom status bar ──────────────────────────────────────────
	if height < 3 {
		return
	}
	by := height - 1
	for x := 0; x < w; x++ {
		s.SetContent(x, by, ' ', nil, StyleDefault.Background(topBg))
	}

	statusStr := fmt.Sprintf(" Tick: %06d  |  [Q/Esc] Exit  |  --speed --density --theme day/night ", tick)
	h.drawStr(s, 0, by, statusStr, sepCol, topBg)
}

func (h *HUD) drawStr(s Screen, x, y int, str string, fg, bg Color) {
	sw, sh := s.Size()
	style := StyleDefault.Foreground(fg).Background(bg)
	for _, ch := range str {
		if x >= sw || y >= sh {
			break
		}
		s.SetContent(x, y, ch, nil, style)
		x++
	}
}
