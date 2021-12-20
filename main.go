package main

import (
	"fmt"
	"image/color"
	_ "image/png" //image functions import

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" //needed for Hello World screen print
)

var windowWidth, windowHeight int = 1280, 960

var spell1 = Action{"Basic Attack", 2, 0}
var spell2 = Action{"Heavy Attack", 4, 2}
var spell3 = Action{"Bite", 1, 0}
var spell4 = Action{"Scratch", 2, 2}

var player Entity
var enemy Entity

var healthBar *ebiten.Image
var box *ebiten.Image
var img *ebiten.Image //variable declared for pointer image

func init() { //init function grabbing image from directory
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")

	player.Name = "Blue Gopher"
	player.Actions = [2]Action{spell1, spell2}
	player.Health = 25
	player.Image = *img
	var x, y = player.Image.Size()
	player.Size = [2]int{x, y}
	player.Position = [2]int{50, 200}

	enemy.Name = "Red Gopher"
	enemy.Actions = [2]Action{spell3, spell4}
	enemy.Health = 10
	enemy.Image = *img
	enemy.Size = [2]int{x, y}
	enemy.Position = [2]int{1230, 200}

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
	npc.GeoM.Scale(-1.5, 1.5)

	npc.GeoM.Translate(float64(enemy.Position[0]), float64(enemy.Position[1]))

	//fmt.Println(npc.GeoM)

	char := &ebiten.DrawImageOptions{}
	char.GeoM.Scale(1.5, 1.5)
	char.GeoM.Translate(float64(player.Position[0]), float64(player.Position[1]))

	b := &ebiten.DrawImageOptions{}
	b.GeoM.Translate(float64(windowWidth)*2/20, float64(windowHeight)*15/20)
	box = ebiten.NewImage(windowWidth*16/20, windowHeight*3/20)
	box.Fill(color.RGBA{0xb0, 0xb0, 0xb0, 0x2f})
	//box = ebiten.NewImage(50, 50)

	hp := &ebiten.DrawImageOptions{}
	hp.GeoM.Translate(float64(enemy.Position[0]), float64(enemy.Position[1]+100))
	//hp.GeoM.Translate(float64(enemy.Position[0]+enemy.Size[0]/2), float64(enemy.Position[1]+enemy.Size[1]+100))
	healthBar = ebiten.NewImage(200, 40)
	healthBar.Fill(color.RGBA{0xb0, 0xb0, 0xb0, 0x4f})

	screen.DrawImage(healthBar, hp)
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
