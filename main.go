package main

import (
	"fmt"
	"image/color"
	_ "image/png" //image functions import


	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" //needed for Hello World screen print
)


var windowWidth, windowHeight int = 1000, 800

var box *ebiten.Image

var img *ebiten.Image //variable declared for pointer image

func init() { //init function grabbing image from directory
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")

	spell1 := Action{"Basic Attack", 2, 0}
	spell2 := Action{"Heavy Attack", 4, 2}

	player := new(Entity)
	player.Name = "Blue Gopher"
	player.Actions = [2]Action{spell1, spell2}
	player.Health = 100
	player.Image = *img
	//player.Size = [2]int player.Image.Size()

	fmt.Println(player)
	fmt.Println(player.Image.Size())

	enemy := new(Entity)
	enemy.Name = "Red Gopher"

	if err != nil {
		log.Fatal(err)
	}
}

type Action struct {
	Name     string
	Damage   int
	CoolDown int
}

type Entity struct {
	Name     string
	Position [2]int
	Size     [2]int
	Health   int
	Actions  [2]Action
	Image    ebiten.Image
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

	npc.GeoM.Translate(950, 50)

	char := &ebiten.DrawImageOptions{}
	char.GeoM.Translate(50, 50)


	b := &ebiten.DrawImageOptions{}
	b.GeoM.Translate(float64(windowWidth)*1/20, float64(windowHeight)/2)
	box = ebiten.NewImage(windowWidth*9/10, windowHeight*9/20)
	//box = ebiten.NewImage(50, 50)
	box.Fill(color.RGBA{0xb0, 0xb0, 0xb0, 0x2f})

	screen.DrawImage(box, b)
	screen.DrawImage(img, char)
	screen.DrawImage(img, npc)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight //sets window size
}

func main() {

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("2D Game")


	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
