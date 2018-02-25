package main

import (
	"github.com/felixangell/strife"
)

type Player struct {
	x, y int
	size int
}

func (p *Player) tickAndRender(ctx *strife.Renderer) {
	if strife.KeyPressed(strife.KEY_W) {
		p.y -= 1
	}
	if strife.KeyPressed(strife.KEY_A) {
		p.x -= 1
	}
	if strife.KeyPressed(strife.KEY_S) {
		p.y += 1
	}
	if strife.KeyPressed(strife.KEY_D) {
		p.x += 1
	}

	ctx.SetColor(strife.White)
	ctx.Rect(p.x, p.y, p.size, p.size, strife.Fill)
}

type MyGame struct {
	p *Player
}

func (g *MyGame) tickAndRender(ctx *strife.Renderer) {
	g.p.tickAndRender(ctx)
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

	winWidth, winHeight := window.GetSize()
	playerSize := 64

	game := &MyGame{}
	game.p = &Player{
		size: playerSize,
		x:    (winWidth / 2) - (playerSize / 2),
		y:    (winHeight / 2) - (playerSize / 2),
	}

	for {
		window.PollEvents()
		if window.CloseRequested() {
			break
		}

		ctx := window.GetRenderContext()
		ctx.Clear()
		game.tickAndRender(ctx)
		ctx.Display()
	}
}
