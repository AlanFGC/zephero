package shared

import (
	"testing"
)

func makeOrderedChunk() *WorldChunk {
	arr := make([][]GNode, 10)
	for i := range arr {
		arr[i] = make([]GNode, 10)
	}

	for i := range arr {
		for j := range arr[i] {
			arr[i][j] = GNode{Id: uint64(i*len(arr) + j)}
		}
	}

	chunk := WorldChunk{
		chunkId: 0,
		data:    arr,
	}

	return &chunk
}

func TestSerializeChunk(t *testing.T) {

	chunk := makeOrderedChunk()

	serializedChunk, err := SerializeChunkData(chunk)
	if err != nil || serializedChunk == nil {
		t.Error(err)
	}
}

func TestDeserializeChunk(t *testing.T) {
	chunk := makeOrderedChunk()
	serializedChunk, err := SerializeChunkData(chunk)
	if err != nil || serializedChunk == nil {
		t.Error(err)
	}
	readChunk, err := DeserializeChunkData(serializedChunk)
	if err != nil || readChunk == nil {
		t.Error(err)
	}
}

func TestCodeDecodeData(t *testing.T) {
	chunk := makeOrderedChunk()
	serializedChunk, err := SerializeChunkData(chunk)
	if err != nil || serializedChunk == nil {
		t.Error(err)
	}
	readChunk, err := DeserializeChunkData(serializedChunk)
	if err != nil || readChunk == nil {
		t.Error(err)
	}

	if len(readChunk) != len(chunk.data) || len(readChunk[0]) != len(chunk.data[0]) {
		t.Error("Chunk length mismatch")
	}

	for i := 0; i < len(readChunk); i++ {
		for j := 0; j < len(readChunk[i]); j++ {
			if readChunk[i][j] != chunk.data[i][j] {
				t.Error("Chunk data mismatch")
			}
		}
	}
}
