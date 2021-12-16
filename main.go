package main

import (
	"fmt"
	_ "image/png" //image functions import

	//"image/color" //needed for color fill
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" //needed for Hello World screen print
)

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
	// Display the information with "X: xx, Y: xx" format
	x, y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))
	npc := &ebiten.DrawImageOptions{}
	npc.ColorM.ChangeHSV(2.85, 2.00, 1.00)
	npc.GeoM.Scale(-1, 1)
	npc.GeoM.Translate(615, 30)

	char := &ebiten.DrawImageOptions{}
	char.GeoM.Translate(25, 30)

	screen.DrawImage(img, char)
	screen.DrawImage(img, npc) //Draw same Image with declareed changes 2nd argument (scaled 1.5 translated 50,50 = op)
	//screen.DrawImage(img, nil) //Draw Default Image on Screen every tick
	//screen.Fill(color.RGBA{0xff, 0, 0, 0xff}) //fills screen color Red
	//ebitenutil.DebugPrint(screen, "Hello, World!") //Draws Hello World Text
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480 //sets window size
	//return 320, 240 // sets window size
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Geometry Matrix") //sets Window Title to "Geometry Matrix"
	//ebiten.SetWindowTitle("Render an image") //sets Window Title to "Render an image"
	//ebiten.SetWindowTitle("Fill") //sets Window Title to Fill
	//ebiten.SetWindowTitle("Hello, World!") //sets Window Title to Hello World!
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
