package mm

import (
	"math"
	"os"
	"sync"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

var pageSize = os.Getpagesize()
var bigChunkThreshold = 1000000 // 1mb

// Blocks allocated from MMap
type block struct {
	data      []byte // The actual data to allocate from
	usedBytes int
	size      int
	nextBlock *block // Pointer to the next block
	prevBlock *block
}

// Metadata for each allocation
type metadata struct {
	size  int
	free  int8 // 0 free 1 not free
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
	// block data is the remainder of the block struct size
	block.data = unsafe.Slice(
		(*byte)(unsafe.Pointer(&bytes[sizeOfBlockStruct])),
		blockSize-int(sizeOfBlockStruct),
	)
	block.size = blockSize

	return block
}

// mutex to avoid race conditions (thread safety)
var mallocMutex sync.Mutex

func malloc(size int) unsafe.Pointer {
	mallocMutex.Lock()
	defer mallocMutex.Unlock()

	// initialize
	if headBlock == nil {
		headBlock = createBlock(size)
	}

	// search for a free block
	currentBlock := headBlock
	for {
		// get the metadata struct
		metaPtr := (*metadata)(unsafe.Pointer(&currentBlock.data[0]))

		// checks if the metaPtr can be inside the data field in the block
		for uintptr(unsafe.Pointer(metaPtr))+sizeOfMetaStruct+uintptr(size)-uintptr(unsafe.Pointer(&currentBlock.data[0])) <= uintptr(len(currentBlock.data)) {
			if metaPtr.free == 0 {
				if metaPtr.size == 0 {
					// first time (free and size is zero)
					metaPtr.size = size // sets the size to the wanted size
					metaPtr.block = currentBlock
				}

				// checks if the size is greater than the wanted size
				if metaPtr.size >= size {
					metaPtr.free = 1 // make it not available
					metaPtr.block.usedBytes += metaPtr.size

					ptr := unsafe.Pointer(uintptr(unsafe.Pointer(metaPtr)) + sizeOfMetaStruct)

					// finally returns the pointer to the data after the meta struct
					return ptr
				}
			}

			// if not found check the next meta struct
			metaPtr = (*metadata)(
				unsafe.Pointer(
					uintptr(unsafe.Pointer(metaPtr)) +
						sizeOfMetaStruct +
						uintptr(metaPtr.size),
				),
			)
		}

		// creates another block if the current one have no next block
		if currentBlock.nextBlock == nil {
			newBlock := createBlock(size)
			newBlock.prevBlock = currentBlock
			currentBlock.nextBlock = newBlock
		}

		// if the block is all occupied check the next block
		currentBlock = currentBlock.nextBlock
	}
}

// Free frees a pointer from Malloc
func free(ptr unsafe.Pointer) {
	// gets the metadata struct
	metaData := (*metadata)(unsafe.Pointer(uintptr(ptr) - sizeOfMetaStruct))

	// zero out memory
	ptrData := unsafe.Slice(
		(*byte)(ptr),
		metaData.size,
	)
	for i := range ptrData {
		ptrData[i] = 0
	}

	// makes it available
	metaData.free = 0
	metaData.block.usedBytes -= metaData.size

	// free the chunk if it's bigger than or equal the threshold
	if metaData.block.usedBytes == 0 && metaData.block.size >= bigChunkThreshold {
		if metaData.block == headBlock {
			headBlock = nil
		}

		nextBlock := metaData.block.nextBlock
		prevBlock := metaData.block.prevBlock
		if prevBlock != nil {
			prevBlock.nextBlock = nextBlock
		}
		if nextBlock != nil {
			nextBlock.prevBlock = prevBlock
		}

		mmap := mmap.MMap(unsafe.Slice(
			(*byte)(unsafe.Pointer(metaData.block)),
			metaData.block.size,
		))
		err := mmap.Unmap()
		if err != nil {
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
