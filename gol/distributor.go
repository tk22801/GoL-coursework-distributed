package gol

import (
	"fmt"
	"net/rpc"
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
func makeCall(client *rpc.Client, message string) {
	request := stubs.Request{Message: message}
	response := new(stubs.Response)
	client.Call(stubs.ReverseHandler, request, response)
	fmt.Println("Responded: " + response.Message)
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	filename := fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)
	c.ioCommand <- ioInput
	c.ioFilename <- filename
	//Events := make(chan Event)
	//Run(p, c.events, keyPresses)
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

	// TODO: Report the final state using FinalTurnCompleteEvent.

	aliveCells := []util.Cell{}
	for i := 0; i < p.ImageHeight; i++ {
		for j := 0; j < p.ImageWidth; j++ {
			if world[i][j] == 255 {
				newCell := []util.Cell{{j, i}}
				aliveCells = append(aliveCells, newCell...)
			}
		}
	}
	//fmt.Println(aliveCells)
	c.events <- FinalTurnComplete{
		CompletedTurns: turn, Alive: aliveCells,
	}
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
