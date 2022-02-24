package main

import (
	_ "embed"
	_ "image/png"
	"reflect"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestActionEffects(t *testing.T) {
	type args struct {
		first  Entity
		second Entity
		spell  Action
	}
	tests := []struct {
		name string
		args args
		want struct {
			maxHealth     int
			currentHealth int
			attack        int
			defense       int
		}
		want1 struct {
			maxHealth     int
			currentHealth int
			attack        int
			defense       int
		}
		want2 [4]Action
		want3 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := ActionEffects(tt.args.first, tt.args.second, tt.args.spell)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ActionEffects() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ActionEffects() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("ActionEffects() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("ActionEffects() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func TestGame_Update(t *testing.T) {
	tests := []struct {
		name    string
		g       *Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Update(); (err != nil) != tt.wantErr {
				t.Errorf("Game.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_Draw(t *testing.T) {
	type args struct {
		screen *ebiten.Image
	}
	tests := []struct {
		name string
		g    *Game
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.Draw(tt.args.screen)
		})
	}
}

func TestGame_Layout(t *testing.T) {
	type args struct {
		outsideWidth  int
		outsideHeight int
	}
	tests := []struct {
		name             string
		g                *Game
		args             args
		wantScreenWidth  int
		wantScreenHeight int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScreenWidth, gotScreenHeight := tt.g.Layout(tt.args.outsideWidth, tt.args.outsideHeight)
			if gotScreenWidth != tt.wantScreenWidth {
				t.Errorf("Game.Layout() gotScreenWidth = %v, want %v", gotScreenWidth, tt.wantScreenWidth)
			}
			if gotScreenHeight != tt.wantScreenHeight {
				t.Errorf("Game.Layout() gotScreenHeight = %v, want %v", gotScreenHeight, tt.wantScreenHeight)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
