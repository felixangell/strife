package main

import (
	"math/rand"

	"github.com/felixangell/strife"
)

type Particle struct {
	x, y   float64
	dx, dy float64
	colour uint32
	s      int
}

func randFloat(max float64, min float64) float64 {
	return (rand.Float64() * (max - min)) + min
}

func main() {
	window := strife.SetupRenderWindow(1280, 720, strife.DefaultConfig())
	window.SetTitle("Hello world!")
	window.SetResizable(true)
	window.Create()

	window.HandleEvents(func(evt strife.StrifeEvent) {
		switch event := evt.(type) {
		case *strife.CloseEvent:
			println("closing window!")
			window.Close()
		case *strife.WindowResizeEvent:
			println("resize to ", event.Width, "x", event.Height)
		}
	})

	const numParticles int = 1000
	const speedCap float64 = 4.0

	w, h := window.GetSize()

	particles := [numParticles]*Particle{}
	for i := 0; i < numParticles; i++ {
		size := rand.Intn(32) + 8
		particles[i] = &Particle{
			x:      randFloat(0, float64(w-size)),
			y:      randFloat(0, float64(h-size)),
			dx:     randFloat(-(speedCap), speedCap),
			dy:     randFloat(-(speedCap), speedCap),
			s:      size,
			colour: uint32(rand.Intn(0xffffff) / 0xff),
		}
	}

	for {
		window.PollEvents()
		if window.CloseRequested() {
			break
		}

		ctx := window.GetRenderContext()
		ctx.SetColor(strife.White)
		ctx.Rect(0, 0, w, h, strife.Fill)

		{
			for _, p := range particles {
				p.x += p.dx
				p.y += p.dy
				if p.x < 0 || int(p.x) >= w-p.s {
					p.dx *= -1
				}
				if p.y < 0 || int(p.y) >= h-p.s {
					p.dy *= -1
				}

				ctx.SetColor(strife.HexRGB(p.colour))
				ctx.Rect(int(p.x), int(p.y), p.s, p.s, strife.Fill)
			}
		}
		ctx.Display()
	}
}
