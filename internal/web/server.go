package web

import (
	"embed"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/yuzeyzer/gomoku/internal/gomoku"
)

//go:embed static/*
var staticFS embed.FS

type Server struct {
	mu   sync.Mutex
	game *gomoku.Game
}

func NewServer(size int) *Server {
	return &Server{
		game: gomoku.NewGame(size),
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	mux.HandleFunc("/api/state", s.handleState)
	mux.HandleFunc("/api/move", s.handleMove)

	return mux
}

type stateResponse struct {
	Size  int              `json:"size"`
	Turn  gomoku.Stone     `json:"turn"`
	Moves int              `json:"moves"`
	Board [][]gomoku.Stone `json:"board"`
}

type moveRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	writeJSON(w, http.StatusOK, s.snapshotLocked())
}

func (s *Server) handleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.game.Play(gomoku.Point{X: req.X, Y: req.Y}); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, s.snapshotLocked())
}

func (s *Server) snapshotLocked() stateResponse {
	size := s.game.Board.Size()
	board := make([][]gomoku.Stone, size)

	for y := 0; y < size; y++ {
		row := make([]gomoku.Stone, size)
		for x := 0; x < size; x++ {
			st, _ := s.game.Board.Get(gomoku.Point{X: x, Y: y})
			row[x] = st
		}
		board[y] = row
	}

	return stateResponse{
		Size:  size,
		Turn:  s.game.Turn,
		Moves: s.game.Moves,
		Board: board,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
