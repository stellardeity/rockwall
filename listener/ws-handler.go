package listener

import (
	"encoding/hex"
	"encoding/json"
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

func handleWs(w http.ResponseWriter, r *http.Request, p *proto.Proto) {
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	br := make(chan bool)

	go waitMessageForWs(p, c, br)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %v", err)
			break
		}
		log.Printf("ws read: [%v] %s", mt, message)

		decodedMessage := &proto.WsMessage{}
		err = json.Unmarshal(message, decodedMessage)

		if err != nil {
			log.Printf("error on unmarshal message: %v", err)
			continue
		}

		switch decodedMessage.Cmd {
		case "HELLO":
			{
				myName := p.MyName()
				name := proto.WsMyName{
					WsCmd: proto.WsCmd{
						Cmd: "NAME",
					},
					Name:   myName.Name,
					PubKey: myName.PubKey,
				}
				writeToWs(c, mt, name.ToJson())
			}
		case "PEERS":
			{
				peerList := p.Peers.PeerList()

				peerListJson, err := json.Marshal(peerList)

				if err != nil {
					panic(err)
				}

				writeToWs(c, mt, peerListJson)
			}
		case "MESS":
			{
				hexPubKey, err := hex.DecodeString(decodedMessage.To)
				if err != nil {
					log.Printf("decode error: %s", err)
					continue
				}
				peer, found := p.Peers.Get(string(hexPubKey))
				if found {
					writeToWs(c, mt, message)
					p.SendMessage(peer, decodedMessage.Content)
				}

			}
		}
	}

	br <- true
}

func waitMessageForWs(p *proto.Proto, c *websocket.Conn, br chan bool) {
	for {
		select {
		case envelope := <-p.Broker:
			{
				log.Printf("New message: %s", envelope.Cmd)
				if string(envelope.Cmd) == "MESS" {

					wsCmd := proto.WsMessage{
						WsCmd: proto.WsCmd{
							Cmd: "MESS",
						},
						From:    hex.EncodeToString(envelope.From),
						To:      hex.EncodeToString(envelope.To),
						Content: string(envelope.Content),
					}

					wsCmdBytes, err := json.Marshal(wsCmd)

					if err != nil {
						panic(err)
					}

					writeToWs(c, 1, wsCmdBytes)
				}
			}
		case _ = <-br:
			{
				log.Printf("ws is broken")
				return
			}
		}
	}
}

func writeToWs(c *websocket.Conn, mt int, message []byte) {
	err := c.WriteMessage(mt, message)
	if err != nil {
		log.Printf("ws write error: %s", err)
	}
}
