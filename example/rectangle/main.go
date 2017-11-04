package main

import (
	"github.com/felixangell/strife"
)

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

	for {
		window.PollEvents()
		if window.CloseRequested() {
			break
		}

		ctx := window.GetRenderContext()
		ctx.Clear()
		{
			ctx.SetColor(strife.Red)
			ctx.Rect(10, 10, 50, 50, strife.Fill)
		}
		ctx.Display()
	}
}
