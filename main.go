package main

import (
	"fmt"
	"net"
)

const defaultListenerAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln        net.Listener
	peer      map[*Peer]bool
	peerAddCh chan *Peer
	quitCh    chan struct{}
	readChan  chan []byte
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = defaultListenerAddr
	}
	return &Server{
		Config:    cfg,
		peer:      make(map[*Peer]bool),
		peerAddCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		readChan:  make(chan []byte),
	}
}

func (s *Server) Start() error {
	var err error
	s.ln, err = net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	fmt.Printf("Server listening on %s\n", s.ListenAddr)
	go s.loop()

	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Errorf("Error accepting connection: %v", err)
			continue
		}
		go s.handleNewConnection(conn)
	}
}

func (s *Server) handleNewConnection(conn net.Conn) {
	peer := NewPeer(conn, s.readChan)
	s.peerAddCh <- peer
	fmt.Printf("New peer connected: %s\n", conn.RemoteAddr().String())
	err := peer.readLoop()
	if err != nil {
		fmt.Printf("Error reading from peer %s: %v", conn.RemoteAddr().String(), err)
	}

}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.peerAddCh:
			s.peer[peer] = true
		case msg := <-s.readChan:
			fmt.Printf("Server received message: %s\n", string(msg))
		case <-s.quitCh:
			return
		}
	}
}

func main() {
	server := NewServer(Config{})
	if err := server.Start(); err != nil {
		panic(err)
	}
}
