package main

import (
	"encoding/json"
	"fmt"
	"net"

	"byted/DB_engine/constants"
	"byted/DB_engine/core/auth"
	"byted/DB_engine/core/bucket"
	"byted/DB_engine/structs"
	"byted/DB_engine/cmd/cli"
)

type Server struct {
	ListenAddr string

	Listener net.Listener
}
type ClientContext struct {
	BucketManager *bucket.BucketManager
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
	return &Server{
		ListenAddr: listenAddr,
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
		bm, _ := bucket.NewBucketManager(constants.DBBUCKETSPATH)

		ctx := &ClientContext{
			BucketManager: bm,
		}

		if !auth.HandleAuthenticatedConnection(comm) {
			conn.Close()
		}

		go s.readLoop(comm, ctx, conn)

	}
}

func (s *Server) readLoop(comm *structs.Communicators, ctx *ClientContext, conn net.Conn) {
	for {
		var msg structs.Message

		// Read message from client
		if err := comm.Dec.Decode(&msg); err != nil {
			fmt.Println("Client disconnected:", err)
			return
		}
		ActiveBucket, _ := ctx.BucketManager.GetActiveBucket()

		switch msg.Type {
		case "command":
			// Execute the command
			var data []string
			var err error

			if ActiveBucket == nil {
				data, err = cli.ExecuteGlobalCommmand(msg.Command, ctx.BucketManager, conn)
			} else {
				data, err = cli.ExecuteCommand(msg.Command, ActiveBucket, ctx.BucketManager)
			}
			ActiveBucket, _ = ctx.BucketManager.GetActiveBucket()
			var currentBkt string
			if ActiveBucket == nil {
				currentBkt = ""
			} else {
				currentBkt = ActiveBucket.Name
			}
			if err != nil {
				// Send error back to client
				comm.Enc.Encode(structs.Message{
					Type:    "error",
					Message: err.Error(),
					Bucket:  currentBkt,
				})
				continue
			} else {
				comm.Enc.Encode(structs.Message{
					Type:   "success",
					Data:   data,
					Bucket: currentBkt,
				})
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
