package gpu

import (
	"fmt"
	"time"
)

func GetProviders() []ComputeProvider {
	return cGetProviders()
}

func ScryptPositions(providerId uint, id, salt []byte, startPosition, endPosition uint64, hashLenBits uint8, options uint32) ([]byte, error, int) {
	if hashLenBits < 1 || hashLenBits > 8 {
		return nil, fmt.Errorf("invalid hashLenBits value; expected: 1-8, given: %v", hashLenBits), 0
	}

	const n, r, p = 512, 1, 1
	outputSize := calcOutputSize(startPosition, endPosition, hashLenBits)

	output, retVal, count := cScryptPositions(providerId, id, salt, startPosition, endPosition, hashLenBits, options, outputSize, n, r, p)
	switch retVal {
	case 0:
		return output, nil, count
	case -1:
		return nil, fmt.Errorf("invalid provider id: %v", providerId), 0
	default:
		panic("unreachable")
	}
}

func Benchmark(providerId uint) (uint64, error) {
	// TODO(moshababo): once fixed, use the stop function, and make the benchmark
	// function to run for a defined time duration.
	id := []byte("id")
	salt := []byte("salt")
	startPosition := uint64(1)
	endPosition := uint64(1 << 14)
	hashLenBits := uint8(8)
	options := uint32(0)

	// TODO(moshababo): refactor ScryptPositions to return the internal time duration, which is more accurate.
	t := time.Now()
	output, err, _ := ScryptPositions(providerId, id, salt, startPosition, endPosition, hashLenBits, options)
	d := time.Since(t)
	if err != nil {
		return 0, err
	}

	return uint64(float64(len(output)) / d.Seconds()), nil
}

func Stop() StopResult {
	return cStop(2000)
}
