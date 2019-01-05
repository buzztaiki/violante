package violante

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Server ...
type Server struct {
	addr string
	det  *Detector
}

// NewServer ...
func NewServer(addr string, det *Detector) *Server {
	return &Server{addr: addr, det: det}
}

// Route ...
func (s *Server) Route() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handler(s.add))
	return mux
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	log.Printf("server started on %s", s.addr)
	return http.ListenAndServe(s.addr, s.Route())
}

func (s *Server) handler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := fn(w, req); err != nil {
			log.Printf("[error] %s %s %v", req.Method, req.RequestURI, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) add(w http.ResponseWriter, req *http.Request) error {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	r := AddRequest{}
	if err := json.Unmarshal(body, &r); err != nil {
		return err
	}

	for _, f := range r.Files {
		s.det.Add(f)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
