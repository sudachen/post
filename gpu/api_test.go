package gpu

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

var (
	id      = []byte("id")
	salt    = []byte("salt")
	options = uint32(0)
)

//
func TestScryptPositions(t *testing.T) {
	req := require.New(t)
	providers := GetProviders()
	fmt.Printf("providers: %+v\n", providers)
	for _, p := range providers {
		providerId := uint(p.Id)
		startPosition := uint64(1)
		endPosition := uint64(1024)
		hashLenBits := uint8(8)
		output, err, count := ScryptPositions(providerId, id, salt, startPosition, endPosition, hashLenBits, options)
		req.NoError(err)
		req.NotNil(output)
		req.Equal(len(output), count)
		//fmt.Printf("hashes_computed: %v\n\n", count)
	}
}

//
//func TestScryptPositions_InvalidProviderId(t *testing.T) {
//	req := require.New(t)
//
//	invalidProviderId := uint(1 << 10)
//	output, err := ScryptPositions(invalidProviderId, id, salt, 1, 1, 8, options)
//	req.EqualError(err, fmt.Sprintf("invalid provider id: %d", invalidProviderId))
//	req.Nil(output)
//}

func TestStop(t *testing.T) {
	req := require.New(t)

	providers := GetProviders()
	for _, p := range providers {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			time.Sleep(2 * time.Second)
			res := cStop(2000)
			req.Equal(StopResultOk, res)
		}()
		go func() {
			defer wg.Done()
			providerId := uint(p.Id)
			startPosition := uint64(1)
			endPosition := uint64(1 << 19)
			hashLenBits := uint8(8)
			output, err, count := ScryptPositions(providerId, id, salt, startPosition, endPosition, hashLenBits, options)
			req.NoError(err)
			req.NotNil(output)
			req.Equal(1<<19, len(output))

			output = getNonEmpty(output)
			//	req.Equal(len(output), count)
			fmt.Printf("hashes_computed: %v\n", count)
			fmt.Printf("output_inspection: %v\n\n", len(output))

		}()
		c := make(chan struct{})
		go func() {
			defer close(c)
			wg.Wait()
		}()
		select {
		case <-c:
		case <-time.After(5 * time.Second):
			req.Fail(fmt.Sprintf("stop timed out; provider: %+v", p))
		}
	}
}

//func TestBenchmark(t *testing.T) {
//	req := require.New(t)
//
//	results := make(map[string]uint64)
//	providers := GetProviders()
//	for _, p := range providers {
//		providerId := uint(p.Id)
//		hps, err := Benchmark(providerId)
//		req.NoError(err)
//		results[p.Model] = hps
//	}
//
//	fmt.Printf("%v\n", results)
//}

func getNonEmpty(b []byte) []byte {
	const repeatedTarget = 10
	var repeated int
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			repeated++
			if repeated == repeatedTarget {
				return b[:i-repeatedTarget+1]
			}
		} else {
			repeated = 0
		}
	}
	return b
}
