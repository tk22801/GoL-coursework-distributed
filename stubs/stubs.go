package stubs

import (
	"uk.ac.bris.cs/gameoflife/util"
)

var GoLWorker = "GameOfLife.GoL"
var AliveWorker = "GameOfLife.Alive"
var KeyPresses = "GameOfLife.Key"
var Worker = "Worker.Worker"

type Response struct {
	World      [][]byte
	AliveCells []util.Cell
	Turn       int
}

type Request struct {
	World         [][]byte
	Turn          int
	ImageHeight   int
	ImageWidth    int
	WorkerAddress string
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
type KeyRequest struct {
	Key rune
}

type KeyResponse struct {
	World [][]byte
	Turn  int
	Pause string
}
type WorkerRequest struct {
	World       [][]byte
	ImageHeight int
	ImageWidth  int
	Quit        string
}

type WorkerResponse struct {
	World [][]byte
}
