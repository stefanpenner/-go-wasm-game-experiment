package main

import (
	"fmt"
	"math"
	"syscall/js"
)

const (
	CanvasWidth   = 400
	CanvasHeight  = 400
	WorldWidth    = 800
	WorldHeight   = 800
	FrameDuration = 8 // ~7.5 FPS, internet says this is good for this type of animation
)

type Rect struct {
	X, Y          float64
	Width, Height float64
}

func (r *Rect) Intersects(other *Rect) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

type Player struct {
	Keys  map[string]bool
	Speed float64
	Rect
	FrameIndex int
	Tick       int
}

func (p *Player) Update(obstacles []Rect) {
	// save snapshot
	prevX, prevY := p.X, p.Y
	moving := false

	// Movement logic
	if p.Keys["ArrowUp"] {
		p.Y -= p.Speed
		moving = true
	}
	if p.Keys["ArrowDown"] {
		p.Y += p.Speed
		moving = true
	}
	if p.Keys["ArrowLeft"] {
		p.X -= p.Speed
		moving = true
	}
	if p.Keys["ArrowRight"] {
		p.X += p.Speed
		moving = true
	}

	// World bounds
	if p.X < 0 {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = 0
	}
	if p.X > WorldWidth-p.Width {
		p.X = WorldWidth - p.Width
	}
	if p.Y > WorldHeight-p.Height {
		p.Y = WorldHeight - p.Height
	}

	// detect collisions
	for _, obstacle := range obstacles {
		if p.Intersects(&obstacle) {
			p.X, p.Y = prevX, prevY // restore
			break
		}
	}

	if moving {
		p.Tick++
		if p.Tick%FrameDuration == 0 {
			p.FrameIndex++
		}
	} else {
		// not moving
		p.FrameIndex = 0
	}
}

func (p *Player) Draw(ctx js.Value, cameraX, cameraY float64) {
	// AI gave this, it's cool
	// Procedural animation: simple "pulsing" effect
	scale := 1.0 + 0.1*math.Sin(float64(p.FrameIndex)*0.5)
	size := p.Width * scale

	color := "rgb(" +
		fmt.Sprintf("%d", (p.FrameIndex*40)%255) + "," +
		fmt.Sprintf("%d", (p.FrameIndex*85)%255) + "," +
		fmt.Sprintf("%d", (p.FrameIndex*60)%255) + ")"

	// Draw the "sprite"
	ctx.Set("fillStyle", color)
	ctx.Call("fillRect",
		p.X-cameraX-(size-p.Width)/2,
		p.Y-cameraY-(size-p.Height)/2,
		size, size)
}

func Log(message string) {
	// js.Global().Get("console").Call("log", "keyDown")
}

func handleKeyDown(p *Player) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Log("keyDown")
		event := args[0]
		key := event.Get("key").String()
		p.Keys[key] = true
		return nil
	})
}

func handleKeyUp(p *Player) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Log("keyUp")
		event := args[0]
		key := event.Get("key").String()
		p.Keys[key] = false
		return nil
	})
}

func startGameLoop(p *Player, ctx js.Value, obstacles []Rect) {
	var loop js.Func
	loop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Log("tik")
		// Update player position
		p.Update(obstacles)

		// Camera follows the player
		cameraX := p.X - CanvasWidth/2 + p.Width/2
		cameraY := p.Y - CanvasHeight/2 + p.Height/2

		// Keep camera within world bounds
		if cameraX < 0 {
			cameraX = 0
		}
		if cameraY < 0 {
			cameraY = 0
		}
		if cameraX > WorldWidth-CanvasWidth {
			cameraX = WorldWidth - CanvasWidth
		}
		if cameraY > WorldHeight-CanvasHeight {
			cameraY = WorldHeight - CanvasHeight
		}

		// Clear the canvas
		ctx.Call("clearRect", 0, 0, CanvasWidth, CanvasHeight)

		// Draw obstacles
		ctx.Set("fillStyle", "gray")
		for _, obstacle := range obstacles {
			ctx.Call("fillRect", obstacle.X-cameraX, obstacle.Y-cameraY, obstacle.Width, obstacle.Height)
		}

		// Draw player
		p.Draw(ctx, cameraX, cameraY)

		// Request the next frame
		js.Global().Call("requestAnimationFrame", loop)
		return nil
	})

	js.Global().Call("requestAnimationFrame", loop)
}

func main() {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "the-canvas")
	ctx := canvas.Call("getContext", "2d")

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

	document.Call("addEventListener", "keydown", handleKeyDown(player))
	document.Call("addEventListener", "keyup", handleKeyUp(player))

	startGameLoop(player, ctx, obstacles)

	select {}
}
