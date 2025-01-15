package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"tcp-aws-crud/config"
	"tcp-aws-crud/internal/db"
)

type Server struct {
	DB  *db.DB
	cfg *config.Server
}

func NewServer(_ context.Context, cfg config.Server, dynamoDB *db.DB) (*Server, error) {
	return &Server{
		DB:  dynamoDB,
		cfg: &cfg,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	cert, err := tls.LoadX509KeyPair(s.cfg.TLS.CertPath, s.cfg.TLS.CertKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %v", err)
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port), tlsCfg)
	if err != nil {
		return fmt.Errorf("failed to start TLS server: %v", err)
	}
	defer func(listener net.Listener) {
		err = listener.Close()
		if err != nil {
			log.Printf("Failed to close listener: %v", err)
		}
	}(listener)

	log.Printf("Server is listening on port %d...", s.cfg.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go func() {
			log.Printf("Handling connection from %s via %s...", conn.RemoteAddr().String(), conn.RemoteAddr().Network())
			if err = s.handleConnection(ctx, conn); err != nil {
				log.Printf("Failed to handle connection: %v", err)
			}
		}()
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) error {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}(conn)

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Connection closed by %s", conn.RemoteAddr().String())
				return nil
			}
			return fmt.Errorf("reading from connection: %w", err)
		}

		log.Printf("Processing request from %s: %s", conn.RemoteAddr().String(), message)
		response, err := s.processRequest(ctx, message)
		if err != nil {
			response = fmt.Sprintf("Failed %s: %s", response, err.Error())
		} else {
			response = fmt.Sprintf("SUCCESS %s", response)
		}
		log.Printf("Response for %s: %s", conn.RemoteAddr().String(), response)

		if _, err = conn.Write([]byte(response + "\n")); err != nil {
			return fmt.Errorf("writing response: %w", err)
		}
	}
}

func (s *Server) processRequest(ctx context.Context, request string) (string, error) {
	var command, id, data string
	n, err := fmt.Sscanf(request, "%s %s %s", &command, &id, &data)
	if err != nil && !errors.Is(err, io.EOF) && n < 2 {
		return "", fmt.Errorf("invalid request format. Use: COMMAND id data")
	}

	switch command {
	case "CREATE":
		if n < 3 {
			return "", fmt.Errorf("CREATE command requires data")
		}
		return "CREATE", s.DB.CreateItem(ctx, id, data)
	case "READ":
		itemData, err := s.DB.ReadItem(ctx, id)
		if err != nil {
			return "READ", err
		}
		return "READ: " + itemData, nil
	case "UPDATE":
		if n < 3 {
			return "", fmt.Errorf("UPDATE command requires data")
		}
		return "UPDATE", s.DB.UpdateItem(ctx, id, data)
	case "DELETE":
		return "DELETE", s.DB.DeleteItem(ctx, id)
	default:
		return "", fmt.Errorf("unknown command. Use CREATE, READ, UPDATE, DELETE")
	}
}
