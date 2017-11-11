package main

// Example window-gradient demonstrates how to create several windows and draw
import (
	"bytes"
	"image"
	"image/color"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// newGradientWindow creates a new X window, paints the initial gradient
// image, and listens for ConfigureNotify events. (A new gradient image must
// be painted in response to each ConfigureNotify event, since a
// ConfigureNotify event corresponds to a change in the window's geometry.)
func newGradientWindow(X *xgbutil.XUtil, width, height int,
	start, end color.RGBA) {

	wid := X.RootWin()

	keybind.Initialize(X)
	if err := keybind.GrabKeyboard(X, wid); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Grabbing keyboard input...")
	}

	// Generate a new window id.
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}

	// Create the window and die if it fails.
	if err := win.CreateChecked(wid, 0, 0, width, height, 0); err != nil {
		log.Fatal(err)
	}

	// Get EventMask events
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease)

	go func() {
		for {
			// Paint the initial gradient to the window and then map the window.
			renderGradient(X, win.Id, width, height, start, end)

			time.Sleep(1 * time.Second)
		}
	}()

	win.Map()

	win.Listen(xproto.EventMaskKeyPress)

	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			// keybind.LookupString does the magic of implementing parts of
			// the X Keyboard Encoding to determine an english representation
			// of the modifiers/keycode tuple.
			// N.B. It's working for me, but probably isn't 100% correct in
			// all environments yet.
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			if len(modStr) > 0 {
				log.Printf("Key: %s-%s", modStr, keyStr)
			} else {
				log.Printf("Key: %s", keyStr)
			}

			if keybind.KeyMatch(X, "Escape", e.State, e.Detail) {
				if e.State&xproto.ModMaskControl > 0 {
					log.Println("Control-Escape detected. Quitting...")
					xevent.Quit(X)
				}
			}
		}).Connect(X, wid)

	if err = ewmh.WmStateReq(X, win.Id, ewmh.StateToggle,
		"_NET_WM_STATE_FULLSCREEN"); err != nil {
		log.Fatal(err)
	}

	if err = ewmh.WmStateReq(X, win.Id, ewmh.StateToggle,
		"_NET_WM_STATE_ABOVE"); err != nil {
		log.Fatal(err)
	}
}

// renderGradient creates a new xgraphics.Image value and draws a gradient
// starting at color 'start' and ending at color 'end'.
//
// Since xgraphics.Image values use pixmaps and pixmaps cannot be resized,
// a new pixmap must be allocated for each resize event.
func renderGradient(X *xgbutil.XUtil, wid xproto.Window, width, height int,
	start, end color.RGBA) {

	ximg := xgraphics.New(X, image.Rect(0, 0, width, height))

	// Now calculate the increment step between each RGB channel in
	// the start and end colors.
	rinc := (0xff * (int(end.R) - int(start.R))) / width
	ginc := (0xff * (int(end.G) - int(start.G))) / width
	binc := (0xff * (int(end.B) - int(start.B))) / width

	// Now apply the increment to each "column" in the image.
	// Using 'ForExp' allows us to bypass the creation of a color.BGRA value
	// for each pixel in the image.
	ximg.ForExp(func(x, y int) (uint8, uint8, uint8, uint8) {
		return uint8(int(start.B) + (binc*x)/0xff),
			uint8(int(start.G) + (ginc*x)/0xff),
			uint8(int(start.R) + (rinc*x)/0xff),
			0xff
	})

	// Set the surface to paint on for ximg.
	// (This simply sets the background pixmap of the window to the pixmap
	// used by ximg.)
	ximg.XSurfaceSet(wid)

	// Render the message text
	renderText(ximg, rand.Intn(width/3), rand.Intn(height-100))

	// XDraw will draw the contents of ximg to its corresponding pixmap.
	ximg.XDraw()

	// XPaint will "clear" the window provided so that it shows the updated
	// pixmap.
	ximg.XPaint(wid)

	// Since we will not reuse ximg, we must destroy its pixmap.
	ximg.Destroy()
}

func renderText(ximg *xgraphics.Image, x, y int) {
	// Load the font.
	fontData, err := Asset("LCD_Solid.ttf")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewReader(fontData)

	// Now parse the font.
	font, err := xgraphics.ParseFont(buf)
	if err != nil {
		log.Fatal(err)
	}

	// Now draw some text
	_, _, err = ximg.Text(x, y, newRandomColor(), size, font, msg)
	if err != nil {
		log.Fatal(err)
	}

	// Now compute extents of the line of text
	secw, sech := xgraphics.Extents(font, size, msg)

	// Now repaint on the region that we drew text on. Then update the screen.
	bounds := image.Rect(x, y, x+secw, y+sech)

	ximg.SubImage(bounds).(*xgraphics.Image).XDraw()
}

// newRandomColor creates a new RGBA color where each channel (except alpha)
// is randomly generated.
func newRandomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 0xff,
	}
}

// rawGeometry isn't smart. It just queries the window given for geometry.
func rawGeometry(xu *xgbutil.XUtil, win xproto.Drawable) (xrect.Rect, error) {
	xgeom, err := xproto.GetGeometry(xu.Conn(), win).Reply()
	if err != nil {
		return nil, err
	}

	return xrect.New(int(xgeom.X), int(xgeom.Y),
		int(xgeom.Width), int(xgeom.Height)), nil
}

// rootGeometry gets the geometry of the root window. It will panic on failure.
func rootGeometry(xu *xgbutil.XUtil) xrect.Rect {
	geom, err := rawGeometry(xu, xproto.Drawable(xu.RootWin()))
	if err != nil {
		log.Fatal(err)
	}

	return geom
}
