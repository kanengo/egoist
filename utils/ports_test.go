package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	port1, err := GetStablePort(40000, "test")
	assert.Nil(t, err)

	port2, err := GetFreePort()
	assert.Nil(t, err)

	fmt.Println(port1, port2)
}
