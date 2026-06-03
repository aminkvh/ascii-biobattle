package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var (
	flagFPS       = flag.Int("speed", 30, "Animation speed in FPS (1-60)")
	flagDensity   = flag.Int("density", 3, "Units per side (1-80)")
	flagIntensity = flag.Int("intensity", 100, "Color intensity percent (10-100)")
	flagTheme     = flag.String("theme", "night", "Color theme: 'day' or 'night'")
	flagHelp      = flag.Bool("help", false, "Show help")
	// Internal flag used when we launch ourselves in a new terminal on Windows
	flagRun = flag.Bool("run", false, "")
)

func main() {
	// Handle Windows screensaver protocol args BEFORE flag.Parse
	if runtime.GOOS == "windows" {
		for _, arg := range os.Args[1:] {
			switch strings.ToLower(strings.TrimLeft(arg, "/-")) {
			case "s":
				launchWindowsScreensaver()
				return
			case "c":
				fmt.Println("ASCII Battle Screensaver — use --help for configuration options.")
				return
			case "p":
				// Preview: just run in current terminal
			}
		}
	}

	flag.Parse()

	if *flagHelp {
		fmt.Print(`ASCII Battle Screen Saver
=========================
A cross-platform terminal screensaver: medieval knights vs AI robots.

Usage:
  screensaver [flags]

Flags:
  --speed INT      FPS (default 30)
  --density INT    Units per side (default 3)
  --intensity INT  Color intensity % (default 100)
  --theme STRING   'night' or 'day' (default 'night')

Controls:
  q / Esc / Ctrl+C   Exit

Windows .scr arguments (used by Windows screensaver system):
  /s    Start screensaver
  /c    Configure (shows this help)
  /p    Preview

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

	runScreensaver()
}

// launchWindowsScreensaver opens a new maximized terminal window
// running this binary in screensaver mode.
func launchWindowsScreensaver() {
	exe, err := os.Executable()
	if err != nil {
		runScreensaver()
		return
	}

	// Try Windows Terminal first (wt.exe), then fall back to cmd.exe
	cmd := exec.Command("wt.exe", "--maximized", "--", exe, "--run")
	if err := cmd.Start(); err != nil {
		cmd = exec.Command("cmd.exe", "/c", "start", "/max", "cmd", "/k", exe, "--run")
		if err2 := cmd.Start(); err2 != nil {
			// Last resort: just run in-process
			runScreensaver()
		}
	}
}

func runScreensaver() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		s.Fini()
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
		}
	}()

	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	s.Clear()

	eng := NewEngine(s, *flagFPS, *flagDensity, *flagIntensity, *flagTheme)
	eng.Run()
}
