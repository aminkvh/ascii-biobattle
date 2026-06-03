package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// HUD renders the battle statistics overlay.
type HUD struct {
	theme string
}

func NewHUD(theme string) *HUD {
	return &HUD{theme: theme}
}

func (h *HUD) Draw(s tcell.Screen, w, height, knightKills, robotKills, tick int) {
	// ── Top bar ─────────────────────────────────────────────────────
	topBg := tcell.NewRGBColor(0, 0, 0)
	knightCol := tcell.NewRGBColor(150, 180, 255) // blue-silver
	robotCol := tcell.NewRGBColor(255, 80, 30)    // orange-red
	sepCol := tcell.NewRGBColor(80, 80, 80)

	// Fill top bar background
	for x := 0; x < w; x++ {
		s.SetContent(x, 0, ' ', nil, tcell.StyleDefault.Background(topBg))
	}

	// Left: Knight score
	knightStr := fmt.Sprintf(" KNIGHTS  K:%05d ", knightKills)
	h.drawStr(s, 1, 0, knightStr, knightCol, topBg)

	// Centre: title
	title := "[ ASCII BATTLE SCREENSAVER ]"
	tx := (w - len(title)) / 2
	if tx > 0 {
		h.drawStr(s, tx, 0, title, tcell.NewRGBColor(0, 220, 180), topBg)
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
		s.SetContent(x, by, ' ', nil, tcell.StyleDefault.Background(topBg))
	}

	statusStr := fmt.Sprintf(" Tick: %06d  |  [Q/Esc] Exit  |  --speed --density --theme day/night ", tick)
	h.drawStr(s, 0, by, statusStr, sepCol, topBg)
}

func (h *HUD) drawStr(s tcell.Screen, x, y int, str string, fg, bg tcell.Color) {
	sw, sh := s.Size()
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	for _, ch := range str {
		if x >= sw || y >= sh {
			break
		}
		s.SetContent(x, y, ch, nil, style)
		x++
	}
}
