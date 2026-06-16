package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mskasa/sumika/internal/config"
)

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Run(port int) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", s.handleIndex)
	r.Get("/api/projects", s.handleProjects)
	r.Post("/api/projects/{name}/open", s.handleOpen)

	addr := fmt.Sprintf(":%d", port)
	slog.Info("starting dashboard", "addr", addr)
	return http.ListenAndServe(addr, r)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<html><body><h1>sumika dashboard</h1><p>Phase 2 implementation pending.</p></body></html>")
}

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, p := range s.cfg.Projects {
		fmt.Fprintf(w, "<div>%s — %s</div>", p.Name, p.Description)
	}
}

func (s *Server) handleOpen(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
