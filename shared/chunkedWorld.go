package shared

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math"
)

type WorldChunk struct {
	chunkId int
	data    [][]GNode
}

type chunkedWorld struct {
	world     [][]WorldChunk
	rows      int
	cols      int
	chunkSize int
	chunkLen  int
}

func newChunkedWorld(chunkLenV int, chunkLenH int, chunkLen int) (*chunkedWorld, error) {
	chunkSize := int(math.Pow(float64(chunkLen), 2))
	fmt.Println("chunksize: ", chunkSize)

	if chunkSize&3 != 0 {
		return nil, errors.New("chunk size needs to be a multiple of 4")
	}

	w := chunkedWorld{
		world:     make([][]WorldChunk, chunkLenV),
		rows:      chunkLenV,
		cols:      chunkLenH,
		chunkSize: chunkSize,
		chunkLen:  chunkLen,
	}

	for i := 0; i < chunkLenV; i++ {
		w.world[i] = make([]WorldChunk, chunkLenH)
	}

	for rowIndex, rowChunk := range w.world {
		for colIndex := range rowChunk {
			w.world[rowIndex][colIndex] = WorldChunk{
				chunkId: rowIndex*chunkLen + colIndex,
				data:    newChunk(chunkLen),
			}
		}
	}
	chunksCreated := len(w.world) * len(w.world[0])
	fmt.Println("Total spaces created:", chunksCreated*chunkSize)
	fmt.Println("Total chunks created:", chunksCreated)
	return &w, nil
}

func newChunk(chunkLen int) [][]GNode {
	c := make([][]GNode, chunkLen)
	for i := 0; i < chunkLen; i++ {
		c[i] = make([]GNode, chunkLen)
	}
	return c
}

func (w *chunkedWorld) SetSpace(id uint64, child uint64, row int, col int) error {
	if row < 0 || col < 0 {
		return errors.New("invalid coordinate")
	}
	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	chunk.data[row%w.chunkLen][col%w.chunkLen] = GNode{
		Id    uint64
		Child uint64
	}
	return nil
}

func (w *chunkedWorld) GetSpace(row int, col int) (Node, error) {
	if row < 0 || col < 0 {
		return nil, errors.New("invalid coordinate")
	}
	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	return chunk.data[row%w.chunkLen][col%w.chunkLen], nil
}

func (w *chunkedWorld) GetSize() (int, int) {
	return w.rows * w.chunkLen, w.cols * w.chunkLen
}

// private functions

func (w *chunkedWorld) getChunkFromRowCol(row int, col int) (*WorldChunk, error) {
	if row < 0 || col < 0 {
		return nil, errors.New("invalid coordinate")
	}
	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	return &chunk, nil
}


func (w *chunkedWorld) getEntityView(row int, col int) {
	panic("Not implemented")
	return
}

func (w *chunkedWorld) loadWorld(url string) {
	panic("Not implemented")
	return
}

func serializeChunk(chunk *WorldChunk) ([]byte, error) {
	data := &chunk.data
	var buffer bytes.Buffer
	gob.Register(GNode{})
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func deserializeChunk(data []byte) ([][]GNode, error) {
	var chunk [][]GNode
	buffer := bytes.NewBuffer(data)
	gob.Register(GNode{})
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(chunk); err != nil {
		return nil, err
	}
	return chunk, nil
}
