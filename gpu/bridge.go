package gpu

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR} -lgpu-setup
// #include "./api.h"
// #include <stdlib.h>
//
import "C"
import "unsafe"

type cUchar = C.uchar
type cUint = C.uint

type Options uint32

const (
	CPU       Options = C.SPACEMESH_API_CPU
	GPUCuda           = C.SPACEMESH_API_CUDA
	GPUOpenCL         = C.SPACEMESH_API_OPENCL
	GPUVulkan         = C.SPACEMESH_API_VULKAN
	GPUAll            = C.SPACEMESH_API_GPU
	ALL               = C.SPACEMESH_API_ALL
)

func cScryptPositions(id, salt []byte, startPosition, endPosition uint64, options Options, hashLenBits uint8, outputSize uint64, n, r, p uint32) ([]byte, int) {
	cId := (*C.uchar)(GoBytes(id).CBytesClone().data)
	cStartPosition := C.uint64_t(startPosition)
	cEndPosition := C.uint64_t(endPosition)
	cHashLenBits := C.uchar(hashLenBits)
	cSalt := (*C.uchar)(GoBytes(salt).CBytesClone().data)
	cOptions := C.uint(options)
	cOut := (*C.uchar)(GoBytes(make([]byte, outputSize)).CBytesClone().data)
	cN := C.uint(n)
	cR := C.uint(r)
	cP := C.uint(p)

	retVal := C.scryptPositions(
		cId,
		cStartPosition,
		cEndPosition,
		cHashLenBits,
		cSalt,
		cOptions,
		cOut,
		cN,
		cR,
		cP,
	)

	//const uint8_t *id,			// 32 bytes
	//uint64_t start_position,	// e.g. 0
	//	uint64_t end_position,		// e.g. 49,999
	//	uint8_t hash_len_bits,		// (1...8) for each hash output, the number of prefix bits (not bytes) to copy into the buffer
	//const uint8_t *salt,		// 32 bytes
	//uint32_t options,			// throttle etc.
	//	uint8_t *out,				// memory buffer large enough to include hash_len_bits * number of requested hashes
	//	uint32_t N,					// scrypt N
	//	uint32_t R,					// scrypt r
	//	uint32_t P					// scrypt p
	//);

	return cBytesCloneToGoBytes(cOut, int(outputSize)), int(retVal)
}

func cStats() int {
	return int(C.stats())
}

func cGPUCount(onlyAvailable bool) int {
	var cOnlyAvailable C.int
	if onlyAvailable {
		cOnlyAvailable = 1
	}

	return int(C.spacemesh_api_get_gpu_count(0, cOnlyAvailable))
}

func cFree(p unsafe.Pointer) {
	C.free(p)
}
