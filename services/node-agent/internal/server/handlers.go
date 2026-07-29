package server

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func (s *Server) handleSystem(w http.ResponseWriter, r *http.Request) {
	info := s.systemService.GetInfo()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(info)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{
		Status:  "ok",
		Service: "sentinel-node-agent",
		Version: "0.1.0-dev",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}