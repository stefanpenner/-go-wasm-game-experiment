package main

import (
	"syscall/js"
)

const FrameDuration = 8

type World struct {
	Width, Height float64
}

type Canvas struct {
	Width, Height float64
}

func Log(message string) {
	js.Global().Get("console").Call("log", message)
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func handleKey(p *Player, pressed bool) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		key := args[0].Get("key").String()
		p.Keys[key] = pressed
		return nil
	})
}

func startGameLoop(w *World, c *Canvas, p *Player, ctx js.Value, obstacles []Rect) {
	var loop js.Func
	loop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		p.Update(w, obstacles)

		cameraX := clamp(p.X-c.Width/2+p.Rect.Width/2, 0, w.Width-c.Width)
		cameraY := clamp(p.Y-c.Height/2+p.Rect.Height/2, 0, w.Height-c.Height)

		ctx.Call("clearRect", 0, 0, c.Width, c.Height)

		ctx.Set("fillStyle", "gray")
		for _, obstacle := range obstacles {
			ctx.Call("fillRect", obstacle.X-cameraX, obstacle.Y-cameraY, obstacle.Width, obstacle.Height)
		}

		p.Draw(ctx, cameraX, cameraY)
		js.Global().Call("requestAnimationFrame", loop)
		return nil
	})

	js.Global().Call("requestAnimationFrame", loop)
}

func main() {
	world := &World{Width: 800, Height: 800}
	canvas := &Canvas{Width: 400, Height: 400}

	document := js.Global().Get("document")
	canvasElement := document.Call("getElementById", "the-canvas")
	ctx := canvasElement.Call("getContext", "2d")

	player := &Player{
		Rect:  Rect{X: 10, Y: 10, Width: 30, Height: 30},
		Speed: 3,
		Keys:  make(map[string]bool),
	}

	obstacles := []Rect{
		{X: 150, Y: 150, Width: 100, Height: 100},
		{X: 500, Y: 300, Width: 150, Height: 50},
		{X: 300, Y: 600, Width: 200, Height: 100},
	}

	document.Call("addEventListener", "keydown", handleKey(player, true))
	document.Call("addEventListener", "keyup", handleKey(player, false))

	startGameLoop(world, canvas, player, ctx, obstacles)
	select {}
}
