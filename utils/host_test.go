package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostAdress(t *testing.T) {
	t.Run("EGOIST_HOST_IP present", func(t *testing.T) {
		hostIP := "test.local"
		t.Setenv(HostIPEnvVar, hostIP)

		address, err := GetHostAddress()
		assert.Nil(t, err)
		assert.Equal(t, hostIP, address)

		fmt.Println(address)
	})

	t.Run("EGOIST_HOST_IP not present, non-empty response", func(t *testing.T) {
		address, err := GetHostAddress()
		assert.Nil(t, err)
		assert.NotEmpty(t, address)

		fmt.Println(address)
	})
}
