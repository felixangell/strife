# strife
A simple game framework that wraps around SDL2.

## example
The biggest example use of the Strife framework is the [Phi text editor](//phi.felixangell.com)

<p align="center"><img src="https://raw.githubusercontent.com/felixangell/phi/master/screenshot.png"></p>

Though there are some smaller examples demonstrating components of the Strife API in the `examples/` folder.

## note/disclaimer
This is a work in progress. It provides a very minimal toolset for rendering shapes, images, and text
as well as capturing user input. This is not at a production level and is mostly being worked on when the
needs of my other projects (that depend on this) evolve.

There is no documentation either! If you want to use it you will have to check out the examples. I may
get round to writing some documentation but the API is very volatile at the moment.

## installing
Simple as

	$ go get github.com/felixangell/strife

Make sure you have SDL2 installed as well as the ttf and img addons:

	$ brew install SDL2 SDL2_ttf SDL2_img

## getting started
This is a commented code snippet to help you get started:

```go
func main() {
	// create a nice shiny window
	window := strife.SetupRenderWindow(1280, 720, strife.DefaultConfig())
	window.SetTitle("Hello world!")
	window.SetResizable(true)
	window.Create()

	// this is our event handler
	window.HandleEvents(func(evt strife.StrifeEvent) {
		switch event := evt.(type) {
		case *strife.CloseEvent:
			println("closing window!")
			window.Close()
		case *strife.WindowResizeEvent:
			println("resize to ", event.Width, "x", event.Height)
		}
	})

	// game loop
	for {
		// handle the events before we do any
		// rendering etc.
		window.PollEvents()

		// if we have a window close event
		// from the previous poll, break out of
		// our game loop
		if window.CloseRequested() {
			break
		}

		// rendering context stuff here
		// clear and display is typical
		// all your rendering code should
		// go...
		ctx := window.GetRenderContext()
		ctx.Clear()
		{
			// ...in this section here
			ctx.SetColor(strife.Red)
			ctx.Rect(10, 10, 50, 50, strife.Fill)

			// check out some other examples!
		}
		ctx.Display()
	}
}
```

## license
[MIT](/LICENSE)
