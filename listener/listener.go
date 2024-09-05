package listener

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
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

func StartListener(proto *proto.Proto, port int) {
	if port <= 0 || port > 65535 {
		port = 35035
	}

	service := fmt.Sprintf("0.0.0.0:%v", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		log.Printf("ResolveTCPAddr: %s", err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("ListenTCP: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n\tService start on %s\n\n", tcpAddr.String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go onConnection(conn, proto)
	}

}

func onConnection(conn net.Conn, p *proto.Proto) {
	defer func() {
		conn.Close()
	}()

	log.Printf("New connection from: %v", conn.RemoteAddr().String())

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
		handleHttp(readWriter, conn, p)
	} else {
		peer := proto.NewPeer(conn)
		p.HandleProto(readWriter, peer)
	}
}

func handleHttp(rw *bufio.ReadWriter, conn net.Conn, p *proto.Proto) {
	request, err := http.ReadRequest(rw.Reader)

	if err != nil {
		log.Printf("Read request ERROR: %s", err)
		return
	}

	response := http.Response{
		StatusCode: 200,
	}

	s := conn.RemoteAddr().String()[0:3]

	if !strings.EqualFold(s, "127") && !strings.EqualFold(s, "[::") {
		response.Body = ioutil.NopCloser(strings.NewReader("Rockwall"))
	} else {
		if path.Clean(request.URL.Path) == "/ws" {
			handleWs(NewMyWriter(conn), request, p)
			return
		} else {
			processRequest(request, &response)
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
