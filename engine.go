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
	world     *World
	grid      *GridScreen
	fps       int
	density   int
	intensity int
	theme     string

	// Window dimensions
	winW, winH int
	cols, rows int

	tickCount int
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

	// 1. Ask world to draw to our abstract grid
	a.grid.Clear()
	a.world.Draw(a.grid)

	// 2. Render the grid to the Ebiten image
	cols, rows := a.grid.Size()
	
	// Create text options
	op := &text.DrawOptions{}
	op.LineSpacing = charHeight

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			cell := a.grid.GetCell(x, y)
			
			// Draw background rect
			if cell.Style.bg != ColorBlack { 
				vector.DrawFilledRect(screen, float32(x)*charWidth, float32(y)*charHeight, charWidth, charHeight, cell.Style.bg.ToRGBA(), false)
			}

			// Draw character
			if cell.Primary != ' ' {
				op.GeoM.Reset()
				op.GeoM.Translate(float64(x)*charWidth, float64(y)*charHeight-3.0) // tweak baseline
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
		
		a.cols = int(float64(outsideWidth) / charWidth)
		a.rows = int(float64(outsideHeight) / charHeight)
		
		if a.cols < 10 { a.cols = 10 }
		if a.rows < 10 { a.rows = 10 }

		a.grid.Resize(a.cols, a.rows)
		if a.world == nil {
			a.world = NewWorld(a.cols, a.rows, a.density, a.intensity, a.theme)
		} else {
			a.world.Resize(a.cols, a.rows)
		}
	}
	return outsideWidth, outsideHeight
}
