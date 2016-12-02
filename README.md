# strife
A simple wrapper around veandco's SDL2 bindings. 

I wrote this out of frustration with using SDL2. It's a bit awkward to use as you need a renderer instance to draw images, etc. The code isn't pretty, but it should simplify usage to make it easy to render things intuitively.

## notes
I pretty much copy/pasted this from one of my game projects into a separate repository, so it's all bundled into one module. Maybe I'll clean it up later, but hey ho, it works for now.

This is pretty feature-less right now. I'll be expanding it as my needs progress in a little project I'm working on.

## example
Here's a little example:

	package main

	import (
		"github.com/felixangell/strife/gfx"
	)

	func main() {
		window, failed := gfx.CreateRenderWindow(1280, 720)
		if failed {
			panic("failed to create render window")
		}

		for !window.CloseRequested() {
			ctx := window.GetRenderContext()
			ctx.Clear()

			{
				ctx.SetColor(gfx.RGB(255, 0, 255))
				ctx.Rect(10, 10, 50, 50, gfx.Fill)			
			}

			ctx.Display()
		}
	}

## license
[MIT](/LICENSE)