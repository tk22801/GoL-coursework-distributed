package stubs

import (
	"uk.ac.bris.cs/gameoflife/util"
)

var GoLWorker = "GameOfLife.GoL"
var AliveWorker = "GameofLife.Alive"

type Response struct {
	World      [][]byte
	AliveCells []util.Cell
}

type Request struct {
	World       [][]byte
	AcrossWorld chan [][]byte
	AcrossTurn  chan int
	Turn        int
	ImageHeight int
	ImageWidth  int
}
type AliveRequest struct {
	AcrossWorld chan [][]byte
	AcrossTurn  chan int
	Turn        int
	ImageHeight int
	ImageWidth  int
}
type AliveResponse struct {
	World           [][]byte
	Turn            int
	AliveCellsCount int
}
