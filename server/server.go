package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Server struct {
	http *http.Server
	mux  *http.ServeMux
}

func NewServer() *Server {
	return &Server{
		http: &http.Server{},
		mux:  http.NewServeMux(),
	}
}

func setCorsHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (s *Server) Func(f func(s *Server)) *Server {
	f(s)
	return s
}

// HandleFunc registers access routes
func (s *Server) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

// Handle registers access routes
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) HttpRun(addr string) {
	s.http.Addr = addr
	s.http.Handler = setCorsHeaders(s.mux) // Wrap the mux with CORS headers

	logrus.Debugf("http listen %s", addr)

	// Start listening
	err := s.http.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
