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
		go onConnection(node, conn)
	}
}

func onConnection(node *proto.Node, conn net.Conn) {
	defer conn.Close()
	var (
		buffer  = make([]byte, 512)
		message string
		pack    proto.Package
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		message += string(buffer[:length])
	}
	err := json.Unmarshal([]byte(message), &pack)
	if err != nil {
		return
	}
	node.ConnectTo([]string{pack.From})

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readWriter := bufio.NewReadWriter(reader, writer)

	buf, err := readWriter.Peek(4)
	if err != nil {
		if err != io.EOF {
			log.Printf("Read peak ERROR: %s", err)
		}
		return
	}

	if ItIsHttp(buf) {
		fmt.Println("IS HTTP")
	} else {
		fmt.Println(pack.Data)
	}
}
