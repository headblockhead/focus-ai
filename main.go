package main

import (
	"github.com/headblockhead/focus-ai/game"
	"github.com/headblockhead/focus-ai/visualizer"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Focus AI Visualizer")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	board := game.NewBoard()

	vis := visualizer.NewVisualizer(&board)

	// An ungraceful quit
	if err := ebiten.RunGame(vis); err != nil && err.Error() != "quit" {
		panic(err)
	}
}
