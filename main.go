package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png" //image functions import
	"log"
	"math/rand"

	"strconv"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" //needed for Hello World screen print
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const windowWidth, windowHeight int = 1280, 960

var spell1 = Action{"Basic Attack", 2, [2]int{0, 0}, "damage"}
var spell2 = Action{"Heavy Attack", 4, [2]int{1, 0}, "damage"}
var spell3 = Action{"Drink Potion", 3, [2]int{2, 0}, "heal"}
var spell4 = Action{"Attack Up", 1, [2]int{2, 0}, "attackBuff"}
var spell5 = Action{"Bite", 1, [2]int{0, 0}, "damage"}
var spell6 = Action{"Scratch", 2, [2]int{1, 0}, "damage"}
var spell7 = Action{"Enrage", 1, [2]int{2, 0}, "attackBuff"}
var spell8 = Action{"Block", 1, [2]int{2, 0}, "defenseBuff"}

var DurationOfTime = time.Duration(3) * time.Second

var FadeTime = time.Duration(2) * time.Second

var player Entity
var enemy Entity

var turnText string = "Battle Starts!"
var turn int = 0
var gameState string = "Battle" // Menu - Battle - Uppgrade after battle

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
//go:embed HeroKnight.png
var heroKnightPng []byte
var hero *ebiten.Image

//go:embed HeavyBandit.png
var heavyBanditPng []byte
var bandit *ebiten.Image

const (
	iconSize = 32
	iconXNum = 16
)

// const (
// 	frameOX     = 100
// 	frameOY     = 0
// 	frameWidth  = 100
// 	frameHeight = 55
// 	frameNum    = 7
// )

//go:embed icons.png
var iconsPng []byte
var iconsImg *ebiten.Image

var arrowPos int = 0
var fadeAlpha uint8 = 255
var fading bool = false

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)
var startingPosP [2]int
var startingPosE [2]int
var currentAnimationP Anim
var currentAnimationE Anim
var playerIdlePos [2]int
var enemyIdlePos [2]int
var playerAttackPos [2]int
var enemyAttackPos [2]int
var playerDestination [2]int
var enemyDestination [2]int

