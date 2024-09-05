package discover

import (
	"bufio"
	"log"
	"net"
	"os"
	"rockwall/proto"
)

var peers = make(map[string]string)

func StartDiscover(p *proto.Proto, peersFile string) {
	file, err := os.Open(peersFile)
	if err != nil {
		log.Printf("DISCOVER: Open peers.txt error: %s", err)
		return
	}

	var lastPeers []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastPeers = append(lastPeers, scanner.Text())
	}

	log.Printf("DISCOVER: Start peer discovering. Last seen peers: %v", len(lastPeers))
	for _, peerAddress := range lastPeers {
		go connectToPeer(p, peerAddress)
	}
}

func connectToPeer(p *proto.Proto, peerAddress string) {
	if _, exist := peers[peerAddress]; exist {
		log.Printf("peer %s already exist", peerAddress)
		return
	}
	peers[peerAddress] = peerAddress
	log.Printf("try to connect peer: %s", peerAddress)

	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		log.Printf("Dial ERROR: " + err.Error())
		return
	}

	defer conn.Close()

	peer := handShake(p, conn)

	if peer == nil {
		log.Printf("Fail on handshake")
		return
	}

	p.RegisterPeer(peer)

	p.ListenPeer(peer)

	p.UnregisterPeer(peer)

	delete(peers, peerAddress)
}

func handShake(p *proto.Proto, conn net.Conn) *proto.Peer {
	log.Printf("DISCOVERY: try handshake with %s", conn.RemoteAddr())
	peer := proto.NewPeer(conn)

	p.SendName(peer)

	envelope, err := proto.ReadEnvelope(bufio.NewReader(conn))
	if err != nil {
		log.Printf("Error on read Envelope: %s", err)
		return nil
	}

	if string(envelope.Cmd) == "HAND" {
		if _, found := p.Peers.Get(string(envelope.From)); found {
			log.Printf(" - - - - - - - - - - - - - - - --  -- - - - - Peer (%s) already exist", peer)
			return nil
		}
	}

	err = peer.UpdatePeer(envelope)
	if err != nil {
		log.Printf("HandShake error: %s", err)
	}

	return peer
}
