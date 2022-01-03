package main

import (
	//"bytes"
	//_ "embed"
	"fmt"
	//"image"
	"image/color"
	_ "image/png" //image functions import
	"log"
	"math/rand"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" //needed for Hello World screen print
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var windowWidth, windowHeight int = 1280, 960

var spell1 = Action{"Basic Attack", 2, [2]int{0, 0}, "damage"}
var spell2 = Action{"Heavy Attack", 4, [2]int{1, 0}, "damage"}
var spell3 = Action{"Drink Potion", 3, [2]int{2, 0}, "heal"}
var spell4 = Action{"Attack Up", 1, [2]int{2, 0}, "attackBuff"}
var spell5 = Action{"Bite", 1, [2]int{0, 0}, "damage"}
var spell6 = Action{"Scratch", 2, [2]int{1, 0}, "damage"}
var spell7 = Action{"Enrage", 1, [2]int{2, 0}, "attackBuff"}
var spell8 = Action{"Block", 1, [2]int{2, 0}, "defenseBuff"}

var DurationOfTime = time.Duration(3) * time.Second

var player Entity
var enemy Entity

var turnText string = "Battle Starts!"
var turn int = 0

var healthBarGreenE *ebiten.Image
var healthBarGreenP *ebiten.Image
var healthBarRed *ebiten.Image

var playerTurn bool = true
var playerDead bool = false
var enemyDead bool = false

var box *ebiten.Image
var arrow *ebiten.Image
var background *ebiten.Image

//var img *ebiten.Image //variable declared for pointer image
var hero *ebiten.Image
var bandit *ebiten.Image

/*const (
	iconSize = 32
	tileXNum = 25
)

var (
	iconsImage *ebiten.Image
)*/

var arrowPos = [2]int{0, 0}
var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() { //init function grabbing image from directory
	var err error
	/*img, _, err = ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}*/

	/*img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	iconsImage = ebiten.NewImageFromImage(img)*/

	background, _, err = ebitenutil.NewImageFromFile("background.png")
	if err != nil {
		log.Fatal(err)
	}

	hero, _, err = ebitenutil.NewImageFromFile("hero.png")
	if err != nil {
		log.Fatal(err)
	}

	bandit, _, err = ebitenutil.NewImageFromFile("bandit.png")
	if err != nil {
		log.Fatal(err)
	}

	player.Name = "Hero Knight"
	player.Actions = [4]Action{spell1, spell2, spell3, spell4}
	player.Stats.maxHealth = 25
	player.Stats.currentHealth = 25
	player.Stats.attack = 0
	player.Stats.defense = 0
	player.Image = *hero
	var x, y = player.Image.Size()
	player.Size = [2]int{x, y}
	player.Position = [2]int{windowWidth * 1 / 20, windowHeight * 4 / 10}

	enemy.Name = "Red Gopher"
	enemy.Actions = [4]Action{spell5, spell6, spell7, spell8}
	enemy.Stats.maxHealth = 25
	enemy.Stats.currentHealth = 25
	enemy.Stats.attack = 0
	enemy.Stats.defense = 0
	enemy.Image = *bandit
	enemy.Size = [2]int{x, y}
	enemy.Position = [2]int{windowWidth*16/20 - enemy.Size[0], windowHeight * 4 / 10}

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
	CoolDown       [2]int
	Effect         string
}

type Entity struct {
	Name     string
	Position [2]int
	Size     [2]int
	Stats    struct {
		maxHealth     int
		currentHealth int
		attack        int
		defense       int
	}
	Actions [4]Action
	Image   ebiten.Image
}

func ActionEffects(first Entity, second Entity, spell Action) (struct {
	maxHealth     int
	currentHealth int
	attack        int
	defense       int
}, struct {
	maxHealth     int
	currentHealth int
	attack        int
	defense       int
}, [4]Action,
	bool) {
	for i := 0; i < len(first.Actions); i++ {
		if first.Actions[i].CoolDown[1] < turn {
			if first.Actions[i] == spell {
				if spell.Effect == "attackBuff" {
					first.Stats.attack = 0
				} else if spell.Effect == "defenseBuff" {
					first.Stats.defense = 0
				}
				first.Actions[i].CoolDown[1] += spell.CoolDown[0]
			}
			first.Actions[i].CoolDown[1] += 1
		}
		if first.Actions[i].CoolDown[1] == turn {
			if first.Actions[i].Effect == "attackBuff" && first.Stats.attack >= 1 {
				first.Stats.attack = 0
				fmt.Println("attack buff wore off")
				fmt.Println("player attack", first.Stats.attack)
			} else if first.Actions[i].Effect == "defenseBuff" && first.Stats.defense > 0 {
				first.Stats.defense = 0
				fmt.Println("defense buff wore off")
				fmt.Println("player defense", first.Stats.defense)
			}
		}
	}
	var dead bool = false
	if spell.Effect == "damage" {
		if (second.Stats.currentHealth - spell.HealthModifier + first.Stats.attack - second.Stats.defense) >= 1 {
			second.Stats.currentHealth -= (spell.HealthModifier + first.Stats.attack - second.Stats.defense)
		} else {
			second.Stats.currentHealth = 1
			dead = true
		}
		//fmt.Println(first.Name)
		//fmt.Println("used Damage spell")
	} else if spell.Effect == "heal" {
		first.Stats.currentHealth += spell.HealthModifier
		if first.Stats.currentHealth > first.Stats.maxHealth {
			first.Stats.currentHealth = first.Stats.maxHealth
		}
		//fmt.Println(first.Name)
		//fmt.Println("used heal spell")
	} else if spell.Effect == "attackBuff" {
		first.Stats.attack += spell.HealthModifier
		//fmt.Println(first.Name)
		//fmt.Println("used attack buff spell")
	} else if spell.Effect == "defenseBuff" {
		first.Stats.defense += spell.HealthModifier
		//fmt.Println(first.Name)
		//fmt.Println("used defense buff spell")
	}

	return first.Stats, second.Stats, first.Actions, dead
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && playerTurn {
		if arrowPos == [2]int{0, 0} {
			fmt.Println("current turn:", turn)
			fmt.Println("current cd:", player.Actions[0].CoolDown[1])

			if player.Actions[0].CoolDown[1] == turn {
				playerTurn = false
				turn += 1
				player.Stats, enemy.Stats, player.Actions, enemyDead = ActionEffects(player, enemy, player.Actions[0])
				turnText = "Enemy's Turn"
				time.AfterFunc(DurationOfTime, func() {
					var x int
					for {
						x = rand.Intn(4)
						if (enemy.Actions[x].CoolDown[1] + 1) == turn {
							break
						}
					}

					enemy.Stats, player.Stats, enemy.Actions, playerDead = ActionEffects(enemy, player, enemy.Actions[x])
					turnText = "Your Turn"
					playerTurn = true
					//fmt.Println(player.Stats)
					//fmt.Println(enemy.Stats)
				})
			} else {
				turnText = "That spell is on COOLDOWN!"
			}
		} else if arrowPos == [2]int{11, 0} {
			fmt.Println("current turn:", turn)
			fmt.Println("current cd:", player.Actions[1].CoolDown[1])
			if player.Actions[1].CoolDown[1] == turn {
				playerTurn = false
				turn += 1
				player.Stats, enemy.Stats, player.Actions, enemyDead = ActionEffects(player, enemy, player.Actions[1])
				turnText = "Enemy's Turn"
				time.AfterFunc(DurationOfTime, func() {
					x := rand.Intn(4)
					enemy.Stats, player.Stats, enemy.Actions, playerDead = ActionEffects(enemy, player, enemy.Actions[x])
					turnText = "Your Turn"
					playerTurn = true
					//fmt.Println(player.Stats)
					//fmt.Println(enemy.Stats)
				})
			} else {
				turnText = "That spell is on COOLDOWN!"
			}
		} else if arrowPos == [2]int{0, 12} {
			fmt.Println("current turn:", turn)
			fmt.Println("current cd:", player.Actions[2].CoolDown[1])
			if player.Actions[2].CoolDown[1] == turn {
				playerTurn = false
				turn += 1
				player.Stats, enemy.Stats, player.Actions, enemyDead = ActionEffects(player, enemy, player.Actions[2])
				turnText = "Enemy's Turn"
				time.AfterFunc(DurationOfTime, func() {
					x := rand.Intn(4)
					enemy.Stats, player.Stats, enemy.Actions, playerDead = ActionEffects(enemy, player, enemy.Actions[x])
					turnText = "Your Turn"
					playerTurn = true
					//fmt.Println(player.Stats)
					//fmt.Println(enemy.Stats)
				})
			} else {
				turnText = "That spell is on COOLDOWN!"
			}
		} else if arrowPos == [2]int{11, 12} {
			fmt.Println("current turn:", turn)
			fmt.Println("current cd:", player.Actions[3].CoolDown[1])
			if player.Actions[3].CoolDown[1] == turn {
				playerTurn = false
				turn += 1
				player.Stats, enemy.Stats, player.Actions, enemyDead = ActionEffects(player, enemy, player.Actions[3])
				turnText = "Enemy's Turn"
				time.AfterFunc(DurationOfTime, func() {
					x := rand.Intn(4)
					enemy.Stats, player.Stats, enemy.Actions, playerDead = ActionEffects(enemy, player, enemy.Actions[x])
					turnText = "Your Turn"
					playerTurn = true
					//fmt.Println(player.Stats)
					//fmt.Println(enemy.Stats)
				})
			} else {
				turnText = "That spell is on COOLDOWN!"
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		player.Stats.currentHealth -= 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Capture Mouse Position in variables
	x, y := ebiten.CursorPosition()
	// Display Mouse Position
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	npc := &ebiten.DrawImageOptions{}
	npc.GeoM.Scale(4.5, 4.5)
	npc.GeoM.Translate(float64(enemy.Position[0]), float64(enemy.Position[1]))

	//fmt.Println(npc.GeoM)

	char := &ebiten.DrawImageOptions{}
	char.GeoM.Scale(4, 4)
	char.GeoM.Translate(float64(player.Position[0]), float64(player.Position[1]))

	b := &ebiten.DrawImageOptions{}
	b.GeoM.Translate(float64(windowWidth)*4/20, float64(windowHeight)*25/40)

	bg := &ebiten.DrawImageOptions{}
	bg.GeoM.Scale(1.3, 1.3)

	box = ebiten.NewImage(windowWidth*47/80, windowHeight*25/80)
	box.Fill(color.RGBA{0x00, 0x00, 0x00, 0x7f})

	arr := &ebiten.DrawImageOptions{}
	arr.GeoM.Translate(float64(windowWidth*(11+arrowPos[0])/40), float64(windowHeight*(56+arrowPos[1])/80))
	arrow = ebiten.NewImage(12, 12)
	arrow.Fill(color.White)

	hpEnemy := &ebiten.DrawImageOptions{}
	hpEnemy.GeoM.Translate(float64(enemy.Position[0]-enemy.Size[0]+40), float64(enemy.Position[1]-40))
	hpPlayer := &ebiten.DrawImageOptions{}
	hpPlayer.GeoM.Translate(float64(player.Position[0]+40), float64(player.Position[1]-40))
	healthBarGreenP = ebiten.NewImage(300*player.Stats.currentHealth/player.Stats.maxHealth, 15)
	healthBarGreenP.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
	//fmt.Println(player.Health)
	healthBarGreenE = ebiten.NewImage(300*enemy.Stats.currentHealth/enemy.Stats.maxHealth, 15)
	healthBarGreenE.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})

	healthBarRed = ebiten.NewImage(300, 15)
	healthBarRed.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})

	screen.DrawImage(background, bg)

	screen.DrawImage(healthBarRed, hpEnemy)
	if !enemyDead {
		screen.DrawImage(healthBarGreenE, hpEnemy)
	}
	screen.DrawImage(healthBarRed, hpPlayer)
	if !playerDead {
		screen.DrawImage(healthBarGreenP, hpPlayer)
	}
	screen.DrawImage(box, b)
	screen.DrawImage(arrow, arr)
	screen.DrawImage(hero, char)
	screen.DrawImage(bandit, npc)

	text.Draw(screen, player.Actions[0].Name, mplusNormalFont, windowWidth*12/40, windowHeight*57/80, color.White)
	text.Draw(screen, player.Actions[1].Name, mplusNormalFont, windowWidth*(12+11)/40, windowHeight*57/80, color.White)
	text.Draw(screen, player.Actions[2].Name, mplusNormalFont, windowWidth*12/40, windowHeight*(57+12)/80, color.White)
	text.Draw(screen, player.Actions[3].Name, mplusNormalFont, windowWidth*(12+11)/40, windowHeight*(57+12)/80, color.White)

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
