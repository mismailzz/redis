package main

import "net"

type Peer struct {
	conn    net.Conn
	msgChan chan []byte
}

func NewPeer(conn net.Conn, readChan chan []byte) *Peer {
	return &Peer{
		conn:    conn,
		msgChan: readChan,
	}
}

func (p *Peer) readLoop() error {
	buff := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buff)
		if err != nil {
			return err
		}
		rbuff := make([]byte, n)
		copy(rbuff, buff[:n])
		p.msgChan <- rbuff
	}
}
