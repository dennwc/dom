// +build wasm

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/svg"
)

func init() {
	dom.Body.Style().SetMarginsRaw("0")
}

func main() {
	fmt.Println("running")

	handler := func(e dom.Event) {
		dom.ConsoleLog(e)
		fmt.Printf("event: %T %v\n", e, e.JSValue())
	}

	inp := dom.Doc.NewInput("text")
	dom.Body.AppendChild(inp)
	inp.OnChange(handler)

	btn := dom.Doc.NewButton("Add")
	dom.Body.AppendChild(btn)

	const (
		w, h = 300, 300
		pad  = 12
	)
	root := svg.New(dom.Perc(100), dom.Px(h))

	center := root.NewG()
	center.Translate(w/2, h/2)
	center.NewCircle(10)

	type Sat struct {
		G       *svg.G
		HPeriod float64
	}

	var (
		mu   sync.Mutex
		sats []Sat
	)

	addSat := func(r, orb, hper float64, s string) {
		g := center.NewG()
		g.NewCircle(int(r)).Translate(orb, 0)
		g.NewLine().SetPos(0, 0, int(orb), 0)
		if s != "" {
			g.NewText(s).Translate(orb+pad, 0)
		}
		mu.Lock()
		sats = append(sats, Sat{G: g, HPeriod: hper})
		mu.Unlock()
	}

	btn.OnClick(func(_ dom.Event) {
		txt := inp.Value()
		r := 50 + rand.Float64()*75
		hper := 0.1 + rand.Float64()*3
		addSat(7, r, hper, txt)
	})

	addSat(5, 27, 1.5, "A")
	addSat(5, 40, 2.5, "B")

	start := time.Now()
	const interval = time.Millisecond * 30
	for {
		dt := time.Since(start).Seconds()
		tr := dt * 180

		for _, s := range sats {
			t := tr / s.HPeriod
			t -= float64(360 * int(t/360))
			s.G.Transform(svg.Rotate{A: t})
		}

		time.Sleep(interval)
	}
}
