package listener

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"rockwall/proto"
)

var itHttp = map[string]bool{
	"GET ": true,
	"HEAD": true,
	"POST": true,
	"PUT ": true,
	"DELE": true,
	"CONN": true,
	"OPTI": true,
	"TRAC": true,
	"PATC": true,
}

func ItIsHttp(ba []byte) bool {
	return itHttp[string(ba)]
}

func StartListener(node *proto.Node) {
	listen, err := net.Listen("tcp", "0.0.0.0"+node.Address.Port)
	if err != nil {
		panic("listen error")
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			break
		}
		go handleConnection(node, conn)
	}
}

func handleConnection(node *proto.Node, conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	readWriter := bufio.NewReadWriter(reader, writer)

	data, err := readWriter.Peek(4)
	if err != nil {
		log.Printf("error: %s", err)
		return
	}

	if ItIsHttp(data) {
		log.Printf("HTTP-request")
		return
	}

	message, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("data reading error: %s", err)
		return
	}

	var pack proto.Package
	err = json.Unmarshal(message, &pack)
	if err != nil {
		log.Printf("JSON deserialization error: %s", err)
		return
	}

	node.ConnectTo([]string{pack.From})

	fmt.Println(pack.Data)
}
