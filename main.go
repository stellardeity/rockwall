package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"rockwall/discover"
	"rockwall/listener"
	"rockwall/proto"
	"strings"
)

type InitParams struct {
	Name    *string
	Address *string
}

var initParams InitParams

func init() {
	currentUser, _ := user.Current()
	hostName, _ := os.Hostname()

	initParams = InitParams{
		Name:    flag.String("name", currentUser.Username+"@"+hostName, "your name"),
		Address: flag.String("address", "192.168.0.100:8080", "the address of the peer to connect to"),
	}

	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	NewNode().Run(listener.StartListener, discover.StartDiscover)
}

func NewNode() *proto.Node {
	splited := strings.Split(*initParams.Address, ":")
	if len(splited) != 2 {
		return nil
	}
	return &proto.Node{
		Connections: make(map[string]bool),
		Address: proto.Address{
			IPv4: splited[0],
			Port: ":" + splited[1],
		},
	}
}
