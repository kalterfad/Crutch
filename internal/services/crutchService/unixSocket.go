package crutchService

import (
	types "crutch/internal/services"
	"encoding/json"
	"log"
	"net"
	"os"
)

type unixSocket struct {
	socket         net.Listener
	cancelWorkerCh chan types.ChannelID
	newWorkerCh    chan types.DirData
}

func newSocket(newWorkerCh chan types.DirData, cancelChan chan types.ChannelID) *unixSocket {
	return &unixSocket{cancelWorkerCh: cancelChan, newWorkerCh: newWorkerCh}
}

func (s *unixSocket) run() error {
	if err := s.createSocket(); err != nil {
		return err
	}
	if err := s.waitEvent(); err != nil {
		return err
	}
	return nil
}

func (s *unixSocket) createSocket() error {
	_ = os.Remove("/tmp/crutch.sock")
	log.Println("start unix socket")
	sock, err := net.Listen("unix", "/tmp/crutch.sock")
	if err != nil {
		return err
	}
	s.socket = sock
	return nil
}

func (s *unixSocket) waitEvent() error {
	for {
		conn, err := s.socket.Accept()
		if err != nil {
			return err
		}
		go s.eventHandler(conn)
	}
}

func (s *unixSocket) eventHandler(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	dec := json.NewDecoder(conn)
	var data types.SocketMsg

	if err := dec.Decode(&data); err != nil {
		log.Println(err)
		return
	}

	switch data.Command {
	case types.RmWorker:
		log.Printf("worker removal command '%v'", data.ChannelID)
		s.cancelWorkerCh <- data.ChannelID
		s.write(conn, []byte("Directory scan stopped\n"))

	case types.NewWorker:
		log.Printf("command to create a new worker '%v'", data.ChannelID)
		s.newWorkerCh <- types.DirData{
			ChannelID: data.ChannelID,
			Path:      data.Path,
			Rules:     data.Rules,
		}
		s.write(conn, []byte("Catalog sent for scanning\n"))
	}
}

func (s *unixSocket) write(conn net.Conn, message []byte) {
	_, err := conn.Write(message)
	if err != nil {
		log.Fatalln(err)
	}
}
