package main

import (
	"fmt"
	"math"
	"syscall/js"
)

type Keys map[string]bool

type Point struct {
	X float64
	Y float64
}

type Player struct {
	Keys
	Speed            float64
	FrameIndex, Tick int
	Rect
	Point
}

func (p *Player) Update(w *World, obstacles []Rect) {
	saved := p.Point
	moving := p.handleMovement()

	p.clampToWorldBounds(w)

	for _, obstacle := range obstacles {
		if p.Intersects(obstacle) {
			// there was an intersection, so we must restore as a collision did occure
			p.Point = saved
			break
		}
	}

	p.updateAnimation(moving)
}

func (p *Player) handleMovement() bool {
	moving := false
	if p.Keys["ArrowUp"] {
		p.Point.Y -= p.Speed
		moving = true
	}
	if p.Keys["ArrowDown"] {
		p.Point.Y += p.Speed
		moving = true
	}
	if p.Keys["ArrowLeft"] {
		p.Point.X -= p.Speed
		moving = true
	}
	if p.Keys["ArrowRight"] {
		p.Point.X += p.Speed
		moving = true
	}
	return moving
}

func (p *Player) clampToWorldBounds(w *World) {
	p.Point.X = clamp(p.Point.X, 0, w.Width-p.Width)
	p.Point.Y = clamp(p.Point.Y, 0, w.Height-p.Height)
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
	size := p.Width * scale
	color := fmt.Sprintf("rgb(%d,%d,%d)", (p.FrameIndex*40)%255, (p.FrameIndex*85)%255, (p.FrameIndex*60)%255)

	ctx.Set("fillStyle", color)
	ctx.Call("fillRect",
		p.Point.X-cameraX-(size-p.Width)/2,
		p.Point.Y-cameraY-(size-p.Height)/2,
		size, size,
	)
}
