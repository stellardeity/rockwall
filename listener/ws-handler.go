package listener

import (
	"log"
	"net/http"
	"rockwall/proto"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWs(w http.ResponseWriter, r *http.Request, node *proto.Node) {
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %v", err)
			break
		}
		log.Printf("ws read: [%v] %s", mt, message)
		node.SendMessageToAll(string(message))
		writeToWs(c, mt, message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Printf("ws write error: %v", err)
			break
		}
	}
}

func writeToWs(c *websocket.Conn, mt int, message []byte) {
	err := c.WriteMessage(mt, message)
	if err != nil {
		log.Printf("ws write error: %s", err)
	}
}
