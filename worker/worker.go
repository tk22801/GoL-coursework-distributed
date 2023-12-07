package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
)

var Quit = "no"

type Worker struct{}

func makeWorld(height int, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

func (s *Worker) Worker(req stubs.WorkerRequest, res *stubs.WorkerResponse) (err error) {
	if req.Quit == "yes" {
		//allows the distributed element to be of the worker to be shutdown on request
		res.World = req.World
		Quit = "yes"
	}
	//implements GameOfLife to the provided world and outputs a new world
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
	pAddr := "8040"
	rand.Seed(time.Now().UnixNano())
	err := rpc.Register(&Worker{})
	if err != nil {
		return
	}
	listener, _ := net.Listen("tcp", ":"+pAddr)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)
	go func() {
		for {
			if Quit == "yes" {
				fmt.Println("Quitting")
				err := listener.Close()
				if err != nil {
					return
				}
				os.Exit(0)
			}
		}
	}()
	rpc.Accept(listener)
}
