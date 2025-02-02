package main

import (
	"testing"
)

func TestPlayerUpdate(t *testing.T) {
	tests := []struct {
		obstacles []Rect
		name      string
		player    Player
		world     World
		want      Point
	}{
		{
			name: "move right without obstacles",
			player: Player{
				Point: Point{
					X: 0,
					Y: 0,
				},
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  Keys{"ArrowRight": true},
			},
			world:     World{Width: 100, Height: 100},
			obstacles: []Rect{},
			want: Point{
				X: 5,
				Y: 0,
			},
		},
		{
			name: "blocked by obstacle",
			player: Player{
				Point: Point{
					X: 0,
					Y: 0,
				},
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  Keys{"ArrowRight": true},
			},
			world: World{Width: 100, Height: 100},
			obstacles: []Rect{
				{X: 5, Y: 0, Width: 10, Height: 10},
			},
			want: Point{
				X: 0,
				Y: 0,
			},
		},
		{
			name: "world boundary collision",
			player: Player{
				Point: Point{
					X: 95,
					Y: 95,
				},
				Speed: 10,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  Keys{"ArrowRight": true, "ArrowDown": true},
			},
			world:     World{Width: 100, Height: 100},
			obstacles: []Rect{},
			want: Point{
				X: 90, // 100 (world width) - 10 (player width)
				Y: 90, // 100 (world height) - 10 (player height)
			},
		},
		{
			name: "diagonal movement",
			player: Player{
				Point: Point{
					X: 50,
					Y: 50,
				},
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  Keys{"ArrowRight": true, "ArrowDown": true},
			},
			world:     World{Width: 100, Height: 100},
			obstacles: []Rect{},
			want: Point{
				X: 55,
				Y: 55,
			},
		},
		{
			name: "no movement when no keys pressed",
			player: Player{
				Point: Point{
					X: 50,
					Y: 50,
				},
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  Keys{},
			},
			world:     World{Width: 100, Height: 100},
			obstacles: []Rect{},
			want: Point{
				X: 50,
				Y: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.Update(&tt.world, tt.obstacles)

			if tt.player.Point != tt.want {
				t.Errorf("Player.Point = %v, want %v", tt.player.Point, tt.want)
			}
		})
	}
}

func TestPlayerHandleMovement(t *testing.T) {
	tests := []struct {
		Keys
		name       string
		wantMoving bool
	}{
		{
			name:       "no movement",
			Keys:       Keys{},
			wantMoving: false,
		},
		{
			name:       "moving right",
			Keys:       Keys{"ArrowRight": true},
			wantMoving: true,
		},
		{
			name:       "moving diagonally",
			Keys:       Keys{"ArrowRight": true, "ArrowUp": true},
			wantMoving: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				Keys:  tt.Keys,
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
			}
			moving := p.handleMovement()
			if moving != tt.wantMoving {
				t.Errorf("handleMovement() = %v, want %v", moving, tt.wantMoving)
			}
		})
	}
}

func TestPlayerClampToWorldBounds(t *testing.T) {
	tests := []struct {
		name   string
		player Player
		world  World
		want   Point
	}{
		{
			name: "within bounds",
			player: Player{
				Point: Point{
					X: 50,
					Y: 50,
				},
				Rect: Rect{Width: 10, Height: 10},
			},
			world: World{Width: 100, Height: 100},
			want: Point{
				X: 50,
				Y: 50,
			},
		},
		{
			name: "outside right bound",
			player: Player{
				Point: Point{
					X: 95,
					Y: 50,
				},
				Rect: Rect{Width: 10, Height: 10},
			},
			world: World{Width: 100, Height: 100},
			want: Point{
				X: 90, // 100 - 10
				Y: 50,
			},
		},
		{
			name: "outside bottom bound",
			player: Player{
				Point: Point{
					X: 50,
					Y: 95,
				},
				Rect: Rect{Width: 10, Height: 10},
			},
			world: World{Width: 100, Height: 100},
			want: Point{
				X: 50,
				Y: 90, // 100 - 10
			},
		},
		{
			name: "outside negative bounds",
			player: Player{
				Point: Point{
					X: -5,
					Y: -5,
				},
				Rect: Rect{Width: 10, Height: 10},
			},
			world: World{Width: 100, Height: 100},
			want: Point{
				X: 0,
				Y: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.clampToWorldBounds(&tt.world)
			if tt.player.Point != tt.want {
				t.Errorf("Player.X = %v, want %v", tt.player.Point, tt.want)
			}
		})
	}
}
