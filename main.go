package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	flagFPS       = flag.Int("speed", 30, "Animation speed in FPS (1-60)")
	flagDensity   = flag.Int("density", 3, "Units per side (1-80)")
	flagIntensity = flag.Int("intensity", 100, "Color intensity percent (10-100)")
	flagTheme     = flag.String("theme", "night", "Color theme: 'day' or 'night'")
	flagWindowed  = flag.Bool("windowed", false, "Run in a window instead of fullscreen (for development)")
	flagMonitor   = flag.Int("monitor", -1, "Monitor index to go fullscreen on (0=primary, 1, 2...)")
	flagSeed      = flag.Int64("seed", 0, "Random seed for world generation (same seed = same world)")
	flagWorldCols = flag.Int("world-cols", 0, "Total world width in columns (for multi-monitor continuous mode)")
	flagViewportX = flag.Int("viewport-x", 0, "Column offset into the world for this monitor's viewport")
	flagListMonitors = flag.Bool("list-monitors", false, "List connected monitors and exit")
	flagHelp      = flag.Bool("help", false, "Show help")
)

func main() {
	// Handle Windows screensaver protocol args BEFORE flag.Parse
	if runtime.GOOS == "windows" {
		for _, arg := range os.Args[1:] {
			switch strings.ToLower(strings.TrimLeft(arg, "/-")) {
			case "s":
				flag.Parse()
				runScreensaver(true)
				return
			case "c":
				fmt.Println("ASCII Battle Screensaver — configuration not implemented natively yet.")
				return
			case "p":
				flag.Parse()
				runScreensaver(false)
				return
			}
		}
	}

	flag.Parse()

	if *flagListMonitors {
		monitors := ebiten.AppendMonitors(nil)
		for i, m := range monitors {
			w, h := m.Size()
			fmt.Printf("Monitor %d: Name=%q, Size=%dx%d, Scale=%.2f\n", i, m.Name(), w, h, m.DeviceScaleFactor())
		}
		return
	}

	if *flagHelp {
		fmt.Print(`ASCII Battle Screen Saver
=========================
A cross-platform native screensaver: medieval knights vs AI robots.

Usage:
  ascii-biobattle [flags]

Flags:
  --speed INT        FPS (default 30)
  --density INT      Units per side (default 3)
  --intensity INT    Color intensity % (default 100)
  --theme STRING     'night' or 'day' (default 'night')
  --windowed         Run in a window (for development)
  --monitor INT      Monitor index for fullscreen (0=primary, 1, 2...)
  --seed INT         Random seed (same seed = identical world across monitors)
  --world-cols INT   Total world width in columns (multi-monitor continuous mode)
  --viewport-x INT   Column offset for this monitor's viewport

Controls:
  q / Esc / mouse move   Exit (in screensaver mode)
  q / Esc                Exit (in windowed mode)
`)
		return
	}

	if *flagFPS < 1 {
		*flagFPS = 1
	}
	if *flagFPS > 60 {
		*flagFPS = 60
	}
	if *flagDensity < 1 {
		*flagDensity = 1
	}
	if *flagDensity > 80 {
		*flagDensity = 80
	}
	if *flagIntensity < 10 {
		*flagIntensity = 10
	}
	if *flagIntensity > 100 {
		*flagIntensity = 100
	}

	runScreensaver(!*flagWindowed)
}

func runScreensaver(fullscreen bool) {
	app := NewApp(*flagFPS, *flagDensity, *flagIntensity, *flagTheme)
	app.screensaverMode = fullscreen
	app.seed = *flagSeed
	app.worldCols = *flagWorldCols
	app.viewportX = *flagViewportX

	ebiten.SetWindowTitle("ASCII Biobattle")

	if fullscreen {
		// Select the target monitor
		monIdx := *flagMonitor
		monitors := ebiten.AppendMonitors(nil)
		if monIdx >= 0 && monIdx < len(monitors) {
			ebiten.SetMonitor(monitors[monIdx])
		}

		ebiten.SetFullscreen(true)
		ebiten.SetWindowDecorated(false)
	} else {
		ebiten.SetWindowSize(1280, 720)
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	}

	// Hide cursor so it feels like a real screensaver
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(app); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
