package stubs

import "uk.ac.bris.cs/gameoflife/util"

var GoLWorker = "GameOfLife.GoL"

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
