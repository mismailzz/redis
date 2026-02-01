package main

import (
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
	}
}

func (s *Server) Start() error {
	var err error
	s.ln, err = net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	go s.loop()

	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.handleNewConnection(conn)
	}
}

func (s *Server) handleNewConnection(conn net.Conn) {
	peer := NewPeer(conn)
	s.peerAddCh <- peer
	peer.readLoop()
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.peerAddCh:
			s.peer[peer] = true
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
