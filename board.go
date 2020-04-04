package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type board struct {
	turn         int
	width        int
	height       int
	imageSize    int
	playCount    int
	pieces       [][]int
	renderer     *sdl.Renderer
	player1Image *sdl.Texture
	player2Image *sdl.Texture
}

const (
	GRID_EMPTY = 0
	PLAYER_1   = 1
	PLAYER_2   = 2
	MAX_PLAYS  = 9
)

func newBoard(sqrSize int, renderer *sdl.Renderer, images []*sdl.Texture) *board {
	return &board{
		pieces: [][]int{
			make([]int, 3),
			make([]int, 3),
			make([]int, 3),
		},
		renderer:     renderer,
		turn:         PLAYER_1,
		player1Image: images[0],
		player2Image: images[1],
		width:        sqrSize,
		height:       sqrSize,
		imageSize:    sqrSize / 4,
	}
}

func (b *board) getPiecePosition(row, col int) (int32, int32) {
	square := b.height / 3
	padding := square - b.imageSize
	x := col*square + padding/2
	y := row*square + padding/2
	return int32(x), int32(y)
}

func (b *board) renderPiece(row int, col int, img *sdl.Texture) {
	x, y := b.getPiecePosition(row, col)
	b.renderer.Copy(img, nil, &sdl.Rect{x, y, int32(b.imageSize), int32(b.imageSize)})
}

func (b *board) renderPieces() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch b.pieces[i][j] {
			case PLAYER_1:
				b.renderPiece(i, j, b.player1Image)
			case PLAYER_2:
				b.renderPiece(i, j, b.player2Image)
			}
		}
	}
}

func (b *board) drawBoard() {
	b.renderer.SetDrawColor(255, 255, 255, 255)
	b.renderer.Clear()

	b.renderer.SetDrawColor(0, 0, 0, 255)
	// vertical lines
	b.renderer.DrawLine(screenWidth/3, 0, screenWidth/3, screenHeight)
	b.renderer.DrawLine(screenWidth/3*2, 0, screenWidth/3*2, screenHeight)

	// horizontal lines
	b.renderer.DrawLine(0, screenHeight/3, screenWidth, screenHeight/3)
	b.renderer.DrawLine(0, screenHeight/3*2, screenWidth, screenHeight/3*2)

	b.renderPieces()

	b.renderer.Present()
}

func (b *board) handleMouseButtonEvent(e *sdl.MouseButtonEvent) (bool, int) {
	if e.Type != sdl.MOUSEBUTTONDOWN || e.Button != sdl.BUTTON_LEFT {
		return false, 0
	}

	row := e.Y / (screenHeight / 3)
	col := e.X / (screenWidth / 3)

	if !b.validateMove(row, col) {
		return false, 0
	}

	b.setMove(row, col)
	b.playCount++

	gameOver, winner := b.checkStatus()

	if b.playCount == MAX_PLAYS && !gameOver {
		return true, 0
	}

	if !gameOver {
		b.changePlayer()
	}

	return gameOver, winner
}

func (b *board) validateMove(row, col int32) bool {
	return b.pieces[row][col] == GRID_EMPTY
}

func (b *board) setMove(row, col int32) {
	b.pieces[row][col] = b.turn
}

func (b *board) changePlayer() {
	if b.turn == PLAYER_1 {
		b.turn = PLAYER_2
		return
	}

	b.turn = PLAYER_1
}

func (b *board) checkStatus() (bool, int) {
	gameOver := b.checkRows() || b.checkColumns() || b.checkDiagonals()

	if gameOver {
		return gameOver, b.turn
	}

	return false, 0
}

func (b *board) checkRows() bool {
	for i := 0; i < 3; i++ {
		if b.pieces[i][0] == b.turn && b.pieces[i][0] == b.pieces[i][1] && b.pieces[i][0] == b.pieces[i][2] {
			return true
		}
	}

	return false
}

func (b *board) checkColumns() bool {
	for i := 0; i < 3; i++ {
		if b.pieces[0][i] == b.turn && b.pieces[0][i] == b.pieces[1][i] && b.pieces[0][i] == b.pieces[2][i] {
			return true
		}
	}

	return false
}

func (b *board) checkDiagonals() bool {
	if b.pieces[0][0] == b.turn && b.pieces[0][0] == b.pieces[1][1] && b.pieces[0][0] == b.pieces[2][2] {
		return true
	}
	if b.pieces[0][2] == b.turn && b.pieces[0][2] == b.pieces[1][1] && b.pieces[0][0] == b.pieces[2][0] {
		return true
	}

	return false
}
