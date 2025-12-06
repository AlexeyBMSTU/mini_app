package server

import (
	"fmt"
	"log"
	"net/http"
	"mini-app-backend/internal/config"
)

type Server struct {
	config *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", s.rootHandler)
	mux.HandleFunc("/health", s.healthHandler)
	
	port := s.config.ServerPort
	
	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Starting HTTP server on PORT with %s", port)
	
	return http.ListenAndServe(addr, mux)
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	fmt.Fprintf(w, "🤖 Server start! PORT: %s", s.config.ServerPort)
	log.Printf("📝 Request to the root path: %s", r.RemoteAddr)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "port": "%s"}`, s.config.ServerPort)
	log.Printf("💓 Check health от %s", r.RemoteAddr)
}