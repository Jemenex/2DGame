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

var healthBarGreenE *ebiten.Image
var healthBarGreenP *ebiten.Image
var healthBarRed *ebiten.Image

var PlayerDead bool = false
var EnemyDead bool = false

var box *ebiten.Image
var img *ebiten.Image //variable declared for pointer image

func init() { //init function grabbing image from directory
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")

	player.Name = "Blue Gopher"
	player.Actions = [2]Action{spell1, spell2}
	player.Health = [2]int{25, 25}
	player.Image = *img
	var x, y = player.Image.Size()
	player.Size = [2]int{x, y}
	player.Position = [2]int{windowWidth * 1 / 20, windowHeight * 2 / 10}

	enemy.Name = "Red Gopher"
	enemy.Actions = [2]Action{spell3, spell4}
	enemy.Health = [2]int{10, 10}
	enemy.Image = *img
	enemy.Size = [2]int{x, y}
	enemy.Position = [2]int{windowWidth * 19 / 20, windowHeight * 2 / 10}

	fmt.Println("-Player Position-")
	fmt.Println(player.Position)
	fmt.Println("-Player Size-")
	fmt.Println(player.Size)
	fmt.Println("-Enemy Position-")
	fmt.Println(enemy.Position)
	fmt.Println("-Enemy Size-")
	fmt.Println(enemy.Size)
	fmt.Println(img.Bounds())

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
	Health   [2]int
	Actions  [2]Action
	Image    ebiten.Image
}

type Game struct {
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) && player.Health[0] > 1 {
		player.Health[0] -= 1
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) && player.Health[0] == 1 {
		PlayerDead = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && enemy.Health[0] > 1 {
		enemy.Health[0] -= 1
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && enemy.Health[0] == 1 {
		EnemyDead = true
	}
	//fmt.Println(player.Health)
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

	hpEnemy := &ebiten.DrawImageOptions{}
	hpEnemy.GeoM.Translate(float64(enemy.Position[0]-enemy.Size[0]*3/2+30), float64(enemy.Position[1]-40))
	hpPlayer := &ebiten.DrawImageOptions{}
	hpPlayer.GeoM.Translate(float64(player.Position[0]+30), float64(player.Position[1]-40))
	healthBarGreenP = ebiten.NewImage(300*player.Health[0]/player.Health[1], 30)
	healthBarGreenP.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
	fmt.Println(player.Health)
	healthBarGreenE = ebiten.NewImage(300*enemy.Health[0]/enemy.Health[1], 30)
	healthBarGreenE.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})

	healthBarRed = ebiten.NewImage(300, 30)
	healthBarRed.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})

	screen.DrawImage(healthBarRed, hpEnemy)
	if !EnemyDead {
		screen.DrawImage(healthBarGreenE, hpEnemy)
	}
	screen.DrawImage(healthBarRed, hpPlayer)
	if !PlayerDead {
		screen.DrawImage(healthBarGreenP, hpPlayer)
	}
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
