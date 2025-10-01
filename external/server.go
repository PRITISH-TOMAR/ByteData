package main

import (
	"encoding/json"
	"fmt"
	"net"

	"byted/constants"
	"byted/internal/auth"
	"byted/internal/bucket"
	"byted/structs"
)

type Server struct {
	ListenAddr    string
	BucketManager *bucket.BucketManager
	Listener      net.Listener
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	s.Listener = ln
	defer s.Listener.Close()
	fmt.Printf("Server listening on %s\n", s.ListenAddr)

	s.acceptLoop()
	return nil
}

func NewServer(listenAddr string) *Server {
	bm, _ := bucket.NewBucketManager(constants.DBBUCKETSPATH)
	return &Server{
		ListenAddr:    listenAddr,
		BucketManager: bm,
	}
}
func communicators(conn net.Conn) *structs.Communicators {
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	return &structs.Communicators{Enc: enc, Dec: dec}
}

func (s *Server) acceptLoop() {

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Fprintf(conn, "failed to accept connection: %v\n", err)
			continue
		}
		comm := communicators(conn)
		if !auth.HandleAuthenticatedConnection(comm) {
			conn.Close()
		}

		go s.readLoop(conn)

	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	// comm := communicators(conn)

	// for {
	// 	fmt.Fprintf(conn, "\nByteData: ") // prompt
	// 	n, err := conn.Read(buff)
	// 	if err != nil {
	// 		fmt.Printf("failed to read from connection: %v\n", err)
	// 		return // stop reading on error
	// 	}

	// 	message := string(buff[:n])
	// 	// Execute command with net.Conn directly, not pointer
	// 	err = cli.ExecuteGlobalCommmand(ctx.Connection, message, s.BucketManager)
	// 	if err != nil {

	// 		fmt.Printf("command execution error: %v\n", err)
	// 	}
	// }
}

func main() {
	server := NewServer(":8080")
	server.Start()
}
