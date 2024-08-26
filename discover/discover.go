package discover

import (
	"bufio"
	"os"
	"rockwall/proto"
	"strings"
)

func StartDiscover(node *proto.Node) {
	for {
		message := InputString()
		splited := strings.Split(message, " ")
		switch splited[0] {
		case "/exit":
			os.Exit(0)
		case "/connect":
			node.ConnectTo(splited[1:])
		case "/network":
			node.PrintNetwork()
		default:
			node.SendMessageToAll(message)
		}
	}
}

func InputString() string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", 1)
}
