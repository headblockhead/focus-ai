package game

import (
	"fmt"
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()
	if b.ReservesR != 0 {
		t.Errorf("Expected 0, got %d", b.ReservesR)
	}
	if b.ReservesG != 0 {
		t.Errorf("Expected 0, got %d", b.ReservesG)
	}
	unusableTiles := [][2]int{
		{0, 0},
		{0, 1},
		{1, 0},
		{6, 0},
		{7, 0},
		{7, 1},
		{0, 6},
		{0, 7},
		{1, 7},
		{6, 7},
		{7, 6},
		{7, 7},
	}
	for i := 0; i < len(unusableTiles); i++ {
		if b.Tiles[unusableTiles[i][0]][unusableTiles[i][1]].useable {
			t.Errorf("Expected unusable tile at %d, %d", unusableTiles[i][0], unusableTiles[i][1])
		}
	}

}

func runForEveryTile(f func(x int, y int)) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			f(i, j)
		}
	}
}

func TestGetTile(t *testing.T) {
	b := NewBoard()
	runForEveryTile(func(i int, j int) {
		tile, err := b.GetTile(i, j)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if tile == nil {
			t.Errorf("Expected tile, got nil")
		}
	})
	runForEveryTile(func(i int, j int) {
		b.Tiles[i][j] = Tile{Pieces: [5]Piece{{Color: RED}, {Color: GREEN}}}
		tile, err := b.GetTile(i, j)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if tile.Pieces[0].Color != RED {
			t.Errorf("Expected tile with red piece at 0, got %v", tile)
		}
		if tile.Pieces[1].Color != GREEN {
			t.Errorf("Expected tile with green piece at 1, got %v", tile)
		}
	})
	tile, err := b.GetTile(-1, 0)
	if err != ErrTileOutOfBounds {
		t.Errorf("Expected GetTileErrorOutOfBounds, got %v", err)
	}
	if tile != nil {
		t.Errorf("Expected nil, got %v", tile)
	}
}

func TestGetReserves(t *testing.T) {
	b := NewBoard()
	if b.GetReserves(RED) == nil {
		t.Errorf("Expected reserves, got nil")
	}
	if b.GetReserves(GREEN) == nil {
		t.Errorf("Expected reserves, got nil")
	}
	b.ReservesR = 5
	b.ReservesG = 6
	if *b.GetReserves(RED) != 5 {
		t.Errorf("Expected 5, got %d", *b.GetReserves(RED))
	}
	if *b.GetReserves(GREEN) != 6 {
		t.Errorf("Expected 6, got %d", *b.GetReserves(GREEN))
	}
}

func TestSetReserves(t *testing.T) {
	b := NewBoard()
	b.SetReserves(RED, 5)
	b.SetReserves(GREEN, 6)
	if b.ReservesR != 5 {
		t.Errorf("Expected 5, got %d", b.ReservesR)
	}
	if b.ReservesG != 6 {
		t.Errorf("Expected 6, got %d", b.ReservesG)
	}
}

func TestAddToReserves(t *testing.T) {
	b := NewBoard()
	b.AddToReserves(RED, 5)
	b.AddToReserves(GREEN, 6)
	if b.ReservesR != 5 {
		t.Errorf("Expected 5, got %d", b.ReservesR)
	}
	if b.ReservesG != 6 {
		t.Errorf("Expected 6, got %d", b.ReservesG)
	}
	b.AddToReserves(RED, 5)
	b.AddToReserves(GREEN, 6)
	if b.ReservesR != 10 {
		t.Errorf("Expected 10, got %d", b.ReservesR)
	}
	if b.ReservesG != 12 {
		t.Errorf("Expected 12, got %d", b.ReservesG)
	}
}

func TestSetTile(t *testing.T) {
	b := NewBoard()
	b.SetTile(6, 5, Tile{useable: false})
	if b.Tiles[6][5].useable {
		t.Errorf("Expected unusable tile at 6, 5")
	}
	b.SetTile(6, 5, Tile{useable: true})
	if !b.Tiles[6][5].useable {
		t.Errorf("Expected usable tile at 6, 5")
	}
}

func checkPieces(pieces [5]Piece, expected [5]Piece) (err error) {
	for i := 0; i < 5; i++ {
		if pieces[i].Color != expected[i].Color {
			return fmt.Errorf("Expected %v, got %v", expected, pieces)
		}
		if pieces[i].Exists != expected[i].Exists {
			return fmt.Errorf("Expected %v, got %v", expected, pieces)
		}
	}
	return nil
}

