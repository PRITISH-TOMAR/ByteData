package main

import (
	"encoding/json"
	"fmt"
	"net"

	"byted/constants"
	"byted/internal/auth"
	"byted/internal/bucket"
	"byted/structs"
	"byted/cmd/cli"
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
		defer conn.Close()

		comm := communicators(conn)
		if !auth.HandleAuthenticatedConnection(comm) {
			conn.Close()
		}

		go s.readLoop(comm)

	}
}

func (s *Server) readLoop(comm *structs.Communicators) {
    for {
        var msg structs.Message

        // Read message from client
        if err := comm.Dec.Decode(&msg); err != nil {
            fmt.Println("Client disconnected:", err)
            return
        }

        switch msg.Type {
        case "command":
            // Execute the command
            err := cli.ExecuteGlobalCommmand(comm, msg.Command, s.BucketManager)
            if err != nil {
                // Send error back to client
                comm.Enc.Encode(structs.Message{
                    Type:    "error",
                    Message: err.Error(),
                })
                continue
            }

        default:
            // Unknown message type
            comm.Enc.Encode(structs.Message{
                Type:    "error",
                Message: "Unknown message type: " + msg.Type,
            })
        }
    }
}


func main() {
	server := NewServer(":8080")
	server.Start()
}
