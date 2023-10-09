package manager

import (
	"google.golang.org/grpc"
)

type Channel struct {
	appId string
	cli   *grpc.ClientConn
}

func (c *Channel) Close() error {
	if c.cli == nil {
		return nil
	}

	return c.cli.Close()
}
