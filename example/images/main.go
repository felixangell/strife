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
			println("HUR DUR WE CLOSING")
			window.Close()
		case *strife.WindowResizeEvent:
			println("resize to ", event.Width, "x", event.Height)
		}
	})

	masterpiece, err := strife.LoadImage("./res/masterpiece.png")
	if err != nil {
		panic(err)
	}
	scaledW := masterpiece.Width / 4
	scaledH := masterpiece.Height / 4

	var x, y int
	var dx, dy float64 = 6, 6

	for {
		window.PollEvents()
		if window.CloseRequested() {
			break
		}

		w, h := window.GetSize()

		dx *= 0.9999
		dy *= 0.9999

		x += int(dx)
		y += int(dy)

		if x > w-scaledW || x < 0 {
			dx *= -1
		}
		if y > h-scaledH || y < 0 {
			dy *= -1
		}

		ctx := window.GetRenderContext()
		ctx.Clear()
		{
			ctx.SetColor(strife.Red)
			ctx.Rect(50, 50, 50, 50, strife.Fill)

			ctx.ImageScale(masterpiece, x, y, scaledW, scaledH)

			// renders some arbitrary section of the image
			ctx.SubImage(masterpiece, 500, 40, 50, 50, 90, 40)
		}
		ctx.Display()
	}

	masterpiece.Destroy()
}
