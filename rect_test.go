package main

import "testing"

func TestRectIntersects(t *testing.T) {
	tests := []struct {
		name     string
		rect1    Rect
		rect2    Rect
		expected bool
	}{
		{
			/*
				*+-----------+
				|  rect1    |
				|   +-----------+
				|   |   rect2   |
				|   |           |
				+---+-----------+
			*/
			name:     "overlapping rectangles",
			rect1:    Rect{Point: Point{X: 0, Y: 0}, Width: 10, Height: 10},
			rect2:    Rect{Point: Point{X: 5, Y: 5}, Width: 10, Height: 10},
			expected: true,
		},
		{
			/*
				+-----+          +-----+
				|rect1|          |rect2|
				+-----+          +-----+
			*/
			name:     "non-overlapping rectangles",
			rect1:    Rect{Point: Point{X: 0, Y: 0}, Width: 5, Height: 5},
			rect2:    Rect{Point: Point{X: 10, Y: 10}, Width: 5, Height: 5},
			expected: false,
		},
		{
			/*
			   +-----+-----+
			   |rect1|rect2|
			   +-----+-----+
			*/
			name:     "touching rectangles",
			rect1:    Rect{Point: Point{X: 0, Y: 0}, Width: 5, Height: 5},
			rect2:    Rect{Point: Point{X: 5, Y: 0}, Width: 5, Height: 5},
			expected: false,
		},
		{
			/*
			  +-----------+
			  |  rect1    |
			  |  +-----+  |
			  |  |rect2|  |
			  |  +-----+  |
			  +-----------+
			*/
			name:     "one rectangle inside another",
			rect1:    Rect{Point: Point{X: 0, Y: 0}, Width: 10, Height: 10},
			rect2:    Rect{Point: Point{X: 2, Y: 2}, Width: 5, Height: 5},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect1.Intersects(tt.rect2)
			if result != tt.expected {
				t.Errorf("Rect.Intersects() = %v, want %v\nRect1: %+v\nRect2: %+v",
					result, tt.expected, tt.rect1, tt.rect2)
			}
		})
	}
}
