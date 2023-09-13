package grpc

import (
	"testing"

	"github.com/kanengo/egoist/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewSever(t *testing.T) {
	hostAddress, err := utils.GetHostAddress()
	assert.Nil(t, err)
	s := NewAPIServer(NewAPI(APIOptions{}), ServerConfig{
		AppID:                "test",
		HostAddress:          hostAddress,
		Port:                 54167,
		APIListenAddresses:   nil,
		NameSpace:            "test",
		MaxRequestBodySizeMB: 4,
		UnixDomainSocket:     "/Users/kanonlee/project/github/egoist",
		ReadBufferSizeKB:     4,
	}, nil)
	assert.Nil(t, s.StartNonBlocking())
	s.Close()
}
