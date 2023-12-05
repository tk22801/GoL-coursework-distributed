package stubs

import (
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
)

type Worker struct{}

func makeWorld(height int, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

func (s *Worker) Worker(req stubs.WorkerRequest, res *stubs.WorkerResponse) (err error) {
	newWorld := makeWorld(req.ImageHeight, req.ImageWidth)
	for x := 0; x < req.ImageWidth; x++ {
		for y := 0; y < req.ImageHeight; y++ {
			numNeighbours := 0
			xBack := x - 1
			xForward := x + 1
			yUp := y - 1
			yDown := y + 1
			if x == 0 {
				xBack = req.ImageWidth - 1
			}
			if x == req.ImageWidth-1 {
				xForward = 0
			}
			if y == 0 {
				yUp = req.ImageHeight - 1
			}
			if y == req.ImageHeight-1 {
				yDown = 0
			}
			if req.World[xBack][y] == 255 { //Horizontal
				numNeighbours += 1
			}
			if req.World[xForward][y] == 255 {
				numNeighbours += 1
			}
			if req.World[x][yUp] == 255 { //Vertical
				numNeighbours += 1
			}
			if req.World[x][yDown] == 255 {
				numNeighbours += 1
			}
			if req.World[xBack][yDown] == 255 { //Diagonal
				numNeighbours += 1
			}
			if req.World[xForward][yUp] == 255 {
				numNeighbours += 1
			}
			if req.World[xBack][yUp] == 255 {
				numNeighbours += 1
			}
			if req.World[xForward][yDown] == 255 {
				numNeighbours += 1
			}
			if numNeighbours == 2 && req.World[x][y] == 255 || numNeighbours == 3 {
				newWorld[x][y] = 255
			} else {
				newWorld[x][y] = 0
			}
		}
	}
	res.World = newWorld
	return
}

func main() {
	pAddr := "8031"
	rand.Seed(time.Now().UnixNano())
	err := rpc.Register(&Worker{})
	if err != nil {
		return
	}
	listener, _ := net.Listen("tcp", ":"+pAddr)
	//fmt.Println("Test 4")
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)
	rpc.Accept(listener)
}
