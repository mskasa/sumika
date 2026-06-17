package server

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mskasa/sumika/internal/config"
	"github.com/mskasa/sumika/internal/git"
	"github.com/mskasa/sumika/internal/launcher"
	"github.com/mskasa/sumika/internal/project"
	webfiles "github.com/mskasa/sumika/web"
)

type Server struct {
	cfg  *config.Config
	tmpl *template.Template
}

type ProjectView struct {
	Project        config.Project
	Status         *git.Status
	LastCommitRel  string
	ContextFileMod string
}

type PageData struct {
	Projects []ProjectView
}

func New(cfg *config.Config) *Server {
	tmpl := template.Must(template.ParseFS(webfiles.FS, "templates/index.html", "templates/cards.html"))
	return &Server{cfg: cfg, tmpl: tmpl}
}

func (s *Server) Run(port int) error {
	staticFS, err := fs.Sub(webfiles.FS, "static")
	if err != nil {
		return fmt.Errorf("sub static fs: %w", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	r.Get("/", s.handleIndex)
	r.Get("/api/projects", s.handleProjects)
	r.Post("/api/projects/{name}/open", s.handleOpen)

	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: addr, Handler: r}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		slog.Info("shutting down dashboard")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	slog.Info("starting dashboard", "addr", addr, "url", fmt.Sprintf("http://localhost%s", addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := s.buildPageData(s.cfg)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, data); err != nil {
		slog.Error("render index", "err", err)
	}
}

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		cfg = s.cfg
	}
	data := s.buildPageData(cfg)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "cards", data); err != nil {
		slog.Error("render cards", "err", err)
	}
}

func (s *Server) handleOpen(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	cfg, err := config.Load()
	if err != nil {
		cfg = s.cfg
	}
	p, err := project.Find(cfg, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	editor := ""
	if p.Launch.Editor {
		editor = cfg.Settings.Editor
	}
	if err := launcher.OpenBackground(p.Path, editor, p.Launch.Commands); err != nil {
		slog.Error("open project", "name", name, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) buildPageData(cfg *config.Config) PageData {
	views := make([]ProjectView, 0, len(cfg.Projects))
	for _, p := range cfg.Projects {
		st, err := git.GetStatus(p.Path, nil)
		if err != nil || st == nil {
			st = &git.Status{}
		}

		view := ProjectView{
			Project:       p,
			Status:        st,
			LastCommitRel: "",
		}
		if st.IsRepo && !st.LastCommitTime.IsZero() {
			view.LastCommitRel = relativeTime(st.LastCommitTime)
		}

		if info, err := os.Stat(filepath.Join(p.Path, "CLAUDE.md")); err == nil {
			view.ContextFileMod = relativeTime(info.ModTime())
		}

		views = append(views, view)
	}
	return PageData{Projects: views}
}

func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "たった今"
	case d < time.Hour:
		return fmt.Sprintf("%d分前", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d時間前", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%d日前", int(d.Hours()/24))
	default:
		return t.Local().Format("2006-01-02")
	}
}
