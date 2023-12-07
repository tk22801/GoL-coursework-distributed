package gol

import (
	"fmt"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
	keyPresses <-chan rune
}

//Makes a World of the specified height and width
func makeWorld(height, width int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

//Makes the call to the server calling the GoL function
func makeCall(client *rpc.Client, world [][]byte, turn int, height int, width int, c distributorChannels, workerServer string) int {
	request := stubs.Request{World: world, Turn: turn, ImageHeight: height, ImageWidth: width, WorkerAddress: workerServer}
	response := new(stubs.Response)
	err := client.Call(stubs.GoLWorker, request, response)
	if err != nil {
		return 0
	}
	//Sends finalTurnComplete event down the events channel, as well as outputs a pgm image of the final world
	c.events <- FinalTurnComplete{CompletedTurns: response.Turn, Alive: response.AliveCells}
	c.ioCommand <- ioOutput
	filename := fmt.Sprintf("%dx%dx%d", width, height, response.Turn)
	c.ioFilename <- filename
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			c.ioOutput <- world[i][j]
		}
	}
	c.events <- ImageOutputComplete{response.Turn, filename}
	return response.Turn
}

//Makes a call to the Alive function
func makeAliveCall(client *rpc.Client, height int, width int, c distributorChannels) {
	request := stubs.AliveRequest{ImageHeight: height, ImageWidth: width}
	response := new(stubs.AliveResponse)
	err := client.Call(stubs.AliveWorker, request, response)
	if err != nil {
		return
	}
	c.events <- AliveCellsCount{response.Turn, response.AliveCellsCount}
}

//Makes a call to the KeyPresses function
func makeKeyCall(client *rpc.Client, key rune, height int, width int, c distributorChannels) {
	request := stubs.KeyRequest{Key: key}
	response := new(stubs.KeyResponse)
	err := client.Call(stubs.KeyPresses, request, response)
	if err != nil {
		return
	}
	if key == 's' || key == 'k' {
		//outputs world as a pgm file
		c.ioCommand <- ioOutput
		filename := fmt.Sprintf("%dx%dx%d", width, height, response.Turn)
		c.ioFilename <- filename
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				c.ioOutput <- response.World[i][j]
			}
		}
		c.events <- ImageOutputComplete{response.Turn, filename}
	}
	if key == 'p' {
		//prints to local controller that the program is paused
		if response.Pause == "Pause" {
			c.events <- StateChange{response.Turn, Paused}
			fmt.Println("Paused at turn ", response.Turn)
		}
		if response.Pause == "Continue" {
			c.events <- StateChange{response.Turn, Executing}
			fmt.Println("Continuing")
		}
	}
	if key == 'q' {
		// Make sure that the Io has finished any output before exiting.
		c.ioCommand <- ioCheckIdle
		<-c.ioIdle
		c.events <- StateChange{response.Turn, Quitting}

	}
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	//requests file with filename from the input
	filename := fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)
	c.ioCommand <- ioInput
	c.ioFilename <- filename
	//asks the user for the Ip Address and gate of the AWS instance since the IP address changes when the instance is restarted
	workerServer := ""
	fmt.Println("Ip address of AWS Node and gate(IpAddress:gate):")
	fmt.Scan(&workerServer)
	server := "127.0.0.1:8030"
	//workerServer := "3.85.89.162:8040"
	//server := "3.85.6.20:8030"
	//workerServer := "3.85.6.20:8040"
	client, _ := rpc.Dial("tcp", server)
	defer func(client *rpc.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)
	//Make world and fills it with the input values
	world := makeWorld(p.ImageHeight, p.ImageWidth)
	for x := 0; x < p.ImageHeight; x++ {
		for y := 0; y < p.ImageWidth; y++ {
			val := <-c.ioInput
			world[x][y] = val
		}
	}
	ticker := time.NewTicker(2 * time.Second)
	//Every 2 seconds it starts the makeAliveCall
	go func() {
		for range ticker.C {
			makeAliveCall(client, p.ImageHeight, p.ImageWidth, c)
		}
	}()
	//Continuously reads key presses and if it is s,q,k or p it then starts the makeKeyCall function
	go func() {
		for {
			key := <-c.keyPresses
			if key == 's' || key == 'q' || key == 'k' || key == 'p' {
				makeKeyCall(client, key, p.ImageHeight, p.ImageWidth, c)
			}
		}
	}()
	lastTurn := makeCall(client, world, p.Turns, p.ImageHeight, p.ImageWidth, c, workerServer)
	ticker.Stop()
	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	c.events <- StateChange{lastTurn, Quitting}
	close(c.events)
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
}
