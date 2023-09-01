package game

import "errors"

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

type Color int

const (
	RED Color = iota
	GREEN
)

type Piece struct {
	Color  Color
	Exists bool
}

type Tile struct {
	useable bool
	Pieces  [5]Piece
}

type Board struct {
	Tiles     [8][8]Tile
	ReservesR int
	ReservesG int
}

func NewBoard() Board {
	b := Board{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			b.SetTile(i, j, Tile{useable: true})
		}
	}
	// Corner triangles.
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
		b.SetTile(unusableTiles[i][0], unusableTiles[i][1], Tile{useable: false})
	}
	return b
}

var (
	ErrTileOutOfBounds = errors.New("Tile out of bounds")
)

func (b *Board) GetTile(x int, y int) (*Tile, error) {
	if x < 0 || x > 7 || y < 0 || y > 7 {
		return nil, ErrTileOutOfBounds
	}
	return &b.Tiles[x][y], nil
}

func (b *Board) GetReserves(color Color) *int {
	if color == RED {
		return &b.ReservesR
	}
	return &b.ReservesG
}

func (b *Board) SetReserves(color Color, amount int) {
	if color == RED {
		b.ReservesR = amount
	} else {
		b.ReservesG = amount
	}
}

func (b *Board) AddToReserves(color Color, amount int) {
	if color == RED {
		b.ReservesR += amount
	} else {
		b.ReservesG += amount
	}
}

func (b *Board) SetTile(x int, y int, tile Tile) {
	b.Tiles[x][y] = tile
}

var (
	ErrPieceNonExistent = errors.New("Cannot add a Piece that does not exist")
)

func (b *Board) AddPiece(x int, y int, piece Piece, playerColor Color) (err error) {
	if !piece.Exists {
		return ErrPieceNonExistent
	}
	tile, err := b.GetTile(x, y)
	if err != nil {
		return err
	}
	for i := 0; i < len(tile.Pieces); i++ {
		if !tile.Pieces[i].Exists {
			tile.Pieces[i] = piece
			return nil
		}
	}
	// No space for piece, shift down to make space.
	if tile.Pieces[0].Color == RED && playerColor == RED {
		b.AddToReserves(RED, 1)
	}
	if tile.Pieces[0].Color == GREEN && playerColor == GREEN {
		b.AddToReserves(GREEN, 1)
	}
	tile.Pieces[0] = tile.Pieces[1]
	tile.Pieces[1] = tile.Pieces[2]
	tile.Pieces[2] = tile.Pieces[3]
	tile.Pieces[3] = tile.Pieces[4]
	tile.Pieces[4] = piece
	return nil
}

func (b *Board) AddFromReserves(color Color, x int, y int) (err error) {
	if color == RED {
		b.ReservesR -= 1
	} else {
		b.ReservesG -= 1
	}
	err = b.AddPiece(x, y, Piece{Color: color, Exists: true}, color) // The player is always the same color as the reserved piece.
	return err
}

func (b *Board) moveDirection(x int, y int, tile *Tile, moveModifier [2]int, piecesToMove int, playerColor Color) (err error) {
	destTile, err := b.GetTile(x+moveModifier[0], y+moveModifier[1])
	if err != nil {
		return err
	}
	if !destTile.useable {
		return ErrTileDestinationUnusable
	}
	piecesMoved := 0
	for i := len(tile.Pieces) - 1; i >= 0; i-- {
		if piecesMoved == piecesToMove {
			break
		}
		if !tile.Pieces[i].Exists {
			continue
		}
		err = b.AddPiece(x+moveModifier[0], y+moveModifier[1], tile.Pieces[i], playerColor)
		if err != nil {
			return err
		}

		tile.Pieces[i] = Piece{Color: 0, Exists: false}

		piecesMoved += 1
	}
	return nil
}

var (
	ErrTileSourceNonUsable     = errors.New("Cannot move from unusable tile")
	ErrTileDestinationUnusable = errors.New("Cannot move to unusable tile")
	ErrNoPieceToMove           = errors.New("No piece to move")
	ErrMustMoveAtLeastOnePiece = errors.New("You must move at least one piece")
	ErrWrongColor              = errors.New("You cannot move your opponent's piece")
)

func (b *Board) Move(x int, y int, piecesToMove int, direction Direction, playerColor Color) (err error) {
	tile, err := b.GetTile(x, y)
	if err != nil {
		return err
	}
	if !tile.useable {
		return ErrTileSourceNonUsable
	}
	allPiecesCheck := len(tile.Pieces)
	for i := len(tile.Pieces) - 1; i >= 0; i-- {
		if tile.Pieces[i].Exists {
			allPiecesCheck -= 1
		}
	}
	if allPiecesCheck == 5 {
		return ErrNoPieceToMove
	}

	for i := len(tile.Pieces) - 1; i >= 0; i-- {
		if tile.Pieces[i].Exists {
			if tile.Pieces[i].Color != playerColor {
				return ErrWrongColor
			}
		}
	}

	if piecesToMove == 0 {
		return ErrMustMoveAtLeastOnePiece
	}

	switch direction {
	case UP:
		err = b.moveDirection(x, y, tile, [2]int{0, -1 * piecesToMove}, piecesToMove, playerColor)
		if err != nil {
			return err
		}
	case DOWN:
		err = b.moveDirection(x, y, tile, [2]int{0, 1 * piecesToMove}, piecesToMove, playerColor)
		if err != nil {
			return err
		}
	case LEFT:
		err = b.moveDirection(x, y, tile, [2]int{-1 * piecesToMove, 0}, piecesToMove, playerColor)
		if err != nil {
			return err
		}
	case RIGHT:
		err = b.moveDirection(x, y, tile, [2]int{1 * piecesToMove, 0}, piecesToMove, playerColor)
		if err != nil {
			return err
		}
	}
	return nil
}
