package main

import (
	"math/rand"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	log "github.com/sirupsen/logrus"
)

const (
	Name = "kbdgrab"
	Path = "github.com/tonylambiris/go-kbdgrab"
)

var (
	// compile-time variables
	Build string
	Type  string

	// The size of the text (scaled dynamically below)
	size = 0.0

	// The text to draw.
	msg = "Capturing keyboard input, type CTRL-ESC to exit."
)

func main() {
	log.Printf("%s %s (%s) [%s]", Name, Build, Type, Path)

	rand.Seed(time.Now().UnixNano())

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Create a gradient window with random colors.
	geom := rootGeometry(X)
	size = float64(geom.Height()) / float64(geom.Width()) * 100

	newGradientWindow(X, geom.Width(), geom.Height(),
		newRandomColor(), newRandomColor())

	xevent.Main(X)
}
