package main

import (
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

var BigWorld = makeWorld(0, 0)
var BigTurn = 0
var Pause = "Continue"
var Quit = "No"

func makeWorld(height int, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}
func worker(world [][]byte, out chan<- [][]byte, ImageHeight int, ImageWidth int) {
	//print("worker")
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
		}
	}
	out <- newWorld
}

type GameOfLife struct{}

func (s *GameOfLife) Alive(req stubs.AliveRequest, res *stubs.AliveResponse) (err error) {
	//fmt.Println("Test 1")
	if Pause == "Pause" {
		for Pause == "Pause" {
			time.Sleep(1 * time.Second)
		}
	}
	aliveCount := 0
	world := BigWorld
	turn := BigTurn
	for i := 0; i < req.ImageHeight; i++ {
		for j := 0; j < req.ImageWidth; j++ {
			if world[i][j] == 255 {
				aliveCount += 1
			}
		}
	}
	res.Turn = turn
	res.World = world
	res.AliveCellsCount = aliveCount
	//fmt.Println("Test 2")
	return
}

func (s *GameOfLife) Key(req stubs.KeyRequest, res *stubs.KeyResponse) (err error) {
	//fmt.Println("Test 1")
	world := BigWorld
	turn := BigTurn
	//if req.Key == 's' {
	//	res.World = world
	//	res.Turn = turn
	//}
	if req.Key == 'p' {
		if Pause == "Continue" {
			Pause = "Pause"
		} else {
			if Pause == "Pause" {
				Pause = "Continue"
			}
		}
		res.Pause = Pause
	}
	if req.Key == 'k' {
		Quit = "Yes"
		res.Pause = "Quit"
	}
	if req.Key == 'k' {
		Quit = "Yes"
		res.Pause = "Quit"
		BigWorld = makeWorld(0, 0)
		BigTurn = 0
	}
	res.Turn = turn
	res.World = world
	//fmt.Println("Test 2")
	return
}

func (s *GameOfLife) GoL(req stubs.Request, res *stubs.Response) (err error) {
	Pause = "Continue"
	Quit = "No"
	BigWorld = req.World
	world := req.World
	BigTurn = 0
	aliveCells := []util.Cell{}
	for turn := 0; turn < req.Turn; turn++ {
		if Pause == "Pause" {
			for Pause == "Pause" {
				time.Sleep(1 * time.Second)
			}
		}
		if Quit == "Yes" {
			res.World = BigWorld
			for i := 0; i < req.ImageHeight; i++ {
				for j := 0; j < req.ImageWidth; j++ {
					if BigWorld[i][j] == 255 {
						newCell := []util.Cell{{j, i}}
						aliveCells = append(aliveCells, newCell...)
					}
				}
			}
			res.AliveCells = aliveCells
			return
		}
		BigTurn = turn
		out := make(chan [][]byte)
		go worker(world, out, req.ImageHeight, req.ImageWidth)
		newWorld := makeWorld(0, 0) // Rebuilds world from sections
		section := <-out
		newWorld = append(newWorld, section...)
		world = newWorld
		BigWorld = world
	}
	res.World = BigWorld
	for i := 0; i < req.ImageHeight; i++ {
		for j := 0; j < req.ImageWidth; j++ {
			if BigWorld[i][j] == 255 {
				newCell := []util.Cell{{j, i}}
				aliveCells = append(aliveCells, newCell...)
			}
		}
	}
	res.AliveCells = aliveCells

	//c.events <- TurnComplete{Turn}
	return
}
func main() {
	//pAddr := flag.String("port", "8030", "Port to listen on")
	//flag.Parse()
	pAddr := "8030"
	//pAddr2 := "8031"
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&GameOfLife{})
	listener, _ := net.Listen("tcp", ":"+pAddr)
	//fmt.Println("Test 4")
	defer listener.Close()
	rpc.Accept(listener)
	//listener2, _ := net.Listen("tcp", ":"+pAddr2)
	//defer listener2.Close()
	//rpc.Accept(listener2)
}
