package smtp

import (
	"fmt"
	"log"
	"net"
	"time"

	"email-service/config"

	"github.com/emersion/go-smtp"
)

// Server represents the SMTP server
type Server struct {
	server *smtp.Server
}

// NewServer creates a new SMTP server
func NewServer() *Server {
	backend := &Backend{}

	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf("%s:%d", config.AppConfig.SMTP.Host, config.AppConfig.SMTP.Port)
	server.Domain = "localhost"
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024 // 1MB max message size
	server.MaxRecipients = 50
	server.AllowInsecureAuth = !config.AppConfig.SMTP.TLS
	server.AuthDisabled = !config.AppConfig.SMTP.AuthRequired

	return &Server{
		server: server,
	}
}

// Start starts the SMTP server
func (s *Server) Start() error {
	log.Printf("Starting SMTP server on %s", s.server.Addr)

	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.server.Addr, err)
	}

	log.Printf("SMTP server listening on %s", s.server.Addr)
	return s.server.Serve(listener)
}

// Stop stops the SMTP server
func (s *Server) Stop() error {
	log.Println("Stopping SMTP server...")
	return s.server.Close()
}
