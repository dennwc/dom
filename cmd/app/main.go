// +build wasm

package main

import (
	"fmt"
	"github.com/dennwc/dom"
	"github.com/dennwc/dom/svg"
	"time"
)

func main() {
	fmt.Println("running")

	const (
		w, h = 300, 300
		r2   = 27
		r3   = 40
	)
	root := svg.New(dom.Perc(100), dom.Vh(99))

	center := root.NewG()
	center.Translate(w/2, h/2)
	center.NewCircle(10)

	g1 := center.NewG()
	g1.NewCircle(5).Translate(r2, 0)
	g1.NewLine().SetPos(0, 0, r2, 0)

	g2 := center.NewG()
	g2.NewCircle(5).Translate(r3, 0)
	g2.NewLine().SetPos(0, 0, r3, 0)
	g2.NewText("A2").Translate(r3, 0)

	start := time.Now()
	const interval = time.Millisecond * 30
	for {
		dt := time.Since(start).Seconds()
		tr := dt * 180

		t := tr / 1.5
		t -= float64(360 * int(t/360))
		g1.Transform(svg.Rotate{A: t})

		t = tr / 2.5
		t -= float64(360 * int(t/360))
		g2.Transform(svg.Rotate{A: t})

		time.Sleep(interval)
	}
}
