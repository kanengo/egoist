package manager

import (
	"google.golang.org/grpc"
)

type Channel struct {
	appId string
	cli   grpc.ClientConnInterface
}