func TestAddPiece(t *testing.T) {
	b := NewBoard()
	err := b.AddPiece(6, 5, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{},
		{},
		{},
		{},
	})
	if err != nil {
		t.Error(err)
	}
	err = b.AddPiece(6, 5, Piece{Color: GREEN, Exists: true}, GREEN)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
		{},
		{},
		{},
	})
	if err != nil {
		t.Error(err)
	}
	err = b.AddPiece(6, 5, Piece{Color: GREEN, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
		{},
		{},
	})
	if err != nil {
		t.Error(err)
	}
	err = b.AddPiece(6, 5, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: RED, Exists: true},
		{},
	})
	if err != nil {
		t.Error(err)
	}
	err = b.AddPiece(6, 5, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: RED, Exists: true},
		{Color: RED, Exists: true},
	})
	if err != nil {
		t.Error(err)
	}
	err = b.AddPiece(6, 5, Piece{Color: GREEN, Exists: true}, GREEN)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: RED, Exists: true},
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
	})
	if err != nil {
		t.Error(err)
	}
	// 0 RED reserves, 0 GREEN reserves, as GREEN player took RED piece, so it is captured, not reserved.
	if b.ReservesR != 0 {
		t.Errorf("Expected 0 RED reserves, got %d", b.ReservesR)
	}
	if b.ReservesG != 0 {
		t.Errorf("Expected 0 GREEN reserves, got %d", b.ReservesG)
	}
	err = b.AddPiece(6, 5, Piece{Color: GREEN, Exists: true}, GREEN)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: GREEN, Exists: true},
		{Color: RED, Exists: true},
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
	})
	if err != nil {
		t.Error(err)
	}
	// 0 RED reserves, 1 GREEN reserve, as GREEN player took GREEN piece, so it is reserved.
	if b.ReservesR != 0 {
		t.Errorf("Expected 0 RED reserves, got %d", b.ReservesR)
	}
	if b.ReservesG != 1 {
		t.Errorf("Expected 1 GREEN reserve, got %d", b.ReservesG)
	}

	err = b.AddPiece(6, 5, Piece{Color: GREEN, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
	})
	if err != nil {
		t.Error(err)
	}
	// 0 RED reserves, (1 GREEN reserve), as RED player took GREEN piece, so it is captured, not reserved.
	if b.ReservesR != 0 {
		t.Errorf("Expected 0 RED reserves, got %d", b.ReservesR)
	}
	if b.ReservesG != 1 {
		t.Errorf("Expected 1 GREEN reserve, got %d", b.ReservesG)
	}

	err = b.AddPiece(6, 5, Piece{Color: GREEN, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[6][5].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
		{Color: GREEN, Exists: true},
	})
	if err != nil {
		t.Error(err)
	}
	// 1 RED reserves, (1 GREEN reserve), as RED player took RED piece, so it is reserved.
	if b.ReservesR != 1 {
		t.Errorf("Expected 1 RED reserves, got %d", b.ReservesR)
	}
	if b.ReservesG != 1 {
		t.Errorf("Expected 1 GREEN reserve, got %d", b.ReservesG)
	}
}

func TestAddPieceNonExistant(t *testing.T) {
	b := NewBoard()
	err := b.AddPiece(6, 5, Piece{Color: RED, Exists: false}, RED) // playerColor doesn't matter here.
	if err != ErrPieceNonExistent {
		t.Errorf("Expected AddPieceErrorPieceNonExistant, got %v", err)
	}
}

