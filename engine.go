package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/gomono"
)

var faceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(gomono.TTF))
	if err != nil {
		log.Fatal(err)
	}
	faceSource = s
}

const (
	charWidth  = 10.0 // approx pixels per char
	charHeight = 20.0
)

type App struct {
	world           *World
	grid            *GridScreen
	fps             int
	density         int
	intensity       int
	theme           string
	screensaverMode bool

	// Window dimensions (this monitor only)
	winW, winH int
	cols, rows int

	// Multi-monitor: the world is wider than this monitor.
	// viewportX is the column offset into the world grid for this monitor.
	// worldCols is the total world width (all monitors combined).
	viewportX int
	worldCols int
	seed      int64

	tickCount int

	// For screensaver mouse-exit detection
	lastMouseX, lastMouseY int
	mouseInitialized       bool
}

func NewApp(fps, density, intensity int, theme string) *App {
	return &App{
		fps:       fps,
		density:   density,
		intensity: intensity,
		theme:     theme,
		grid:      NewGridScreen(80, 24),
	}
}

func (a *App) Update() error {
	// In screensaver mode: exit on any key OR mouse movement.
	// But give a 2-second grace period (120 ticks at 60fps) for the
	// window to fully appear before tracking mouse position.
	if a.screensaverMode {
		if a.tickCount < 120 {
			// Grace period: just record current mouse position, don't exit
			a.lastMouseX, a.lastMouseY = ebiten.CursorPosition()
			a.mouseInitialized = true
		} else {
			mx, my := ebiten.CursorPosition()
			if !a.mouseInitialized {
				a.lastMouseX = mx
				a.lastMouseY = my
				a.mouseInitialized = true
			} else {
				dx := mx - a.lastMouseX
				dy := my - a.lastMouseY
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}
				if dx > 10 || dy > 10 {
					return ebiten.Termination
				}
			}
		}
	}

	// ebiten updates 60 times a second by default.
	// We want to update at `fps` times a second.
	a.tickCount++
	targetWait := 60 / a.fps
	if targetWait < 1 {
		targetWait = 1
	}

	if a.tickCount%targetWait == 0 {
		if a.world != nil {
			a.world.Update()
		}
	}

	// Handle exit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	if a.world == nil {
		return
	}

	// 1. Ask world to draw to the FULL-WIDTH grid (all monitors)
	a.grid.Clear()
	a.world.Draw(a.grid)

	// 2. Render only THIS monitor's viewport slice to the Ebiten screen
	_, rows := a.grid.Size()

	op := &text.DrawOptions{}
	op.LineSpacing = charHeight

	for y := 0; y < rows; y++ {
		for screenX := 0; screenX < a.cols; screenX++ {
			// Map screen column to world column
			worldX := screenX + a.viewportX
			cell := a.grid.GetCell(worldX, y)

			// Draw background rect
			if cell.Style.bg != ColorBlack {
				vector.DrawFilledRect(screen, float32(screenX)*charWidth, float32(y)*charHeight, charWidth, charHeight, cell.Style.bg.ToRGBA(), false)
			}

			// Draw character
			if cell.Primary != ' ' {
				op.GeoM.Reset()
				op.GeoM.Translate(float64(screenX)*charWidth, float64(y)*charHeight-3.0)
				op.ColorScale.Reset()
				op.ColorScale.ScaleWithColor(cell.Style.fg.ToRGBA())

				text.Draw(screen, string(cell.Primary), &text.GoTextFace{
					Source: faceSource,
					Size:   18,
				}, op)
			}
		}
	}
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != a.winW || outsideHeight != a.winH || a.world == nil {
		a.winW = outsideWidth
		a.winH = outsideHeight

		// This monitor's columns
		a.cols = int(float64(outsideWidth) / charWidth)
		a.rows = int(float64(outsideHeight) / charHeight)

		if a.cols < 10 {
			a.cols = 10
		}
		if a.rows < 10 {
			a.rows = 10
		}

		// Total world width: if worldCols was set (multi-monitor), use it.
		// Otherwise just use this monitor's cols.
		totalCols := a.cols
		if a.worldCols > 0 {
			totalCols = a.worldCols
		}

		// Grid is the FULL world width (so subsystems draw across all monitors)
		a.grid.Resize(totalCols, a.rows)
		if a.world == nil {
			a.world = NewWorldWithSeed(totalCols, a.rows, a.density, a.intensity, a.theme, a.seed)
		} else {
			a.world.Resize(totalCols, a.rows)
		}
	}
	return outsideWidth, outsideHeight
}
