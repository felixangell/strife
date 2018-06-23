package main

import (
	"github.com/felixangell/strife"
)

func render(ctx *strife.Renderer) {
	ctx.SetColor(strife.Black)
	ctx.Text("Hello World!", 50, 49)

	ctx.SetColor(strife.Red)
	ctx.Text("Hello World!", 50, 50)
}

func main() {
	window := strife.SetupRenderWindow(800, 600, strife.DefaultConfig())
	window.SetTitle("Text stuff!")
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

	ctx := window.GetRenderContext()
	bigger, _ := ctx.GetFont().DeriveFont(92)
	ctx.SetFont(bigger)

	for {
		window.PollEvents()
		if window.CloseRequested() {
			break
		}

		ctx.Clear()
		{
			ctx.SetColor(strife.White)
			ctx.Rect(0, 0, 800, 600, strife.Fill)

			render(ctx)
		}
		ctx.Display()
	}
}
