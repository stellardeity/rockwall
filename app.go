package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"sync"

	"rockwall/discover"
	"rockwall/listener"
	"rockwall/proto"
)

type InitParams struct {
	Name      *string
	Port      *int
	PeersFile *string
}

var initParams InitParams

func init() {
	currentUser, _ := user.Current()
	hostName, _ := os.Hostname()

	initParams = InitParams{
		Name:      flag.String("name", currentUser.Username+"@"+hostName, "your name"),
		Port:      flag.Int("port", 35035, "port that have to listen"),
		PeersFile: flag.String("peers", "peers.txt", "path to file with peer addresses on each line"),
	}

	flag.Parse()
}

func retro() func(a int, b int) int {
	return func(a int, b int) int {
		return a + b
	}
}

func main() {
	fff := retro()

	i := fff(1, 2)
	log.Printf("result = %v", i)

	p := proto.NewProto(*initParams.Name, *initParams.Port)

	var wg sync.WaitGroup
	wg.Add(2)
	go discover.StartDiscover(p, *initParams.PeersFile)
	go listener.StartListener(p, *initParams.Port)
	wg.Wait()
}
