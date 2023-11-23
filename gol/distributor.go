package gol

import (
	"flag"
	"fmt"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

func makeWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}
func makeCall(client *rpc.Client, world [][]byte, turn int, height int, width int, c distributorChannels) {
	request := stubs.Request{World: world, Turn: turn, ImageHeight: height, ImageWidth: width}
	response := new(stubs.Response)
	client.Call(stubs.GoLWorker, request, response)
	fmt.Println(response.World)
	c.events <- FinalTurnComplete{CompletedTurns: turn, Alive: response.AliveCells}
	c.ioCommand <- ioOutput
	filename := fmt.Sprintf("%dx%dx%d", width, height, turn)
	c.ioFilename <- filename
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			c.ioOutput <- world[i][j]
		}
	}
	c.events <- ImageOutputComplete{turn, filename}
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	filename := fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)
	c.ioCommand <- ioInput
	c.ioFilename <- filename
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	turn := 0
	//newWorld := makeWorld(0, 0)
	world := makeWorld(p.ImageHeight, p.ImageWidth)
	for x := 0; x < p.ImageHeight; x++ {
		for y := 0; y < p.ImageWidth; y++ {
			val := <-c.ioInput
			world[x][y] = val
			if val == 255 {
				c.events <- CellFlipped{turn, util.Cell{x, y}}
			}
		}
	}
	// TODO: Execute all turns of the Game of Life.
	fmt.Println("Called")
	makeCall(client, world, p.Turns, p.ImageHeight, p.ImageWidth, c)
	// TODO: Report the final state using FinalTurnCompleteEvent.
	defer client.Close()
	//fmt.Println(aliveCells)
	c.ioCommand <- ioOutput
	filename = fmt.Sprintf("%dx%dx%d", p.ImageWidth, p.ImageHeight, turn)
	c.ioFilename <- filename
	for i := 0; i < p.ImageHeight; i++ {
		for j := 0; j < p.ImageWidth; j++ {
			c.ioOutput <- world[i][j]
		}
	}
	c.events <- ImageOutputComplete{turn, filename}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
