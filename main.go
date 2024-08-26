package main

import (
	"os"
	"rockwall/discover"
	"rockwall/listener"
	"rockwall/proto"
	"strings"
)

func init() {
	if len(os.Args) != 2 {
		panic("len args != 2")
	}
}

func main() {
	NewNode(os.Args[1]).Run(listener.StartListener, discover.StartDiscover)
}

func NewNode(address string) *proto.Node {
	splited := strings.Split(address, ":")
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
