package stubs

var GoLWorker = "GameOfLife.GoL"

type Response struct {
	World [][]byte
	//AliveCells []util.Cell
}

type Request struct {
	World       [][]byte
	Turn        int
	ImageHeight int
	ImageWidth  int
}