package crutchClient

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type UnixSocket struct {
	socket net.Conn
}

func newSocket() *UnixSocket {
	return &UnixSocket{}
}

func (s *UnixSocket) connectToSocket() error {
	conn, err := net.Dial("unix", "/tmp/crutch.sock")
	if err != nil {
		return fmt.Errorf("\x1b[31m\u001B[1mCrutch service is not running\x1b[0m")
	}
	s.socket = conn
	return nil
}

func (s *UnixSocket) readSock() (string, error) {
	reader := bufio.NewReader(s.socket)
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return "\u001B[1m" + str, nil
}

func (s *UnixSocket) closeConn() {
	err := s.socket.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *UnixSocket) write(message []byte) error {
	_, err := s.socket.Write(message)
	if err != nil {
		return err
	}
	return nil
}
