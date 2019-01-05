package violante

import (
	"context"
	"net"

	"github.com/buzztaiki/violante/rpc"
	"google.golang.org/grpc"
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

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	grpcSvr := grpc.NewServer()
	rpc.RegisterViolanteServer(grpcSvr, s)

	return grpcSvr.Serve(l)
}

// AddFiles implements rpc.ViolanteServer.AddFiles.
func (s *Server) AddFiles(ctx context.Context, req *rpc.AddFilesRequest) (*rpc.AddFilesResponse, error) {
	for _, f := range req.Files {
		s.det.Add(f)
	}
	return &rpc.AddFilesResponse{}, nil
}