func init() { //init function grabbing image from directory
	var err error

	img, _, err := image.Decode(bytes.NewReader(iconsPng))
	if err != nil {
		log.Fatal(err)
	}
	iconsImg = ebiten.NewImageFromImage(img)

	img2, _, err := image.Decode(bytes.NewReader(heroKnightPng))
	if err != nil {
		log.Fatal(err)
	}
	hero = ebiten.NewImageFromImage(img2)

	img3, _, err := image.Decode(bytes.NewReader(heavyBanditPng))
	if err != nil {
		log.Fatal(err)
	}
	bandit = ebiten.NewImageFromImage(img3)

	background, _, err = ebitenutil.NewImageFromFile("background.png")
	if err != nil {
		log.Fatal(err)
	}

	/*hero, _, err = ebitenutil.NewImageFromFile("hero.png")
	if err != nil {
		log.Fatal(err)
	}*/

	player.Name = "Hero Knight"
	player.Actions = [4]Action{spell1, spell2, spell3, spell4}
	player.Stats.maxHealth = 25
	player.Stats.currentHealth = 25
	player.Stats.attack = 0
	player.Stats.defense = 0
	player.Image = *hero
	player.Size = [2]int{100, 55}
	player.Position = [2]int{windowWidth * 1 / 20, windowHeight * 4 / 10}
	player.Anims.Idle = Anim{
		0,
		0,
		100,
		55,
		8,
		7,
		1000,
	}
	player.Anims.Attack = Anim{
		800,
		55,
		100,
		55,
		6,
		7,
		1000,
	}
	player.Anims.Attack2 = Anim{
		400,
		110,
		100,
		55,
		6,
		4,
		1000,
	}
	player.Anims.Attack3 = Anim{
		100,
		165,
		100,
		55,
		6,
		7,
		1000,
	}
	player.Anims.Run = Anim{
		800,
		0,
		100,
		55,
		10,
		7,
		1000,
	}
	player.Anims.Hurt = Anim{
		500,
		220,
		100,
		55,
		3,
		7,
		1000,
	}
	player.Anims.Death = Anim{
		800,
		220,
		100,
		55,
		8,
		7,
		1000,
	}
	//fmt.Println(player.Anims.Idle.frameHeight)

	enemy.Name = "Bandit"
	enemy.Actions = [4]Action{spell5, spell6, spell7, spell8}
	enemy.Stats.maxHealth = 25
	enemy.Stats.currentHealth = 25
	enemy.Stats.attack = 0
	enemy.Stats.defense = 0
	enemy.Image = *bandit
	enemy.Size = [2]int{48, 48}
	enemy.Position = [2]int{windowWidth*16/20 - enemy.Size[0]*2, windowHeight * 4 / 10}
	enemy.Anims.Idle = Anim{
		192,
		0,
		48,
		48,
		4,
		12,
		384,
	}
	enemy.Anims.Attack = Anim{
		0,
		96,
		48,
		48,
		8,
		12,
		384,
	}
	enemy.Anims.Hurt = Anim{
		0,
		192,
		48,
		48,
		3,
		12,
		384,
	}
	enemy.Anims.Hurt = Anim{
		0,
		192,
		48,
		48,
		3,
		12,
		384,
	}
	enemy.Anims.Hurt = Anim{
		0,
		192,
		48,
		48,
		3,
		12,
		384,
	}
	enemy.Anims.Death = Anim{
		192,
		192,
		48,
		48,
		1,
		12,
		384,
	}
	enemyIdlePos = enemy.Position
	enemyAttackPos = [2]int{player.Position[0] + enemy.Size[0]*3/2, player.Position[1]}
	enemyDestination = enemyIdlePos
	playerIdlePos = player.Position
	playerAttackPos = [2]int{enemy.Position[0] - player.Size[0]*3/2, enemy.Position[1]}
	playerDestination = playerIdlePos
	currentAnimationP = player.Anims.Idle
	currentAnimationE = enemy.Anims.Idle
	startingPosP = [2]int{windowWidth * 1 / 20, windowHeight * 4 / 10}
	startingPosE = [2]int{windowWidth*16/20 - enemy.Size[0]*2, windowHeight * 4 / 10}
	// fmt.Println(currentAnimationP)
	// fmt.Println(player.Anims.Idle.frameNum)
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

type Anim struct {
	frameOX     int
	frameOY     int
	frameWidth  int
	frameHeight int
	frameNum    int
	framesPer   int
	maxWidth    int
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
	Anims   struct {
		Idle    Anim
		Death   Anim
		Attack  Anim
		Attack2 Anim
		Attack3 Anim
		Run     Anim
		Hurt    Anim
	}
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
	count  int
	layers [][]int
}

func (g *Game) Update() error {
	g.count++
	if player.Position[0] != playerDestination[0] {
		if player.Position[0] < playerDestination[0] {
			player.Position[0]++
		}
		if player.Position[0] > playerDestination[0] {
			player.Position[0]--
		}
	}

	if enemy.Position[0] != enemyDestination[0] {
		if enemy.Position[0] > enemyDestination[0] {
			enemy.Position[0]++
		}
		if enemy.Position[0] < enemyDestination[0] {
			enemy.Position[0]--
		}
	}

	if gameState == "Battle" { //Battle State Controls
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			if arrowPos >= 64 {
				arrowPos -= 64
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			if arrowPos <= 384 {
				arrowPos += 64
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && playerTurn {
			if arrowPos == 0 {
				fmt.Println("current turn:", turn)
				fmt.Println("current cd:", player.Actions[0].CoolDown[1])

				if player.Actions[0].CoolDown[1] == turn {
					playerTurn = false
					turn += 1
					g.count = 0
					currentAnimationP = player.Anims.Attack
					playerDestination = playerAttackPos
					player.Stats, enemy.Stats, player.Actions, enemyDead = ActionEffects(player, enemy, player.Actions[0])
					turnText = "Enemy's Turn"
					time.AfterFunc(DurationOfTime, func() {
						currentAnimationP = player.Anims.Idle
						playerDestination = playerIdlePos
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
			} else if arrowPos == 64 {
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
			} else if arrowPos == 128 {
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
			} else if arrowPos == 192 {
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
			} else if arrowPos > 192 {
				turnText = "No spell assigned!"
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			player.Stats.currentHealth -= 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			fading = true
			fmt.Println(fading)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			currentAnimationP = player.Anims.Run
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			player.Position[0] += 10
		}
		if playerDead {
			turnText = "Battle is over " + player.Name + " has died"
		} else if enemyDead {
			turnText = "Battle is over " + enemy.Name + " has died"
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//Enemy Draw Options - Position . Scale
	npc := &ebiten.DrawImageOptions{}
	npc.GeoM.Scale(4.5, 4.5)
	npc.GeoM.Translate(float64(enemy.Position[0]), float64(enemy.Position[1]))
	//Player Draw Options - Position . Scale
	char := &ebiten.DrawImageOptions{}
	char.GeoM.Scale(4, 4)
	char.GeoM.Translate(float64(player.Position[0]), float64(player.Position[1]))
	//Player
	ap := &ebiten.DrawImageOptions{}
	ap.GeoM.Scale(4, 4)
	ap.GeoM.Translate(float64(player.Position[0]), float64(player.Position[1]))
	//Background Draw Options - Scale
	bg := &ebiten.DrawImageOptions{}
	bg.GeoM.Scale(1.3, 1.3)
	//Arrow/Highlight Draw Options - Position . Size . Color
	arr := &ebiten.DrawImageOptions{}
	arr.GeoM.Translate(float64(384+arrowPos), float64(832))
	arrow = ebiten.NewImage(64, 64)
	arrow.Fill(color.RGBA{0xf9, 0xe8, 0x2f, 0x2f})
	//Hp Bar Draw Options - Position . Size . Fill
	hpEnemy := &ebiten.DrawImageOptions{}
	hpEnemy.GeoM.Translate(float64(enemy.Position[0]-enemy.Size[0]), float64(enemy.Position[1]-40))
	hpPlayer := &ebiten.DrawImageOptions{}
	hpPlayer.GeoM.Translate(float64(player.Position[0]+40), float64(player.Position[1]-40))
	healthBarGreenP = ebiten.NewImage(300*player.Stats.currentHealth/player.Stats.maxHealth, 15)
	healthBarGreenP.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
	healthBarGreenE = ebiten.NewImage(300*enemy.Stats.currentHealth/enemy.Stats.maxHealth, 15)
	healthBarGreenE.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
	healthBarRed = ebiten.NewImage(300, 15)
	healthBarRed.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
	//Background DrawImage function
	screen.DrawImage(background, bg)
	//HP bar DrawImage function and conditionals
	screen.DrawImage(healthBarRed, hpEnemy)
	if !enemyDead {
		screen.DrawImage(healthBarGreenE, hpEnemy)
	}
	screen.DrawImage(healthBarRed, hpPlayer)
	if !playerDead {
		screen.DrawImage(healthBarGreenP, hpPlayer)
	}
	//Player and Enemy DrawImage function
	//screen.DrawImage(hero, char)
	// screen.DrawImage(bandit, npc)

	j := (g.count / currentAnimationE.framesPer) % currentAnimationE.frameNum
	bx := (currentAnimationE.frameOX + j*currentAnimationE.frameWidth) % currentAnimationE.maxWidth
	by := currentAnimationE.frameOY + ((currentAnimationE.frameOX+j*currentAnimationE.frameWidth)/currentAnimationE.maxWidth)*currentAnimationE.frameHeight
	screen.DrawImage(bandit.SubImage(image.Rect(bx, by, bx+currentAnimationE.frameWidth, by+currentAnimationE.frameHeight)).(*ebiten.Image), npc)

	//Array iteration to position a tileset using subimages and positions - Used for hotbar
	const xNum = (windowWidth / 2) / iconSize
	for _, l := range g.layers {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}

			op.GeoM.Translate(float64((i%xNum)*(iconSize)), float64((i/xNum)*(iconSize)))
			op.GeoM.Scale(2, 2)

			sx := (t % iconXNum) * iconSize
			sy := (t / iconXNum) * iconSize
			screen.DrawImage(iconsImg.SubImage(image.Rect(sx, sy, sx+iconSize, sy+iconSize)).(*ebiten.Image), op)
		}
	}
	//Arrow DrawImage function and hotkey Text drawn
	screen.DrawImage(arrow, arr)
	for i := 0; i < 8; i++ {
		text.Draw(screen, strconv.Itoa(i+1), mplusNormalFont, 428+(i*64), 890, color.White)
	}
	// Capture Mouse Position in variables
	x, y := ebiten.CursorPosition()
	//FPS shown top
	// Display Mouse Position
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f   X: %d, Y: %d", ebiten.CurrentTPS(), x, y))

	// i := (g.count / 7) % player.Anims.Idle.frameNum
	// sx, sy := player.Anims.Idle.frameOX+i*player.Anims.Idle.frameWidth, player.Anims.Idle.frameOY
	// screen.DrawImage(hero.SubImage(image.Rect(sx, sy, sx+player.Anims.Idle.frameWidth, sy+player.Anims.Idle.frameHeight)).(*ebiten.Image), ap)

	i := (g.count / currentAnimationP.framesPer) % currentAnimationP.frameNum
	sx := (currentAnimationP.frameOX + i*currentAnimationP.frameWidth) % currentAnimationP.maxWidth
	sy := currentAnimationP.frameOY + ((currentAnimationP.frameOX+i*currentAnimationP.frameWidth)/currentAnimationP.maxWidth)*currentAnimationP.frameHeight
	screen.DrawImage(hero.SubImage(image.Rect(sx, sy, sx+player.Anims.Idle.frameWidth, sy+player.Anims.Idle.frameHeight)).(*ebiten.Image), ap)

	//Fade in box drawn and alpha change
	box = ebiten.NewImage(windowWidth, windowHeight)
	box.Fill(color.RGBA{0, 0, 0, fadeAlpha})
	screen.DrawImage(box, nil) //test for fade in from black transparency
	//var p uint8 = 5
	if !fading {
		if fadeAlpha >= 10 {
			fadeAlpha -= 10
		} else if fadeAlpha < 10 {
			fadeAlpha = 0
		}
	} else if fading {
		if fadeAlpha <= 240 {
			fadeAlpha += 10
		} else if fadeAlpha > 240 {
			fadeAlpha = 255
			time.AfterFunc(FadeTime, func() {
				fading = false
			})
		}
	}

	//Turn Text DrawImage function
	text.Draw(screen, turnText, mplusBigFont, windowWidth/2-(text.BoundString(mplusBigFont, turnText).Dx()/2), windowHeight*1/20, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	z := 400
	g := &Game{
		layers: [][]int{
			{
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
				z, z, z, z, z, z, 81, 48, 148, 59, 11, 11, 11, 11, z, z, z, z, z, z,
				z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z, z,
			},
		},
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("2D Game")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

/*
END BATTLE AT 0 HP
SETUP FADE IN FROM BLACK ANIMATION
SETUP EXPERIENCE BAR AND LEVELS
SETUP STATS INCREASE SCREEN AFTER BATTLE
SETUP CONTINUOUS LOOP FADE IN -> BATTLE -> FADE OUT/FADE IN -> EXP STAT UPGRADE -> FADE OUT/FADE IN -> BATTLE LOOPS
*/

/*
start turn
fade in animation "box"
idle animations
choose spell
g.count reset to start animations clean
< move player position / walk to enemy
< change player animation to walking
< arrive at enemy position
< change player animation to attacking
< flip player model
< change player animation to walking
< move player position / walk back to original spot
< change player animation to idle
< start enemy turn
-change idle animation to attack animation
-goroutine set to 3 second delay for enemy turn
-change attack aniamtion to idle
set turn to player
loop

*/
