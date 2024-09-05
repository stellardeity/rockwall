package proto

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ed25519"
)

type SharedKey struct {
	RemoteKey []byte
	LocalKey  []byte
	Secret    []byte
}

func (sk *SharedKey) Update(remoteKey []byte, localKey []byte) {
	if remoteKey != nil {
		sk.RemoteKey = remoteKey
	}

	if localKey != nil {
		sk.LocalKey = localKey
	}

	if sk.RemoteKey != nil && sk.LocalKey != nil {
		secret := CalcSharedSecret(sk.RemoteKey, sk.LocalKey)
		sk.Secret = secret[:32]
	}
}

type Peer struct {
	PubKey    ed25519.PublicKey
	Conn      *net.Conn
	Name      string
	FirstSeen string
	LastSeen  string
	Peers     *Peers
	SharedKey SharedKey
}

func (p Peer) String() string {
	return string(p.Name) + ":" + hex.EncodeToString(p.PubKey)
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		PubKey:    nil,
		Conn:      &conn,
		Name:      conn.RemoteAddr().String(),
		FirstSeen: time.Now().String(),
		LastSeen:  time.Now().String(),
		Peers:     NewPeers(),
		SharedKey: SharedKey{
			RemoteKey: nil,
			LocalKey:  nil,
			Secret:    nil,
		},
	}
}

func (p *Peer) UpdatePeer(envelope *Envelope) error {
	if string(envelope.Cmd) != "HAND" {
		return errors.New("invalid command")
	}

	handShake := &HandShake{}
	err := json.Unmarshal(envelope.Content, handShake)
	if err != nil {
		return err
	}

	rawPubKey, err := hex.DecodeString(handShake.PubKey)
	if err != nil {
		return err
	}

	rawExKey, err := hex.DecodeString(handShake.ExKey)
	if err != nil {
		return err
	}

	p.Name = handShake.Name
	p.PubKey = rawPubKey

	p.SharedKey.Update(rawExKey, nil)
	return nil
}

type Peers struct {
	sync.RWMutex
	peers map[string]*Peer
}

func NewPeers() *Peers {
	return &Peers{
		peers: make(map[string]*Peer),
	}
}

func (p Peers) Put(peer *Peer) {
	p.Lock()
	defer p.Unlock()

	p.peers[string(peer.PubKey)] = peer
}

func (p Peers) Get(key string) (peer *Peer, found bool) {
	p.RLock()
	defer p.RUnlock()

	peer, found = p.peers[key]
	return
}

func (p Peers) Remove(peer *Peer) (found bool) {
	p.RLock()
	defer p.RUnlock()

	delete(p.peers, string(peer.PubKey))
	return
}

func (p Peers) PeerList() *WsPeerList {

	peerList := &WsPeerList{
		WsCmd: WsCmd{
			Cmd: "PEERS",
		},
		Peers: []PeerName{},
	}

	log.Printf("total peers: %v", len(p.peers))

	for _, el := range p.peers {
		peerList.Peers = append(peerList.Peers, PeerName{
			Name:   el.Name,
			PubKey: hex.EncodeToString(el.PubKey),
		})
	}

	return peerList
}
