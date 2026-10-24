package main

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed web/*
var assets embed.FS

// App exposes the small set of system controls needed by the desktop UI.
type App struct {
	inputMu sync.Mutex
}

var defaultWindowWidth, defaultWindowHeight = 1380, 1000
var minWindowWidth, minWindowHeight = 360, 320

func main() {
	if runtime.GOOS != "windows" {
		log.Fatal("FacePointer system pointer control is currently supported on Windows.")
	}
	if err := enableBackgroundTracking(); err != nil {
		log.Fatal(err)
	}
	screenWidth, screenHeight := screenSize()
	if screenWidth > 0 {
		defaultWindowWidth = min(1380, max(1, screenWidth*9/10))
		minWindowWidth = min(360, defaultWindowWidth)
	}
	if screenHeight > 0 {
		defaultWindowHeight = min(1000, max(1, screenHeight*9/10))
		minWindowHeight = min(320, defaultWindowHeight)
	}

	webAssets, err := fs.Sub(assets, "web")
	if err != nil {
		log.Fatal(err)
	}

	err = wails.Run(&options.App{
		Title:         "FacePointer",
		Width:         defaultWindowWidth,
		Height:        defaultWindowHeight,
		MinWidth:      minWindowWidth,
		MinHeight:     minWindowHeight,
		DisableResize: false,
		AssetServer:   &assetserver.Options{Assets: webAssets},
		Bind:          []interface{}{&App{}},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func enableBackgroundTracking() error {
	const variable = "WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS"
	args := strings.TrimSpace(os.Getenv(variable))
	for _, flag := range []string{
		"--disable-background-timer-throttling",
		"--disable-renderer-backgrounding",
		"--disable-backgrounding-occluded-windows",
	} {
		if !strings.Contains(args, flag) {
			if args != "" {
				args += " "
			}
			args += flag
		}
	}
	return os.Setenv(variable, args)
}

// GetStatus returns the desktop screen size for the pointer display.
func (a *App) GetStatus() map[string]any {
	width, height := screenSize()
	virtualX, virtualY, virtualWidth, virtualHeight := virtualScreenBounds()
	return map[string]any{
		"name": "FacePointer", "screenWidth": width, "screenHeight": height,
		"virtualX": virtualX, "virtualY": virtualY,
		"virtualWidth": virtualWidth, "virtualHeight": virtualHeight,
		"windowWidth": defaultWindowWidth, "windowHeight": defaultWindowHeight,
		"doubleClickTime": doubleClickInterval(), "platform": runtime.GOOS,
	}
}

// Move moves the Windows pointer to normalized screen coordinates.
func (a *App) Move(x, y float64) error {
	if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) || x < 0 || x > 1 || y < 0 || y > 1 {
		return errors.New("cursor coordinates must be between 0 and 1")
	}
	left, top, width, height := virtualScreenBounds()
	if width < 1 || height < 1 {
		return errors.New("could not read the display size")
	}
	a.inputMu.Lock()
	defer a.inputMu.Unlock()
	return moveCursor(left+int(math.Round(x*float64(width-1))), top+int(math.Round(y*float64(height-1))))
}

// Action performs a supported mouse or keyboard action.
func (a *App) Action(action string) error {
	if action != "left" && action != "left-down" && action != "left-up" && action != "double-left" && action != "right" && action != "scroll-up" && action != "scroll-down" && action != "tab" && action != "next-tab" {
		return errors.New("unsupported pointer action")
	}
	a.inputMu.Lock()
	defer a.inputMu.Unlock()
	return performAction(action)
}
