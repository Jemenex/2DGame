package main

import (
	"fmt"
	"image/color"
	_ "image/png" //image functions import

	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" //needed for Hello World screen print
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var windowWidth, windowHeight int = 1280, 960

var spell1 = Action{"Basic Attack", 2, 0, "damage"}
var spell2 = Action{"Heavy Attack", 4, 2, "damage"}
var spell3 = Action{"Drink Potion", 6, 4, "heal"}
var spell4 = Action{"Defense Up", 1, 4, "buff"}
var spell5 = Action{"Bite", 1, 0, "damage"}
var spell6 = Action{"Scratch", 2, 2, "damage"}
var spell7 = Action{"Enrage", 1, 4, "buff"}
var spell8 = Action{"Block", 1, 3, "buff"}

var player Entity
var enemy Entity

var turnText string = "Battle Starts!"

var healthBarGreenE *ebiten.Image
var healthBarGreenP *ebiten.Image
var healthBarRed *ebiten.Image

var PlayerDead bool = false
var EnemyDead bool = false

var box *ebiten.Image
var arrow *ebiten.Image
var img *ebiten.Image //variable declared for pointer image

var arrowPos = [2]int{11, 12}
var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() { //init function grabbing image from directory
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	player.Name = "Blue Gopher"
	player.Actions = [4]Action{spell1, spell2, spell3, spell4}
	player.Health = [2]int{25, 25}
	player.Image = *img
	var x, y = player.Image.Size()
	player.Size = [2]int{x, y}
	player.Position = [2]int{windowWidth * 1 / 20, windowHeight * 2 / 10}

	enemy.Name = "Red Gopher"
	enemy.Actions = [4]Action{spell5, spell6, spell7, spell8}
	enemy.Health = [2]int{10, 10}
	enemy.Image = *img
	enemy.Size = [2]int{x, y}
	enemy.Position = [2]int{windowWidth * 19 / 20, windowHeight * 2 / 10}

	/*
		fmt.Println("-Player Position-")
		fmt.Println(player.Position)
		fmt.Println("-Player Size-")
		fmt.Println(player.Size)
		fmt.Println("-Enemy Position-")
		fmt.Println(enemy.Position)
		fmt.Println("-Enemy Size-")
		fmt.Println(enemy.Size)
		fmt.Println(img.Bounds())*/

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Action struct {
	Name           string
	HealthModifier int //How much it modifies health + or -
	CoolDown       int
	Effect         string
}

type Entity struct {
	Name     string
	Position [2]int
	Size     [2]int
	Health   [2]int //{Current Health, Total Health}
	Actions  [4]Action
	Image    ebiten.Image
}

type Game struct {
}

func (g *Game) Update() error {
	/*if inpututil.IsKeyJustPressed(ebiten.KeyUp) && player.Health[0] > 1 {
		player.Health[0] -= 1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) && player.Health[0] == 1 {
		PlayerDead = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && enemy.Health[0] > 1 {
		enemy.Health[0] -= 1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) && enemy.Health[0] == 1 {
		EnemyDead = true
	}*/
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		arrowPos[1] = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		arrowPos[1] = 12
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		arrowPos[0] = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		arrowPos[0] = 11
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if arrowPos == [2]int{0, 0} {
			enemy.Health[0] -= player.Actions[0].HealthModifier
		} else if arrowPos == [2]int{11, 0} {
			enemy.Health[0] -= player.Actions[1].HealthModifier
		} else if arrowPos == [2]int{0, 12} {
			enemy.Health[0] -= player.Actions[2].HealthModifier
		} else if arrowPos == [2]int{11, 12} {
			enemy.Health[0] -= player.Actions[3].HealthModifier
		}
	}
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
	b.GeoM.Translate(float64(windowWidth)*4/20, float64(windowHeight)*25/40)

	box = ebiten.NewImage(windowWidth*47/80, windowHeight*25/80)
	box.Fill(color.RGBA{0xb0, 0xb0, 0xb0, 0x0f})

	arr := &ebiten.DrawImageOptions{}
	arr.GeoM.Translate(float64(windowWidth*(11+arrowPos[0])/40), float64(windowHeight*(56+arrowPos[1])/80))
	arrow = ebiten.NewImage(12, 12)
	arrow.Fill(color.White)

	hpEnemy := &ebiten.DrawImageOptions{}
	hpEnemy.GeoM.Translate(float64(enemy.Position[0]-enemy.Size[0]*3/2+30), float64(enemy.Position[1]-40))
	hpPlayer := &ebiten.DrawImageOptions{}
	hpPlayer.GeoM.Translate(float64(player.Position[0]+30), float64(player.Position[1]-40))
	healthBarGreenP = ebiten.NewImage(300*player.Health[0]/player.Health[1], 30)
	healthBarGreenP.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
	//fmt.Println(player.Health)
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
	screen.DrawImage(arrow, arr)
	screen.DrawImage(img, char)
	screen.DrawImage(img, npc)

	text.Draw(screen, "first spell", mplusNormalFont, windowWidth*12/40, windowHeight*57/80, color.White)
	text.Draw(screen, "second spell", mplusNormalFont, windowWidth*(12+11)/40, windowHeight*57/80, color.White)
	text.Draw(screen, "third spell", mplusNormalFont, windowWidth*12/40, windowHeight*(57+12)/80, color.White)
	text.Draw(screen, "fourth spell", mplusNormalFont, windowWidth*(12+11)/40, windowHeight*(57+12)/80, color.White)

	text.Draw(screen, turnText, mplusBigFont, windowWidth*15/40, windowHeight*1/20, color.White)

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

/*
SPELLS OPTIONS AND SELECTION
	DRAW TEXT -DONE
	POSITION TEXT - DONE
	DRAW SELECTOR - DONE
	ATTACH VARIABLE WITH VALUE FOR POSITION  0,0  0,1 1,0 1,1 -DONE
	ATTACH VARIABLE CHANGE TO KEY PRESS DIRECTIONALLY - DONE
ENEMY TURN
BATTLE LOOP
*/
