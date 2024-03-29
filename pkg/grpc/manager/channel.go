package manager

import (
	"google.golang.org/grpc"
)

type Channel struct {
	appId  string
	cli    *grpc.ClientConn
	target string
}

func (c *Channel) Close() error {
	if c.cli == nil {
		return nil
	}

	return c.cli.Close()
}
