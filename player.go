package main

import (
	"fmt"
	"math"
	"syscall/js"
)

type Player struct {
	Keys  map[string]bool
	Speed float64
	Rect
	FrameIndex, Tick int
	X, Y             float64
}

func (p *Player) Update(w *World, obstacles []Rect) {
	prevX, prevY := p.X, p.Y
	moving := p.handleMovement()

	p.clampToWorldBounds(w)

	for _, obstacle := range obstacles {
		if p.Rect.Intersects(obstacle) {
			Log(fmt.Sprintf("did Intersect: test person: %v obstacle: %v", p, obstacle))
			p.X, p.Y = prevX, prevY
			break
		}
	}

	p.updateAnimation(moving)
}

func (p *Player) handleMovement() bool {
	moving := false
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
	return moving
}

func (p *Player) clampToWorldBounds(w *World) {
	p.X = clamp(p.X, 0, w.Width-p.Rect.Width)
	p.Y = clamp(p.Y, 0, w.Height-p.Rect.Height)
}

func (p *Player) updateAnimation(moving bool) {
	if moving {
		p.Tick++

		// 8 is the frame duraction, and gives us ~7.5 FPS, good for this type of
		//   animation
		if p.Tick%8 == 0 {
			p.FrameIndex++
		}
	} else {
		p.FrameIndex = 0
	}
}

func (p *Player) Draw(ctx js.Value, cameraX, cameraY float64) {
	scale := 1.0 + 0.1*math.Sin(float64(p.FrameIndex)*0.5)
	size := p.Rect.Width * scale
	color := fmt.Sprintf("rgb(%d,%d,%d)", (p.FrameIndex*40)%255, (p.FrameIndex*85)%255, (p.FrameIndex*60)%255)

	ctx.Set("fillStyle", color)
	ctx.Call("fillRect",
		p.X-cameraX-(size-p.Rect.Width)/2,
		p.Y-cameraY-(size-p.Rect.Height)/2,
		size, size,
	)
}
