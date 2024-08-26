package listener

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"rockwall/proto"
	"strings"
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
		handleHttp(readWriter, conn)
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

func handleHttp(rw *bufio.ReadWriter, conn net.Conn) {
	request, err := http.ReadRequest(rw.Reader)

	if err != nil {
		log.Printf("Read request ERROR: %s", err)
		return
	}

	response := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	s := conn.RemoteAddr().String()[0:3]
	// TODO: сравнение среза со строкой
	if !strings.EqualFold(s, "127") && !strings.EqualFold(s, "[::") {
		response.Body = ioutil.NopCloser(strings.NewReader("Banner"))
	} else {
		if path.Clean(request.URL.Path) == "/ws" {
			handleWs(NewMyWriter(conn), request)
			return
		} else {
			processRequest(request, &response)
			//fileServer := http.FileServer(http.Dir("./front/build/"))
			//fileServer.ServeHTTP(NewMyWriter(conn), request)
		}
	}

	err = response.Write(rw)
	if err != nil {
		log.Printf("Write response ERROR: %s", err)
		return
	}

	err = rw.Writer.Flush()
	if err != nil {
		log.Printf("Flush response ERROR: %s", err)
		return
	}
}
