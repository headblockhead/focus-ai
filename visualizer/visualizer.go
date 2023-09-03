package visualizer

import (
	"errors"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/headblockhead/focus-ai/game"
	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"
)

type Visualizer struct {
	Width, Height  int
	Library        *tetra3d.Library
	Scene          *tetra3d.Scene
	DrawDebugStats bool
	Board          *game.Board
	RedPieces      [8][8][5]*tetra3d.Model
	GreenPieces    [8][8][5]*tetra3d.Model
}

//go:embed startingScene.gltf
var startingGLTF []byte

func NewVisualizer(board *game.Board) *Visualizer {
	vis := &Visualizer{
		Width:  3840,
		Height: 2160,
		Board:  board,
	}
	vis.Init()
	return vis
}

func (vis *Visualizer) Init() {
	if vis.Library == nil {
		options := tetra3d.DefaultGLTFLoadOptions()
		options.CameraWidth = vis.Width
		options.CameraHeight = vis.Height
		library, err := tetra3d.LoadGLTFData(startingGLTF, options)
		if err != nil {
			panic(err)
		}

		vis.Library = library
	}
	vis.Scene = vis.Library.ExportedScene.Clone()

	// Get the green piece
	greenPiece := vis.Scene.Root.Get("GreenPiece").(*tetra3d.Model)
	// Create 320 copies of the green piece
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			for k := 0; k < 5; k++ {
				piece := greenPiece.Clone().(*tetra3d.Model)
				vis.GreenPieces[i][j][k] = piece
				vis.Scene.Root.AddChildren(piece)
				piece.SetWorldPosition((float64(i)*2)-8, (float64(k)*0.4)+0.2, (float64(j)*2)-6)
				piece.SetVisible(false, true)
			}
		}
	}
	// Get the red piece
	redPiece := vis.Scene.Root.Get("RedPiece").(*tetra3d.Model)
	// Create 320 copies of the red piece
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			for k := 0; k < 5; k++ {
				piece := redPiece.Clone().(*tetra3d.Model)
				vis.RedPieces[i][j][k] = piece
				vis.Scene.Root.AddChildren(piece)
				piece.SetWorldPosition((float64(i)*2)-8, (float64(k)*0.4)+0.2, (float64(j)*2)-6)
				piece.SetVisible(false, true)
			}
		}
	}

	// Delete the original RedPiece and GreenPiece
	vis.Scene.Root.RemoveChildren(greenPiece)
	vis.Scene.Root.RemoveChildren(redPiece)
}

func (vis *Visualizer) Update() (err error) {
	// Update the scene based on the current Board
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			tile := vis.Board.Tiles[i][j]
			for k := 0; k < len(tile.Pieces); k++ {
				if !tile.Pieces[k].Exists {
					continue
				}
				if tile.Pieces[k].Color == game.RED {
					vis.RedPieces[i][j][k].SetVisible(true, true)
				} else {
					vis.GreenPieces[i][j][k].SetVisible(true, true)
				}
			}
		}
	}

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
		vis.DrawDebugStats = !vis.DrawDebugStats
	}

	return err
}

func (vis *Visualizer) Draw(screen *ebiten.Image) {

	screen.Fill(vis.Scene.World.ClearColor.ToRGBA64())

	camera := vis.Scene.Root.Get("Camera").(*tetra3d.Camera)

	camera.Clear()
	camera.RenderNodes(vis.Scene, vis.Scene.Root)
	screen.DrawImage(camera.ColorTexture(), nil)

	if vis.DrawDebugStats {
		camera.DrawDebugRenderInfo(screen, 1, colors.White())
	}

}

func (vis *Visualizer) Layout(w, h int) (int, int) {
	// This is a fixed aspect ratio; we can change this to, say, extend for wider displays by using the provided w argument and
	// calculating the height from the aspect ratio, then calling Camera.Resize() on any / all cameras with the new width and height.
	return vis.Width, vis.Height
}
