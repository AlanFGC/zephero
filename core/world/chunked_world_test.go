package world

import (
	"fmt"
	"testing"
)

func makeOrderedChunk() *WorldChunk {
	arr := make([][]GNode, 10)
	for i := range arr {
		arr[i] = make([]GNode, 10)
	}

	for i := range arr {
		for j := range arr[i] {
			arr[i][j] = GNode{Eid: uint64(i*len(arr) + j)}
		}
	}

	chunk := WorldChunk{
		ChunkId: 0,
		Data:    arr,
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

func TestEncodeDecodeData(t *testing.T) {
	chunk := makeOrderedChunk()
	serializedChunk, err := SerializeChunkData(chunk)
	if err != nil || serializedChunk == nil {
		t.Error(err)
	}
	readChunk, err := DeserializeChunkData(serializedChunk)
	if err != nil || readChunk == nil {
		t.Error(err)
	}

	if len(readChunk) != len(chunk.Data) || len(readChunk[0]) != len(chunk.Data[0]) {
		t.Error("Chunk length mismatch")
	}

	for i := 0; i < len(readChunk); i++ {
		for j := 0; j < len(readChunk[i]); j++ {
			if readChunk[i][j] != chunk.Data[i][j] {
				t.Error("Chunk data mismatch")
			}
		}
	}
}

func TestChunkedWorld_GetSize(t *testing.T) {
	const chunkLen = 12
	const chunkLenV = 10
	const chunkLenH = 11
	world, err := NewChunkedWorld(chunkLenV, chunkLenH, chunkLen)
	if err != nil {
		t.Error(err)
	}
	expectedRows := chunkLenV * chunkLen
	expectedCols := chunkLenH * chunkLen
	rows, cols := world.GetSize()
	if rows != expectedRows || cols != expectedCols {
		t.Error(fmt.Sprintf("World size mismatch: rows: %d -> %d cols: %d -> %d",
			rows, expectedRows, cols, expectedCols))
	}
}

func Test_getSpace(t *testing.T) {
	const chunkLen = 12
	const chunkLenV = 10
	const chunkLenH = 11
	world, err := NewChunkedWorld(chunkLenV, chunkLenH, chunkLen)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < chunkLenV*chunkLen; i++ {
		for j := 0; j < chunkLenH*chunkLen; j++ {
			node, err := world.GetSpace(i, j)
			if err != nil {
				t.Error(err)
			}
			if node == nil {
				t.Error("Node is nil")
			}
		}
	}
}

func Test_setSpace(t *testing.T) {
	const chunkLen = 32
	const chunkLenV = 100
	const chunkLenH = 100

	// Create a new chunked world
	world, err := NewChunkedWorld(chunkLenV, chunkLenH, chunkLen)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Loop through all possible spaces and set their values
	for i := 0; i < chunkLenV*chunkLen; i++ {
		for j := 0; j < chunkLenH*chunkLen; j++ {
			err := world.SetSpace(uint64(i*chunkLen+j), uint64(j), i, j)
			if err != nil {
				t.Errorf("Error setting space at (%d, %d): %v", i, j, err)
			}
		}
	}

	rows, cols := world.GetSize()
	// Verify that the spaces were set correctly by querying them again
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			node, err := world.GetSpace(i, j)
			if err != nil {
				t.Errorf("Error querying space at (%d, %d) after setting: %v", i, j, err)
			}
			if node == nil {
				t.Errorf("Node is nil at (%d, %d) after setting", i, j)
			}
			res := node.GetId()
			if res != uint64(i*chunkLen+j) {
				t.Errorf("Wrong id at (%d, %d)", i, j)
			}
		}
	}
}

func TestChunkedWorld_GetPlayerViewByCellCoordinate(t *testing.T) {
	world, err := NewChunkedWorld(100, 100, 16)
	if err != nil {
		t.Error(err)
	}

	view, err := world.GetPlayerViewByCellCoordinate(511, 511)
	if err != nil {
		t.Error(err)
	}
	if len(view) != 9 {
		t.Error("Expected chunk array of size 9, got", len(view))
	}
	for i := 0; i < len(view); i++ {
		if view[i].ChunkId < 0 {
			t.Error("Expected chunk at index", i, view[i].ChunkId)
		}
	}
}

func TestChunkedWorld_SetChunk(t *testing.T) {
	world, err := NewChunkedWorld(10, 10, 16)
	if err != nil {
		t.Error(err)
	}
	newChunk := makeChunkData(16)

	customEid := uint64(53)

	for i := 0; i < 16; i++ {
		newChunk[i] = make([]GNode, 16)
		for j := 0; j < 16; j++ {
			newChunk[i][j] = GNode{
				Eid: customEid,
				Tid: customEid,
			}
		}
	}

	err = world.SetChunk(5, 5, newChunk)
	if err != nil {
		t.Error(err)
	}
	currentChunk, err := world.getChunkByChunkCoordinate(5, 5)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			node := currentChunk.Data[i][j]
			if node.Eid != customEid || node.Tid != customEid {
				t.Error("Unexpected entity ID")
			}
		}
	}
}
