package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
)

func makeWorld(height int, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}
func worker(world [][]byte, out chan<- [][]byte, ImageHeight int, ImageWidth int) {
	newWorld := makeWorld(ImageHeight, ImageWidth)
	for x := 0; x < ImageWidth; x++ {
		for y := 0; y < ImageHeight; y++ {
			numNeighbours := 0
			xBack := x - 1
			xForward := x + 1
			yUp := y - 1
			yDown := y + 1
			if x == 0 {
				xBack = ImageWidth - 1
			}
			if x == ImageWidth-1 {
				xForward = 0
			}
			if y == 0 {
				yUp = ImageHeight - 1
			}
			if y == ImageHeight-1 {
				yDown = 0
			}
			if world[xBack][y] == 255 { //Horizontal
				numNeighbours += 1
			}
			if world[xForward][y] == 255 {
				numNeighbours += 1
			}
			if world[x][yUp] == 255 { //Vertical
				numNeighbours += 1
			}
			if world[x][yDown] == 255 {
				numNeighbours += 1
			}
			if world[xBack][yDown] == 255 { //Diagonal
				numNeighbours += 1
			}
			if world[xForward][yUp] == 255 {
				numNeighbours += 1
			}
			if world[xBack][yUp] == 255 {
				numNeighbours += 1
			}
			if world[xForward][yDown] == 255 {
				numNeighbours += 1
			}
			if numNeighbours == 2 && world[x][y] == 255 || numNeighbours == 3 {
				newWorld[x][y] = 255
			} else {
				newWorld[x][y] = 0
			}
			//if p.Turns == 1 {
			//fmt.Println(x, y, newWorld[x][y])
			//}
		}
	}
	out <- world
}

type GameOfLife struct{}

func (s *GameOfLife) GoL(req stubs.Request, res *stubs.Response) (err error) {
	world := req.World
	for turn := 0; turn < req.Turn; turn++ {
		out := make(chan [][]byte)
		worker(world, out, req.ImageHeight, req.ImageWidth)
		newWorld := makeWorld(0, 0) // Rebuilds world from sections
		section := <-out
		newWorld = append(newWorld, section...)
		world = newWorld
	}
	res.World = world
	//c.events <- TurnComplete{Turn}
	return
}
func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&GameOfLife{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
