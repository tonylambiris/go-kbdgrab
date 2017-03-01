package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	keybind.Initialize(X)

	// Create a gradient window with random colors.
	geom := rootGeometry(X)
	newGradientWindow(X, geom.Width(), geom.Height(), newRandomColor(), newRandomColor())

	xevent.Main(X)
}
