package server

import (
	"fmt"
	"net/http"
	"tender_srevice/internal/app"
	"tender_srevice/internal/config"
	"tender_srevice/internal/repository"

	"github.com/gorilla/mux"
)

type Server struct {
	config *config.Config
	router *mux.Router
	repo   *repository.PostgresRepository
}

func New(cfg *config.Config, repo *repository.PostgresRepository) *Server {
	s := &Server{
		config: cfg,
		repo:   repo,
	}
	s.router = app.SetupRouter(cfg, repo)
	return s
}

func (s *Server) Run() error {
	fmt.Printf("Server is running on %s\n", s.config.ServerAddress)
	return http.ListenAndServe(s.config.ServerAddress, s.router)
}