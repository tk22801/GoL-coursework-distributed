package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

var BigWorld = makeWorld(0, 0)
var BigTurn = 0
var Pause = "Continue"
var Quit = "No"
var Close = "No"

func makeWorld(height int, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}
func worker(world [][]byte, out chan<- [][]byte, ImageHeight int, ImageWidth int) [][]byte {
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
	return newWorld
}

type GameOfLife struct{}

func (s *GameOfLife) Alive(req stubs.AliveRequest, res *stubs.AliveResponse) (err error) {
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
	res.AliveCellsCount = aliveCount
	return
}

func (s *GameOfLife) Key(req stubs.KeyRequest, res *stubs.KeyResponse) (err error) {
	world := BigWorld
	turn := BigTurn
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
	if req.Key == 'q' || req.Key == 'k' {
		Quit = "Yes"
		res.Pause = "Quit"
	}
	if req.Key == 'k' {
		Close = "yes"
	}
	res.Turn = turn
	res.World = world
	return
}

func (s *GameOfLife) GoL(req stubs.Request, res *stubs.Response) (err error) {
	Pause = "Continue"
	Quit = "No"
	BigWorld = req.World
	world := req.World
	BigTurn = 0
	aliveCells := []util.Cell{}
	//client, _ := rpc.Dial("tcp", req.WorkerAddress)
	for turn := 0; turn < req.Turn; turn++ {
		if Pause == "Pause" {
			for Pause == "Pause" {
				time.Sleep(1 * time.Second)
			}
		}
		//if Close == "yes" {
		//	//request := stubs.WorkerRequest{World: world, ImageHeight: req.ImageHeight, ImageWidth: req.ImageWidth, Quit: "yes"}
		//	//response := new(stubs.WorkerResponse)
		//	//err := client.Call(stubs.Worker, request, response)
		//	//if err != nil {
		//	//	return err
		//	//}
		//}
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
			res.Turn = turn
			return
		}
		BigTurn = turn
		out := make(chan [][]byte)
		world = worker(world, out, req.ImageHeight, req.ImageWidth)

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
	res.Turn = BigTurn
	return
}
func main() {
	//pAddr := flag.String("port", "8030", "Port to listen on")
	//flag.Parse()
	pAddr := "8040"
	rand.Seed(time.Now().UnixNano())
	err := rpc.Register(&GameOfLife{})
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
			if Close == "yes" {
				//time.Sleep(1 * time.Second)
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
