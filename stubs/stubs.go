package stubs

import (
	"uk.ac.bris.cs/gameoflife/util"
)

var GoLWorker = "GameOfLife.GoL"
var AliveWorker = "GameOfLife.Alive"

type Response struct {
	World      [][]byte
	AliveCells []util.Cell
}

type Request struct {
	World       [][]byte
	Turn        int
	ImageHeight int
	ImageWidth  int
}
type AliveRequest struct {
	ImageHeight int
	ImageWidth  int
}
type AliveResponse struct {
	World           [][]byte
	Turn            int
	AliveCellsCount int
}
