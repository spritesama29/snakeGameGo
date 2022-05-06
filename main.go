package main

import (
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font/basicfont"
	stdImage "image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

//go:embed assets/*
var EmbeddedAssets embed.FS
var simpleGame Game
var textWidget *widget.Text
var button *widget.Button
var button2 *widget.Button
var teleport = false
var counter = 0
var fillTime = false
var started = false
var counte = 0
var length = 5
var foodCount = 0
var badFoodCount = 0

const (
	GameWidth   = 700
	GameHeight  = 700
	PlayerSpeed = 1
)

type cords struct {
	xloc int
	yloc int
}
type Sprite struct {
	pict         *ebiten.Image
	xloc         int
	yloc         int
	dX           int
	dY           int
	segmentNum   int
	coords       cords
	direction    string
	teleported   bool
	countingTime int

	drawOps ebiten.DrawImageOptions
}

type Game struct {
	player      Sprite
	enemy       Sprite
	enemy2      Sprite
	tail        Sprite
	score       int
	drawOps     ebiten.DrawImageOptions
	enemyList   []Sprite
	enemyList2  []Sprite
	segmentList []*Sprite
	enemyCount  int
	AppUI       *ebitenui.UI
}

func (g *Game) Update() error {
	processPlayerInput(g)

	g.AppUI.Update()
	//print("xloc:", g.player.xloc)

	//print("  yloc:", g.player.yloc)
	//print("\n")

	//timer1 := time.NewTimer(2 * time.Second)
	//<-timer1.C
	return nil
}
func foodTime(g *Game) bool {
	ticker := time.NewTicker(3 * time.Second)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _ = range ticker.C {
		g.enemyList[foodCount].xloc = r.Intn(650)
		g.enemyList[foodCount].yloc = r.Intn(650)

		foodCount += 1

	}
	return true
}
func badFoodTime(g *Game) {
	ticker := time.NewTicker(8 * time.Second)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _ = range ticker.C {
		g.enemyList2[foodCount].xloc = r.Intn(650)
		g.enemyList2[foodCount].yloc = r.Intn(650)

		badFoodCount += 1

	}
}
func appleCollison(player Sprite, apple Sprite, g *Game) bool {
	enemyWidth, enemyHeight := apple.pict.Size()
	playerWidth, playerHeight := g.player.pict.Size()

	if player.xloc < apple.xloc+enemyWidth && player.xloc+playerWidth > apple.xloc-enemyWidth &&
		player.yloc < apple.yloc+enemyHeight && player.yloc+playerHeight > apple.yloc-enemyHeight {

		return true
	}
	return false
}
func yCheck(y1 int, y2 int) int {
	ycheck := y1 - y2
	if ycheck < 0 {
		ycheck = 0 - ycheck
	}
	return ycheck
}
func (g *Game) Draw(screen *ebiten.Image) {

	g.AppUI.Draw(screen)
	g.drawOps.GeoM.Reset()
	g.drawOps.GeoM.Translate(float64(g.segmentList[0].xloc), float64(g.segmentList[0].yloc))
	screen.DrawImage(g.player.pict, &g.drawOps)
	//g.drawOps.GeoM.Scale(.5, .5)
	if length < 5 {
		println("You got too short!")
		os.Exit(69)
	}
	for num, _ := range g.segmentList {

		//_,_ := g.player.pict.Size()
		counte += 1
		//println(counte)
		for j := length; j < 1000; j++ {
			g.segmentList[j].xloc = 10000
			g.segmentList[j].yloc = 10000
		}
		if num < length {

			if g.segmentList[0].direction == "left" {
				ycheck := yCheck(g.segmentList[0].yloc, g.segmentList[num].yloc)
				xcheck := g.segmentList[0].xloc - g.segmentList[num].xloc
				if xcheck < 0 {
					xcheck = 0 - xcheck
				}
				if g.segmentList[0].xloc > g.segmentList[num].xloc && xcheck < 10 && ycheck < 10 {
					os.Exit(69)
				}

			} else if g.segmentList[0].direction == "right" {
				ycheck := g.segmentList[0].yloc - g.segmentList[num].yloc
				if ycheck < 0 {
					ycheck = 0 - ycheck
				}
				xcheck := g.segmentList[0].xloc - g.segmentList[num].xloc
				if xcheck < 0 {
					xcheck = 0 - xcheck
				}
				if g.segmentList[0].xloc < g.segmentList[num].xloc && xcheck < 10 && ycheck < 10 {
					os.Exit(69)
				}

			} else if g.segmentList[0].direction == "up" {
				xcheck := g.segmentList[0].xloc - g.segmentList[num].xloc
				if xcheck < 0 {
					xcheck = 0 - xcheck
				}
				ycheck := g.segmentList[0].yloc - g.segmentList[num].yloc
				if ycheck < 0 {
					ycheck = 0 - ycheck
				}
				if g.segmentList[0].yloc > g.segmentList[num].yloc && ycheck < 10 && xcheck < 10 {
					os.Exit(69)
				}

			} else if g.segmentList[0].direction == "down" {
				xcheck := g.segmentList[0].xloc - g.segmentList[num].xloc
				if xcheck < 0 {
					xcheck = 0 - xcheck
				}
				ycheck := g.segmentList[0].yloc - g.segmentList[num].yloc
				if ycheck < 0 {
					ycheck = 0 - ycheck
				}
				if g.segmentList[0].yloc < g.segmentList[num].yloc && ycheck < 10 && xcheck < 10 {
					os.Exit(69)
				}

			}

			g.drawOps.GeoM.Reset()
			g.drawOps.GeoM.Translate(float64(g.segmentList[num].xloc), float64(g.segmentList[num].yloc))
			screen.DrawImage(g.player.pict, &g.drawOps)

			if num == 0 {

				g.segmentList[num].coords.yloc = g.segmentList[num].yloc
				g.segmentList[num].coords.xloc = g.segmentList[num].xloc

			} else if g.segmentList[length-1].teleported == true {

				for i := 0; i < length; i++ {
					g.segmentList[i].teleported = false

				}
			} else if g.segmentList[num].teleported == false && counte > 1000 {

				g.segmentList[num].xloc = g.segmentList[0].coords.xloc
				g.segmentList[num].yloc = g.segmentList[0].coords.yloc
				//g.segmentList[num].coords.xloc = g.segmentList[num].xloc
				//g.segmentList[num].coords.yloc = g.segmentList[num].yloc
				teleport = false
				g.segmentList[num].teleported = true
				counte = 0

			}

		}

	}
	textWidget.SetLocation(stdImage.Rectangle{
		Min: stdImage.Point{
			X: 0,
			Y: 600},
		Max: stdImage.Point{
			X: 700,
			Y: 700,
		},
	})
	textWidget.Label = fmt.Sprintf("Score:%d", counter)
	button.SetLocation(stdImage.Rectangle{
		Min: stdImage.Point{
			X: 0,
			Y: -100},
		Max: stdImage.Point{
			X: 1000,
			Y: -100,
		},
	})
	button.Text().Label = ""

	button2.SetLocation(stdImage.Rectangle{
		Min: stdImage.Point{
			X: 0,
			Y: -100},
		Max: stdImage.Point{
			X: 1000,
			Y: -100,
		},
	})
	button2.Text().Label = ""
	// This collision detection is from jsantore firstGameDemo
	for num1, enemy := range g.enemyList {
		enemyWidth, enemyHeight := enemy.pict.Size()
		playerWidth, playerHeight := g.player.pict.Size()

		if g.segmentList[0].xloc < enemy.xloc+enemyWidth && g.segmentList[0].xloc+playerWidth > enemy.xloc-enemyWidth &&
			g.segmentList[0].yloc < enemy.yloc+enemyHeight && g.segmentList[0].yloc+playerHeight > enemy.yloc-enemyHeight {
			//Collison from jsantore firstGameDemo repo on github
			remove(g.enemyList, num1)
			length += 10
			counter += 1
			//message := fmt.Sprintf("score: %d", counter)
			//textWidget.Label = message

		} else {
			g.drawOps.GeoM.Reset()
			g.drawOps.GeoM.Translate(float64(enemy.xloc), float64(enemy.yloc))
			screen.DrawImage(enemy.pict, &g.drawOps)

		}

	}
	for num2, enemy := range g.enemyList2 {

		if appleCollison(g.player, enemy, g) {
			//Collison from jsantore firstGameDemo repo on github
			remove(g.enemyList2, num2)
			length -= 20
			counter += 1
			//message := fmt.Sprintf("score: %d", counter)
			//textWidget.Label = message

		} else {
			g.drawOps.GeoM.Reset()
			g.drawOps.GeoM.Translate(float64(enemy.xloc), float64(enemy.yloc))
			screen.DrawImage(enemy.pict, &g.drawOps)

		}

	}
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}
func remove(s []Sprite, index int) []Sprite {
	return append(s[:index], s[index+1:]...)
	//Remove function found here https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
}
func fillList(g *Game, amount int) int {

	for i := 0; i < amount; i++ {
		g.enemy = Sprite{
			pict: loadPNGImageFromEmbedded("apple.png"),
			xloc: 10000,
			yloc: 10000,
			dX:   0,
			dY:   0,
		}
		g.enemyList[i] = g.enemy

	}

	return amount + 1
}
func fillListBad(g *Game) {

	for i := 0; i < 1000; i++ {
		g.enemy2 = Sprite{
			pict: loadPNGImageFromEmbedded("badApple.png"),
			xloc: 10000,
			yloc: 10000,
			dX:   0,
			dY:   0,
		}
		g.enemyList2[i] = g.enemy2

	}
}

func fillListTail(g *Game) {

	for i := 1; i < 9999; i++ {
		tail := Sprite{
			pict:         loadPNGImageFromEmbedded("blueBox.png"),
			xloc:         10000,
			yloc:         10000,
			dX:           0,
			dY:           0,
			teleported:   false,
			countingTime: i,
		}

		g.segmentList[i] = &tail

	}
}
func main() {
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle("Minimal Game")

	simpleGame := Game{AppUI: MakeUIWindow()}

	simpleGame.player = Sprite{
		pict:      loadPNGImageFromEmbedded("blueBox.png"),
		xloc:      200,
		yloc:      300,
		dX:        0,
		dY:        0,
		direction: "null",
	}
	simpleGame.segmentList = make([]*Sprite, 10000)
	simpleGame.segmentList[0] = &simpleGame.player
	simpleGame.enemyList = make([]Sprite, 1001)
	simpleGame.enemyList2 = make([]Sprite, 1001)

	fillList(&simpleGame, 1000)
	fillListTail(&simpleGame)
	fillListBad(&simpleGame)
	simpleGame.enemy = Sprite{
		pict: loadPNGImageFromEmbedded("smallhammer.png"),
		xloc: 10000,
		yloc: 10000,
		dX:   0,
		dY:   0,
	}
	simpleGame.enemy2 = Sprite{
		pict: loadPNGImageFromEmbedded("smallhammer.png"),
		xloc: 10000,
		yloc: 10000,
		dX:   0,
		dY:   0,
	}
	go foodTime(&simpleGame)
	go badFoodTime(&simpleGame)
	simpleGame.enemyList[1000] = simpleGame.enemy
	simpleGame.enemyList2[1000] = simpleGame.enemy2
	if err := ebiten.RunGame(&simpleGame); err != nil {
		log.Fatal("Oh no! something terrible happened and the game crashed", err)
	}
	textInfo := widget.TextOptions{}.Text("score: 0", basicfont.Face7x13, color.White)

	textWidget = widget.NewText(textInfo)
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("assets")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("assets/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func processPlayerInput(theGame *Game) {
	if theGame.player.direction == "null" {
		theGame.player.dX = PlayerSpeed
		theGame.player.dY = 0
		theGame.player.direction = "right"

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && theGame.player.direction != "down" {
		theGame.player.dY = -PlayerSpeed
		theGame.player.dX = 0
		theGame.player.direction = "up"

	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) && theGame.player.direction != "up" {
		theGame.player.dY = PlayerSpeed
		theGame.player.dX = 0
		theGame.player.direction = "down"

	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && theGame.player.direction != "right" {
		theGame.player.dX = -PlayerSpeed
		theGame.player.dY = 0
		theGame.player.direction = "left"

	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) && theGame.player.direction != "left" {
		theGame.player.dX = PlayerSpeed
		theGame.player.dY = 0
		theGame.player.direction = "right"

	} //else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight) {
	//theGame.player.dX = 0
	//}else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
	//		theGame.player.dY = 0
	//	}
	theGame.player.yloc += theGame.player.dY
	theGame.player.xloc += theGame.player.dX
	if theGame.player.yloc <= 0 {
		theGame.player.dY = 0
		theGame.player.yloc = 0
	} else if theGame.player.yloc > 675 {
		theGame.player.dY = 0
		theGame.player.yloc = 675

	}
	if theGame.player.xloc <= 0 {
		theGame.player.dX = 0
		theGame.player.xloc = 0
	} else if theGame.player.xloc > GameWidth-30 {
		theGame.player.dX = 0
		theGame.player.xloc = GameWidth - 30
	}
}

func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	background := image.NewNineSliceColor(color.Gray16{})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),

			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(background))
	textInfo := widget.TextOptions{}.Text("score: 0", basicfont.Face7x13, color.White)

	idle, err := loadImageNineSlice("button-idle.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	hover, err := loadImageNineSlice("button-hover.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pressed, err := loadImageNineSlice("button-pressed.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	disabled, err := loadImageNineSlice("button-disabled.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}

	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("I want to go :(", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		// ... click handler, etc. ...
		widget.ButtonOpts.ClickedHandler(quit),
	)
	button2 = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Let's Play again!", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display

		// ... click handler, etc. ...
		widget.ButtonOpts.ClickedHandler(playAgain),
	)

	rootContainer.AddChild(button)
	rootContainer.AddChild(button2)
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)
	GUIhandler = &ebitenui.UI{Container: rootContainer}

	return GUIhandler
}

func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func playAgain(args *widget.ButtonClickedEventArgs) {

	fillTime = true
	counter = 0
}
func quit(args *widget.ButtonClickedEventArgs) {
	os.Exit(3)
}
