package malloc

import (
	"math"
	"os"
	"sync"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

var pageSize = os.Getpagesize()

// Blocks allocated from MMap
type block struct {
	data           []byte // The actual data to allocate from
	allocatedBytes int
	nOfAllocations int
	size           int
}

// Metadata for each allocation
type metadata struct {
	size  int
	block *block
}

var sizeOfBlockStruct = unsafe.Sizeof(block{})
var sizeOfDataFieldInBlock = uintptr(pageSize) - sizeOfBlockStruct // The size of the data field in block struct (= pageSize - size of block struct)
var sizeOfMetaStruct = unsafe.Sizeof(metadata{})
var headBlock *block // the starting block to search

func createBlock(size int) *block {
	// How many pages to allocate rounded up
	pageMultiplier := int(math.Ceil(float64(size+int(sizeOfMetaStruct)) / float64(sizeOfDataFieldInBlock)))

	// the size to allocate from MMap
	blockSize := pageSize * pageMultiplier

	// byte array from mmap
	bytes, err := mmap.MapRegion(nil, blockSize, mmap.RDWR, mmap.ANON, 0)
	if err != nil {
		panic(err)
	}

	// casting the byte array to block struct
	block := (*block)(unsafe.Pointer(&bytes[0]))
	block.size = blockSize
	// block data is the remainder of the block struct size
	block.data = unsafe.Slice(
		(*byte)(unsafe.Pointer(&bytes[sizeOfBlockStruct])),
		blockSize-int(sizeOfBlockStruct),
	)

	return block
}

var mlock sync.Mutex

func malloc(size int) unsafe.Pointer {
	mlock.Lock()
	defer mlock.Unlock()

	if headBlock == nil {
		headBlock = createBlock(size)
	}

	currentBlock := headBlock

	for {
		if currentBlock.allocatedBytes+int(sizeOfMetaStruct)+size <= len(currentBlock.data) {
			mdata := (*metadata)(unsafe.Pointer(&currentBlock.data[currentBlock.allocatedBytes]))
			mdata.size = size
			mdata.block = currentBlock

			ptr := unsafe.Pointer(&currentBlock.data[currentBlock.allocatedBytes+int(sizeOfMetaStruct)])

			currentBlock.allocatedBytes += int(sizeOfMetaStruct) + size
			currentBlock.nOfAllocations++

			return ptr
		}

		currentBlock = createBlock(size)
	}
}

func free(ptr unsafe.Pointer) {
	mlock.Lock()
	defer mlock.Unlock()

	mdata := (*metadata)(unsafe.Pointer(uintptr(ptr) - sizeOfMetaStruct))
	mdata.block.nOfAllocations--

	if mdata.block.nOfAllocations == 0 {
		if mdata.block == headBlock {
			headBlock = nil
		}

		mmapBytes := mmap.MMap(
			unsafe.Slice(
				(*byte)(unsafe.Pointer(mdata.block)),
				mdata.block.size,
			),
		)

		if err := mmapBytes.Unmap(); err != nil {
			panic(err)
		}
	}
}

func realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	oldMeta := (*metadata)(unsafe.Pointer(uintptr(ptr) - sizeOfMetaStruct))
	oldData := unsafe.Slice(
		(*byte)(ptr),
		oldMeta.size,
	)

	newPtr := malloc(size)
	newPtrData := unsafe.Slice(
		(*byte)(newPtr),
		size,
	)

	copy(newPtrData, oldData)

	free(ptr)

	return newPtr
}

// CMalloc raw binding to c calloc(1, size)
func Malloc(size int) unsafe.Pointer {
	return malloc(size)
}

// CMalloc raw binding to c free
func Free(ptr unsafe.Pointer) {
	free(ptr)
}

// CMalloc raw binding to c realloc
func Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return realloc(ptr, size)
}
