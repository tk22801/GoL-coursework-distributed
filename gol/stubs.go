package gol

import "uk.ac.bris.cs/gameoflife/util"

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
