// Example get-active-window reads the _NET_ACTIVE_WINDOW property of the root
// window and uses the result (a window id) to get the name of the window.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
)

type Monitor struct {
	X      int16
	Y      int16
	Width  uint16
	Height uint16
}

func getCurrentMonitor(X *xgb.Conn, root xproto.Window) (Monitor, error) {
	err := xinerama.Init(X)
	if err != nil {
		return Monitor{}, fmt.Errorf("Failed to start Xinerama: %v", err)
	}

	isActive, err := xinerama.IsActive(X).Reply()
	if err != nil {
		return Monitor{}, fmt.Errorf("Failed to identify if Xinerama is active: %v", err)
	}

	// if it is not active, return fullscreen
	if isActive == nil {
		setup := xproto.Setup(X)
		screen := setup.DefaultScreen(X)
		return Monitor{
			X:      0,
			Y:      0,
			Width:  screen.WidthInPixels,
			Height: screen.HeightInPixels,
		}, nil
	}

	pointerReply, err := xproto.QueryPointer(X, root).Reply()
	if err != nil {
		return Monitor{}, fmt.Errorf("Failed to query pointer: %v", err)
	}

	cursorX, cursorY := pointerReply.RootX, pointerReply.RootY

	screens, err := xinerama.QueryScreens(X).Reply()
	if err != nil {
		return Monitor{}, fmt.Errorf("Failed to query screens: %v", err)
	}

	for i, screen := range screens.ScreenInfo {
		if int16(screen.XOrg) <= cursorX &&
			cursorX < (screen.XOrg+int16(screen.Width)) && // width is the width of the current monitor only
			int16(screen.YOrg) <= cursorY &&
			cursorY < (screen.YOrg+int16(screen.Height)) {
			log.Printf("Cursor is on monitor %d: %dx%d at position (%d,%d)\n",
				i, screen.Width, screen.Height, cursorX, cursorY)

			return Monitor{
				X:      int16(screen.XOrg),
				Y:      int16(screen.YOrg),
				Width:  screen.Width,
				Height: screen.Height,
			}, nil
		}

	}

	// if can't find, use the first one
	if len(screens.ScreenInfo) > 0 {
		log.Println("Failed to find the proper monitor where the cursor is... using the first one")
		screen := screens.ScreenInfo[0]
		return Monitor{
			X:      int16(screen.XOrg),
			Y:      int16(screen.YOrg),
			Width:  screen.Width,
			Height: screen.Height,
		}, nil
	}

	// last resource: use fullscreen
	log.Println("Failed to find any monitor. Using fullscreen.")
	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)
	return Monitor{
		X:      0,
		Y:      0,
		Width:  screen.WidthInPixels,
		Height: screen.HeightInPixels,
	}, nil
}

func createNewImage(imgCookie *xproto.GetImageReply, width int, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			idx := (y*int(width) + x) * 4 // 4 bytes per pixel
			if idx+3 < len(imgCookie.Data) {
				b := imgCookie.Data[idx]
				g := imgCookie.Data[idx+1]
				r := imgCookie.Data[idx+2]
				a := uint8(255)
				img.SetRGBA(x, y, color.RGBA{r, g, b, a})
			}
		}

	}

	file, err := os.Create("screenshot.png")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		log.Fatal("Failed to save image:", err)
	}

	log.Println("Saved screenshot to screenshot.png")
	return nil
}

func main() {
	X, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)
	root := screen.Root
	width, height := screen.WidthInPixels, screen.HeightInPixels

	fmt.Println("width, height: ", width, height)

	mon, err := getCurrentMonitor(X, root)
	if err != nil {
		log.Fatal("Failed to get monitor: ", err)
	}

	imgCookie, err := xproto.GetImage(
		X,
		xproto.ImageFormatZPixmap,
		xproto.Drawable(root),
		mon.X,
		mon.Y,
		mon.Width,
		mon.Height,
		0xffffff,
	).Reply()

	if err != nil {
		log.Fatal("Failed to get screen image: ", err)
	}

	createNewImage(imgCookie, int(mon.Width), int(mon.Height))
}
