package main

import (
	"errors"

	_ "embed"

	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Width, Height  int
	Library        *tetra3d.Library
	Scene          *tetra3d.Scene
	DrawDebugDepth bool
	DrawDebugStats bool
}

//go:embed startingScene.gltf
var startingGLTF []byte

func NewGame() *Game {

	game := &Game{
		Width:  796,
		Height: 448,
	}

	game.Init()

	return game
}

func (g *Game) Init() {

	if g.Library == nil {

		options := tetra3d.DefaultGLTFLoadOptions()
		options.CameraWidth = g.Width
		options.CameraHeight = g.Height
		library, err := tetra3d.LoadGLTFData(startingGLTF, options)
		if err != nil {
			panic(err)
		}

		g.Library = library

	}

	g.Scene = g.Library.ExportedScene.Clone()

}

func (g *Game) Update() error {
	var err error

	plate := g.Scene.Root.Get("Plate").(*tetra3d.Model)
	plate.Rotate(0, 0.5, 0, tetra3d.ToRadians(0.5))

	// Quit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		err = errors.New("quit")
	}

	// Fullscreen
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Stats for nerds
	if inpututil.IsKeyJustPressed(ebiten.KeyF8) {
		g.DrawDebugStats = !g.DrawDebugStats
	}

	return err
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(g.Scene.World.ClearColor.ToRGBA64())

	camera := g.Scene.Root.Get("Camera").(*tetra3d.Camera)

	camera.Clear()
	camera.RenderNodes(g.Scene, g.Scene.Root)
	screen.DrawImage(camera.ColorTexture(), nil)

	if g.DrawDebugStats {
		camera.DrawDebugRenderInfo(screen, 1, colors.White())
	}

}

func (g *Game) Layout(w, h int) (int, int) {
	// This is a fixed aspect ratio; we can change this to, say, extend for wider displays by using the provided w argument and
	// calculating the height from the aspect ratio, then calling Camera.Resize() on any / all cameras with the new width and height.
	return g.Width, g.Height
}

func main() {

	ebiten.SetWindowTitle("Focus AI Visualizer")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	// An ungraceful quit
	if err := ebiten.RunGame(game); err != nil && err.Error() != "quit" {
		panic(err)
	}

}
