package main

import (
	"fmt"
	_ "image/png" //image functions import

	"image/color" //needed for color fill
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" //needed for Hello World screen print
)

//Vars
var windowWidth = 1000
var windowHeight = 800

//Images
var divider *ebiten.Image
var img *ebiten.Image //variable declared for pointer image

func init() { //init function grabbing image from directory
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Capture Mouse Position in variables
	x, y := ebiten.CursorPosition()
	// Display Mouse Position
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	npc := &ebiten.DrawImageOptions{}
	npc.ColorM.ChangeHSV(2.85, 2.00, 1.00)
	npc.GeoM.Scale(-1, 1)
	npc.GeoM.Translate(950, 80)

	char := &ebiten.DrawImageOptions{}
	char.GeoM.Translate(50, 80)

	div := &ebiten.DrawImageOptions{}
	div.GeoM.Translate(float64(windowWidth)/20.00, float64(windowHeight)/2.00)

	screen.DrawImage(img, char)
	screen.DrawImage(img, npc)
	// Create an bottom box
	divider = ebiten.NewImage(int(float64(windowWidth)*.9), int(float64(windowHeight)*.45))
	divider.Fill(color.RGBA{0xb0, 0xb0, 0xb0, 0x5f})
	screen.DrawImage(divider, div)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("2D Game")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
