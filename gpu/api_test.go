package gpu

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScryptPositions(t *testing.T) {
	req := require.New(t)

	id := []byte("id")
	salt := []byte("salt")
	startPosition := uint64(1)
	endPosition := uint64(2)
	hashLenBits := uint8(7)

	output, err := ScryptPositions(id, salt, startPosition, endPosition, CPU, hashLenBits)
	req.NoError(err)
	req.NotNil(output)
	fmt.Printf("%v\n", output)
}

func TestStats(t *testing.T) {
	req := require.New(t)

	c := Stats()
	req.True(c.CPU)
	req.False(c.GPUCuda)
	req.False(c.GPUOpenCL)
	req.True(c.GPUVulkan)
}

func TestGPUCount(t *testing.T) {
	req := require.New(t)

	count := GPUCount(true)
	req.Equal(2, count)
}
