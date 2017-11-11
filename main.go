package main

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

var (
	// The size of the text (scaled dynamically below)
	size = 0.0

	// The text to draw.
	msg = "Capturing keyboard input, type CTRL-ESC to exit."
)

func main() {
	//fmt.Sprintf("%s (%s build)", main.Build, main.Type)
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