func TestAddFromReserves(t *testing.T) {
	b := NewBoard()
	b.SetReserves(RED, 5)
	b.SetReserves(GREEN, 6)
	err := b.AddFromReserves(RED, 4, 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[4][3].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{},
		{},
		{},
		{},
	})
	if err != nil {
		t.Error(err)
	}
	if b.ReservesR != 4 {
		t.Errorf("Expected 4 RED reserves, got %d", b.ReservesR)
	}
	if b.ReservesG != 6 {
		t.Errorf("Expected 6 GREEN reserves, got %d", b.ReservesG)
	}
	err = b.AddFromReserves(GREEN, 4, 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[4][3].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
		{},
		{},
		{},
	})
	if err != nil {
		t.Error(err)
	}
	if b.ReservesR != 4 {
		t.Errorf("Expected 4 RED reserves, got %d", b.ReservesR)
	}
	if b.ReservesG != 5 {
		t.Errorf("Expected 5 GREEN reserves, got %d", b.ReservesG)
	}

	err = b.AddPiece(4, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.AddPiece(4, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.AddPiece(4, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = b.AddFromReserves(GREEN, 4, 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = checkPieces(b.Tiles[4][3].Pieces, [5]Piece{
		{Color: GREEN, Exists: true},
		{Color: RED, Exists: true},
		{Color: RED, Exists: true},
		{Color: RED, Exists: true},
		{Color: GREEN, Exists: true},
	})
	if err != nil {
		t.Error(err)
	}
	if b.ReservesR != 4 {
		t.Errorf("Expected 4 RED reserves, got %d", b.ReservesR)
	}
	if b.ReservesG != 4 {
		t.Errorf("Expected 4 GREEN reserves, got %d", b.ReservesG)
	}
}

func TestMoveUnusableTile(t *testing.T) {
	b := NewBoard()
	for i := 0; i < 4; i++ {
		err := b.Move(0, 0, 1, []Direction{Direction(i)}, RED)
		if err != ErrTileSourceNonUsable {
			t.Errorf("Expected MoveErrorTileSourceNonUsable, got %v", err)
		}
	}

	err := b.AddPiece(0, 2, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(0, 2, 1, []Direction{UP}, RED)
	if err != ErrTileDestinationUnusable {
		t.Errorf("Expected MoveErrorTileDestNonUsable, got %v", err)
	}

	err = b.AddPiece(0, 5, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(0, 5, 1, []Direction{DOWN}, RED)
	if err != ErrTileDestinationUnusable {
		t.Errorf("Expected MoveErrorTileDestNonUsable, got %v", err)
	}

	err = b.AddPiece(2, 0, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(2, 0, 1, []Direction{LEFT}, RED)
	if err != ErrTileDestinationUnusable {
		t.Errorf("Expected MoveErrorTileDestNonUsable, got %v", err)
	}

	err = b.AddPiece(5, 0, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(5, 0, 1, []Direction{RIGHT}, RED)
	if err != ErrTileDestinationUnusable {
		t.Errorf("Expected MoveErrorTileDestNonUsable, got %v", err)
	}

	err = b.AddPiece(5, 1, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.AddPiece(5, 1, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(5, 1, 2, []Direction{UP, RIGHT}, RED)
	if err != ErrTileDestinationUnusable {
		t.Errorf("Expected MoveErrorTileDestNonUsable, got %v", err)
	}
}

func TestMoveOutOfBounds(t *testing.T) {
	b := NewBoard()
	err := b.Move(-1, 0, 1, []Direction{UP}, RED)
	if err != ErrTileOutOfBounds {
		t.Errorf("Expected GetTileErrorOutOfBounds, got %v", err)
	}
	err = b.Move(0, -1, 1, []Direction{UP}, RED)
	if err != ErrTileOutOfBounds {
		t.Errorf("Expected GetTileErrorOutOfBounds, got %v", err)
	}
	err = b.Move(8, 0, 1, []Direction{UP}, RED)
	if err != ErrTileOutOfBounds {
		t.Errorf("Expected GetTileErrorOutOfBounds, got %v", err)
	}
	err = b.Move(0, 8, 1, []Direction{UP}, RED)
	if err != ErrTileOutOfBounds {
		t.Errorf("Expected GetTileErrorOutOfBounds, got %v", err)
	}
}

func TestMoveNoPiece(t *testing.T) {
	b := NewBoard()
	err := b.Move(3, 3, 1, []Direction{UP}, RED)
	if err != ErrNoPieceToMove {
		t.Errorf("Expected ErrNoPieceToMove, got %v", err)
	}
	err = b.AddPiece(3, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(3, 3, 0, []Direction{UP}, RED)
	if err != ErrMustMoveAtLeastOnePiece {
		t.Errorf("Expected MustMoveAtLeastOnePiece, got %v", err)
	}
}

func TestMoveNoDirections(t *testing.T) {
	b := NewBoard()
	err := b.AddPiece(3, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(3, 3, 1, []Direction{}, RED)
	if err != ErrNoDirections {
		t.Errorf("Expected ErrNoDirections, got %v", err)
	}
}

func TestMoveWrongAmountDirections(t *testing.T) {
	b := NewBoard()
	err := b.AddPiece(3, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(3, 3, 1, []Direction{UP, UP}, RED)
	if err != ErrWrongDirectionAmount {
		t.Errorf("Expected ErrWrongDirectionAmount, got %v", err)
	}
}

func TestMoveTooManyPieces(t *testing.T) {
	b := NewBoard()
	err := b.AddPiece(3, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(3, 3, 2, []Direction{UP, UP}, RED)
	if err != ErrTooManyPieces {
		t.Errorf("Expected ErrorTooManyPieces, got %v", err)
	}
}

func TestMoveDifferentColor(t *testing.T) {
	b := NewBoard()
	err := b.AddPiece(3, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(3, 3, 1, []Direction{UP}, GREEN)
	if err != ErrWrongColor {
		t.Errorf("Expected MoveErrorWrongColor, got %v", err)
	}
	b = NewBoard()
	err = b.AddPiece(3, 3, Piece{Color: RED, Exists: true}, GREEN)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.AddPiece(3, 3, Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = b.Move(3, 3, 1, []Direction{UP}, GREEN)
	if err != ErrWrongColor {
		t.Errorf("Expected MoveErrorWrongColor, got %v", err)
	}
}

func MoveDirectionTest(initalLocation [2]int, direction Direction) (err error) {
	locationModifier := [2]int{0, 0}
	if direction == UP {
		locationModifier = [2]int{0, -1}
	} else if direction == DOWN {
		locationModifier = [2]int{0, 1}
	} else if direction == LEFT {
		locationModifier = [2]int{-1, 0}
	} else if direction == RIGHT {
		locationModifier = [2]int{1, 0}
	}

	locationModifier[0] *= 4
	locationModifier[1] *= 4

	b := NewBoard()
	err = b.AddPiece(initalLocation[0], initalLocation[1], Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		return err
	}
	err = b.Move(initalLocation[0], initalLocation[1], 1, []Direction{direction}, RED)
	if err != nil {
		return err
	}
	err = checkPieces(b.Tiles[initalLocation[0]+(locationModifier[0]/4)][initalLocation[1]+(locationModifier[1]/4)].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{},
		{},
		{},
		{},
	})
	if err != nil {
		return err
	}
	err = checkPieces(b.Tiles[initalLocation[0]][initalLocation[1]].Pieces, [5]Piece{
		{},
		{},
		{},
		{},
		{},
	})
	if err != nil {
		return err
	}
	b = NewBoard()
	err = b.AddPiece(initalLocation[0], initalLocation[1], Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		return err
	}
	err = b.AddPiece(initalLocation[0], initalLocation[1], Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		return err
	}
	err = b.Move(initalLocation[0], initalLocation[1], 1, []Direction{direction}, RED)
	if err != nil {
		return err
	}
	err = checkPieces(b.Tiles[initalLocation[0]][initalLocation[1]].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{},
		{},
		{},
		{},
	})
	if err != nil {
		return err
	}
	err = checkPieces(b.Tiles[initalLocation[0]+(locationModifier[0]/4)][initalLocation[1]+(locationModifier[1]/4)].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{},
		{},
		{},
		{},
	})
	if err != nil {
		return err
	}

	b = NewBoard()
	err = b.AddPiece(initalLocation[0], initalLocation[1], Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		return err
	}
	err = b.AddPiece(initalLocation[0], initalLocation[1], Piece{Color: RED, Exists: true}, RED)
	if err != nil {
		return err
	}
	err = b.Move(initalLocation[0], initalLocation[1], 2, []Direction{direction, direction}, RED)
	if err != nil {
		return err
	}
	err = checkPieces(b.Tiles[initalLocation[0]][initalLocation[1]].Pieces, [5]Piece{
		{},
		{},
		{},
		{},
		{},
	})
	if err != nil {
		return err
	}
	err = checkPieces(b.Tiles[initalLocation[0]+(locationModifier[0]/4)][initalLocation[1]+(locationModifier[1]/4)].Pieces, [5]Piece{
		{},
		{},
		{},
		{},
		{},
	})
	if err != nil {
		return err
	}
	err = checkPieces(b.Tiles[initalLocation[0]+(locationModifier[0]/2)][initalLocation[1]+(locationModifier[1]/2)].Pieces, [5]Piece{
		{Color: RED, Exists: true},
		{Color: RED, Exists: true},
		{},
		{},
		{},
	})
	if err != nil {
		return err
	}
	return nil
}
func TestMoveUp(t *testing.T) {
	err := MoveDirectionTest([2]int{3, 3}, UP)
	if err != nil {
		t.Error(err)
	}
}
func TestMoveDown(t *testing.T) {
	err := MoveDirectionTest([2]int{3, 3}, DOWN)
	if err != nil {
		t.Error(err)
	}
}
func TestMoveLeft(t *testing.T) {
	err := MoveDirectionTest([2]int{3, 3}, LEFT)
	if err != nil {
		t.Error(err)
	}
}
func TestMoveRight(t *testing.T) {
	err := MoveDirectionTest([2]int{3, 3}, RIGHT)
	if err != nil {
		t.Error(err)
	}
}

func TestMoveMultipleDirections(t *testing.T) {
}
