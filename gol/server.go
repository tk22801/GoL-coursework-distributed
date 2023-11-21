package gol

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

func worker() {

}
func GOL() {
	for Turn := 0; Turn < p.Turns; Turn++ {
		out := make(chan [][]byte)
		go worker(p, c, out[i], world, newWorld, workerHeight, i, Turn)
		finalWorld := makeWorld(0, 0) // Rebuilds world from sections
		section := <-out
		world = append(finalWorld, section...)
	}
	c.events <- TurnComplete{Turn}
}
func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&SecretStringOperations{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
