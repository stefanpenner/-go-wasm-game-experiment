package main

type Point struct {
	X float64
	Y float64
}

type Rect struct {
	Width, Height float64
	Point
}

func (r *Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}
