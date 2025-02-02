package main

import (
	"testing"
)

func TestPlayerUpdate(t *testing.T) {
	tests := []struct {
		keys      map[string]bool
		obstacles []Rect
		name      string
		player    Player
		world     World
		wantX     float64
		wantY     float64
	}{
		{
			name: "move right without obstacles",
			player: Player{
				X:     0,
				Y:     0,
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  map[string]bool{"ArrowRight": true},
			},
			world:     World{Width: 100, Height: 100},
			obstacles: []Rect{},
			wantX:     5,
			wantY:     0,
		},
		{
			name: "blocked by obstacle",
			player: Player{
				X:     0,
				Y:     0,
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  map[string]bool{"ArrowRight": true},
			},
			world: World{Width: 100, Height: 100},
			obstacles: []Rect{
				{X: 5, Y: 0, Width: 10, Height: 10},
			},
			wantX: 0,
			wantY: 0,
		},
		{
			name: "world boundary collision",
			player: Player{
				X:     95,
				Y:     95,
				Speed: 10,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  map[string]bool{"ArrowRight": true, "ArrowDown": true},
			},
			world:     World{Width: 100, Height: 100},
			obstacles: []Rect{},
			wantX:     90, // 100 (world width) - 10 (player width)
			wantY:     90, // 100 (world height) - 10 (player height)
		},
		{
			name: "diagonal movement",
			player: Player{
				X:     50,
				Y:     50,
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  map[string]bool{"ArrowRight": true, "ArrowDown": true},
			},
			world:     World{Width: 100, Height: 100},
			obstacles: []Rect{},
			wantX:     55,
			wantY:     55,
		},
		{
			name: "no movement when no keys pressed",
			player: Player{
				X:     50,
				Y:     50,
				Speed: 5,
				Rect:  Rect{Width: 10, Height: 10},
				Keys:  map[string]bool{},
			},
			world:     World{Width: 100, Height: 100},
			obstacles: []Rect{},
			wantX:     50,
			wantY:     50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.Update(&tt.world, tt.obstacles)

			if tt.player.X != tt.wantX {
				t.Errorf("Player.X = %v, want %v", tt.player.X, tt.wantX)
			}
			if tt.player.Y != tt.wantY {
				t.Errorf("Player.Y = %v, want %v", tt.player.Y, tt.wantY)
			}
		})
	}
}

func TestPlayerHandleMovement(t *testing.T) {
	tests := []struct {
		keys       map[string]bool
		name       string
		wantMoving bool
	}{
		{
			name:       "no movement",
			keys:       map[string]bool{},
			wantMoving: false,
		},
		{
			name:       "moving right",
			keys:       map[string]bool{"ArrowRight": true},
			wantMoving: true,
		},
		{
			name:       "moving diagonally",
			keys:       map[string]bool{"ArrowRight": true, "ArrowUp": true},
			wantMoving: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				Keys:  tt.keys,
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
		wantX  float64
		wantY  float64
	}{
		{
			name: "within bounds",
			player: Player{
				X:    50,
				Y:    50,
				Rect: Rect{Width: 10, Height: 10},
			},
			world: World{Width: 100, Height: 100},
			wantX: 50,
			wantY: 50,
		},
		{
			name: "outside right bound",
			player: Player{
				X:    95,
				Y:    50,
				Rect: Rect{Width: 10, Height: 10},
			},
			world: World{Width: 100, Height: 100},
			wantX: 90, // 100 - 10
			wantY: 50,
		},
		{
			name: "outside bottom bound",
			player: Player{
				X:    50,
				Y:    95,
				Rect: Rect{Width: 10, Height: 10},
			},
			world: World{Width: 100, Height: 100},
			wantX: 50,
			wantY: 90, // 100 - 10
		},
		{
			name: "outside negative bounds",
			player: Player{
				X:    -5,
				Y:    -5,
				Rect: Rect{Width: 10, Height: 10},
			},
			world: World{Width: 100, Height: 100},
			wantX: 0,
			wantY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.clampToWorldBounds(&tt.world)
			if tt.player.X != tt.wantX {
				t.Errorf("Player.X = %v, want %v", tt.player.X, tt.wantX)
			}
			if tt.player.Y != tt.wantY {
				t.Errorf("Player.Y = %v, want %v", tt.player.Y, tt.wantY)
			}
		})
	}
}
