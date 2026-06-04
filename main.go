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
	flagHelp      = flag.Bool("help", false, "Show help")
)

func main() {
	// Handle Windows screensaver protocol args BEFORE flag.Parse
	if runtime.GOOS == "windows" {
		for _, arg := range os.Args[1:] {
			switch strings.ToLower(strings.TrimLeft(arg, "/-")) {
			case "s":
				// Start Fullscreen
				flag.Parse()
				runScreensaver(true)
				return
			case "c":
				// Configure
				fmt.Println("ASCII Battle Screensaver — configuration not implemented natively yet.")
				return
			case "p":
				// Preview
				flag.Parse()
				runScreensaver(false)
				return
			}
		}
	}

	flag.Parse()

	if *flagHelp {
		fmt.Print(`ASCII Battle Screen Saver
=========================
A cross-platform native screensaver: medieval knights vs AI robots.

Usage:
  screensaver [flags]

Flags:
  --speed INT      FPS (default 30)
  --density INT    Units per side (default 3)
  --intensity INT  Color intensity % (default 100)
  --theme STRING   'night' or 'day' (default 'night')

Controls:
  q / Esc   Exit

Windows .scr arguments (used by Windows screensaver system):
  /s    Start screensaver (fullscreen)
  /c    Configure
  /p    Preview (windowed)
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

	// By default, just run as a regular windowed app
	// On Linux/Mac this allows it to be used as a standalone graphical app.
	runScreensaver(false)
}

func runScreensaver(fullscreen bool) {
	app := NewApp(*flagFPS, *flagDensity, *flagIntensity, *flagTheme)
	
	ebiten.SetWindowTitle("ASCII Biobattle")
	if fullscreen {
		ebiten.SetFullscreen(true)
	} else {
		ebiten.SetWindowSize(800, 600)
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	}
	// Hide cursor so it feels like a real screensaver
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(app); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
