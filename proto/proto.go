package proto

import (
	"encoding/json"
	"fmt"
	"net"
)

type Node struct {
	Connections map[string]bool
	Address     Address
}

type Address struct {
	IPv4 string
	Port string
}

type Package struct {
	To   string
	From string
	Data string
}

func (node *Node) Run(handleServer func(*Node), handleClient func(*Node)) {
	go handleServer(node)
	handleClient(node)
}

func (node *Node) PrintNetwork() {
	for addr := range node.Connections {
		fmt.Println("|", addr)
	}
}

func (node *Node) ConnectTo(addresses []string) {
	for _, addr := range addresses {
		node.Connections[addr] = true
	}
}

func (node *Node) SendMessageToAll(message string) {
	var newPack = &Package{
		From: node.Address.IPv4 + node.Address.Port,
		Data: message,
	}
	for addr := range node.Connections {
		newPack.To = addr
		node.Send(newPack)
	}
}

func (node *Node) Send(pack *Package) {
	conn, err := net.Dial("tcp", pack.To)
	if err != nil {
		delete(node.Connections, pack.To)
		return
	}
	defer conn.Close()
	jsonPack, _ := json.Marshal(*pack)
	conn.Write(jsonPack)
}
