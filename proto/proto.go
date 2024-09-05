package proto

import (
	"bufio"
	"encoding/hex"
	"io"
	"log"
	"os"
	"reflect"

	"golang.org/x/crypto/ed25519"
)

type Proto struct {
	Port     int
	Name     string
	Peers    *Peers
	PubKey   ed25519.PublicKey
	privKey  ed25519.PrivateKey
	Broker   chan *Envelope
	handlers map[string]func(peer *Peer, envelope *Envelope)
}

func (p Proto) MyName() *PeerName {
	return &PeerName{
		Name:   p.Name,
		PubKey: hex.EncodeToString(p.PubKey),
	}
}

func (p Proto) String() string {
	return "proto: " + hex.EncodeToString(p.PubKey) + ": " + p.Name
}

func getSeed() []byte {
	seed := getRandomSeed(32)

	fName := "seed.dat"
	file, err := os.Open(fName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(fName)
			if err != nil {
				panic(err)
			}
		}
	}

	_, err = file.Read(seed)
	if err != nil {
		panic(err)
	}
	return seed
}

func NewProto(name string, port int) *Proto {
	publicKey, privateKey := LoadKey(name)
	proto := &Proto{
		Port:     port,
		Name:     name,
		Peers:    NewPeers(),
		PubKey:   publicKey,
		privKey:  privateKey,
		Broker:   make(chan *Envelope),
		handlers: make(map[string]func(peer *Peer, envelope *Envelope)),
	}

	proto.handlers["HAND"] = proto.onHand
	proto.handlers["MESS"] = proto.onMess
	proto.handlers["LIST"] = proto.onList

	return proto
}

func (p Proto) SendName(peer *Peer) {

	exchPubKey, exchPrivKey := CreateKeyExchangePair()

	handShake := HandShake{
		Name:   p.Name,
		PubKey: hex.EncodeToString(p.PubKey),
		ExKey:  hex.EncodeToString(exchPubKey[:]),
	}.ToJson()

	peer.SharedKey.Update(nil, exchPrivKey[:])

	sign := ed25519.Sign(p.privKey, handShake)

	envelope := NewSignedEnvelope("HAND", p.PubKey[:], make([]byte, 32), sign, handShake)

	envelope.Send(peer)
}

func (p Proto) RequestPeers(peer *Peer) {
	envelope := NewEnvelope("LIST", []byte("TODO"))
	envelope.Send(peer)
}

func (p Proto) SendPeers(peer *Peer) {
	envelope := NewEnvelope("PEER", []byte("TODO"))
	envelope.Send(peer)
}

func (p Proto) SendMessage(peer *Peer, msg string) {
	if peer.SharedKey.Secret == nil {
		log.Fatalf("can't send message!")
	}

	encryptedMessage := Encrypt([]byte(msg), peer.SharedKey.Secret)

	envelope := NewSignedEnvelope("MESS", p.PubKey, peer.PubKey, ed25519.Sign(p.privKey, encryptedMessage), encryptedMessage)

	envelope.Send(peer)
}

func (p Proto) RegisterPeer(peer *Peer) *Peer {
	if reflect.DeepEqual(peer.PubKey, p.PubKey) {
		return nil
	}

	p.Peers.Put(peer)

	log.Printf("Register new peer: %s (%v)", peer.Name, len(p.Peers.peers))

	return peer
}

func (p Proto) UnregisterPeer(peer *Peer) {
	if p.Peers.Remove(peer) {
		log.Printf("UnRegister peer: %s", peer.Name)
	}
}

func (p Proto) ListenPeer(peer *Peer) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(*peer.Conn), bufio.NewWriter(*peer.Conn))
	p.HandleProto(readWriter, peer)
}

func (p Proto) HandleProto(rw *bufio.ReadWriter, peer *Peer) {
	for {
		envelope, err := ReadEnvelope(rw.Reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error on read Envelope: %v", err)
			}
			log.Printf("Disconnect peer %s", peer)
			break
		}

		if ed25519.Verify(envelope.From, envelope.Content, envelope.Sign) {
			log.Printf("Signed envelope!")
		}

		log.Printf("LISTENER: receive envelope from %s", (*peer.Conn).RemoteAddr())

		handler, found := p.handlers[string(envelope.Cmd)]

		if !found {
			log.Printf("LISTENER: UNHANDLED PROTO MESSAGE %v %v %v", envelope.Cmd, envelope.Id, envelope.Content)
			continue
		}

		handler(peer, envelope)
	}

	if peer != nil {
		p.UnregisterPeer(peer)
	}
}
