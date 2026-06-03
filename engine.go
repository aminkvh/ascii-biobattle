package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// Engine manages the main tick loop and screen I/O.
type Engine struct {
	screen    tcell.Screen
	fps       int
	density   int
	intensity int
	theme     string
	world     *World
	quit      chan struct{}
}

func NewEngine(s tcell.Screen, fps, density, intensity int, theme string) *Engine {
	w, h := s.Size()
	return &Engine{
		screen:    s,
		fps:       fps,
		density:   density,
		intensity: intensity,
		theme:     theme,
		world:     NewWorld(w, h, density, intensity, theme),
		quit:      make(chan struct{}),
	}
}

func (e *Engine) Run() {
	go e.pollEvents()

	ticker := time.NewTicker(time.Second / time.Duration(e.fps))
	defer ticker.Stop()

	for {
		select {
		case <-e.quit:
			return
		case <-ticker.C:
			e.world.Update()
			e.screen.Clear()
			e.world.Draw(e.screen)
			e.screen.Show()
		}
	}
}

func (e *Engine) pollEvents() {
	for {
		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch {
			case ev.Key() == tcell.KeyEscape,
				ev.Key() == tcell.KeyCtrlC,
				ev.Rune() == 'q',
				ev.Rune() == 'Q':
				close(e.quit)
				return
			}
		case *tcell.EventResize:
			w, h := ev.Size()
			e.world.Resize(w, h)
			e.screen.Sync()
		}
	}
}
