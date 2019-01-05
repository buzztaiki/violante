package violante

import (
	"context"

	"github.com/buzztaiki/violante/rpc"
	"google.golang.org/grpc"
)

// Client ...
type Client struct {
	addr string
}

// NewClient ...
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

// Add ..
func (c *Client) Add(files []string) error {
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	rpcClient := rpc.NewViolanteClient(conn)

	ctx := context.Background()
	req := rpc.AddFilesRequest{Files: files}
	if _, err := rpcClient.AddFiles(ctx, &req); err != nil {
		return err
	}

	return nil
}
