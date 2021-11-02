package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type Game struct{}

func NewGame() *Game{
	g := &Game{}
	return g
}

func (g *Game) Update() error{
	//If camera is active, update camera
	//
	return nil
}

func (g *Game) Layout(int, int)(screenWidth, screenLength int){
	return 600, 800
}

func (g *Game) Draw(screen *ebiten.Image){}

func main(){
	g := NewGame()
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("App")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}