package shared

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math"
)

type ChunkedWorld struct {
	id        int
	world     [][]WorldChunk
	rows      int
	cols      int
	ChunkSize int
	chunkLen  int
}

func NewChunkedWorld(chunkLenV int, chunkLenH int, chunkLen int) (*ChunkedWorld, error) {
	chunkSize := int(math.Pow(float64(chunkLen), 2))
	log.Println(fmt.Sprintf("chunksize: %d", chunkSize))

	if chunkSize&3 != 0 {
		return nil, errors.New("chunk size needs to be a multiple of 4")
	}

	w := ChunkedWorld{
		world:     make([][]WorldChunk, chunkLenV),
		rows:      chunkLenV * chunkLen,
		cols:      chunkLenH * chunkLen,
		ChunkSize: chunkSize,
		chunkLen:  chunkLen,
	}

	for i := 0; i < chunkLenV; i++ {
		w.world[i] = make([]WorldChunk, chunkLenH)
	}

	var idx = 0
	for rowIndex, rowChunk := range w.world {
		for colIndex := range rowChunk {
			w.world[rowIndex][colIndex] = WorldChunk{
				ChunkId: idx,
				Data:    makeChunkData(chunkLen),
				Row:     rowIndex,
				Col:     colIndex,
			}
			idx += 1
		}
	}
	chunksCreated := len(w.world) * len(w.world[0])
	log.Println(fmt.Sprintf("Rows: %d, Cols: %d", w.rows, w.cols))
	log.Println("Total spaces created:", chunksCreated*chunkSize)
	log.Println("Total chunks created:", chunksCreated)
	return &w, nil
}

func makeChunkData(chunkLen int) [][]GNode {
	c := make([][]GNode, chunkLen)
	for i := 0; i < chunkLen; i++ {
		c[i] = make([]GNode, chunkLen)
	}
	return c
}

func (w *ChunkedWorld) SetSpace(id uint64, child uint64, row int, col int) error {
	if row < 0 || row >= w.rows || col < 0 || col >= w.cols {
		return errors.New("invalid coordinate")
	}
	chunk, err := w.getChunkByCellCoordinate(row, col)
	if err != nil {
		return err
	}

	indexRow := row % w.chunkLen
	indexCol := col % w.chunkLen

	chunk.Data[indexRow][indexCol] = GNode{
		EntityID:  id,
		TerrainID: child,
	}
	return nil
}

func (w *ChunkedWorld) GetSpace(row int, col int) (Node, error) {
	if row < 0 || col < 0 {
		return nil, errors.New("invalid coordinate")
	}
	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	node := chunk.Data[row%w.chunkLen][col%w.chunkLen]
	return &node, nil
}

func (w *ChunkedWorld) GetSize() (int, int) {
	return w.rows, w.cols
}

func (w *ChunkedWorld) SetChunk(chunkRowId int, chunkColId int, chunk [][]GNode) error {
	if chunkRowId < 0 || chunkRowId >= w.rows/w.chunkLen || chunkColId < 0 || chunkColId >= w.cols/w.chunkLen {
		return errors.New("invalid coordinate")
	}

	if len(chunk) != w.chunkLen || len(chunk[0]) != w.chunkLen {
		return errors.New("invalid chunk size")
	}
	oldId := w.world[chunkRowId][chunkColId].ChunkId
	w.world[chunkColId][chunkRowId] = WorldChunk{
		ChunkId: oldId,
		Data:    chunk,
		Row:     chunkRowId,
		Col:     chunkColId,
	}

	return nil
}

func (w *ChunkedWorld) GetChunkData() ([][]WorldChunk, error) {
	if w.world == nil {
		return nil, errors.New("world is empty")
	}

	return w.world, nil
}

func SerializeChunkData(chunk *WorldChunk) ([]byte, error) {
	data := &chunk.Data
	var buffer bytes.Buffer
	gob.Register(GNode{})
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DeserializeChunkData(data []byte) ([][]GNode, error) {
	var chunk [][]GNode
	buffer := bytes.NewBuffer(data)
	gob.Register(GNode{})
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&chunk); err != nil {
		return nil, err
	}
	return chunk, nil
}

func (w *ChunkedWorld) GetPlayerViewByCellCoordinate(row int, col int) ([]WorldChunk, error) {
	log.Println("GetPlayerViewByCellCoordinate")

	centerV := row / w.chunkLen
	centerH := col / w.chunkLen
	chunks := make([]WorldChunk, 9)

	idx := 0
	invalidChunk := WorldChunk{ChunkId: -1}
	for i := centerV - 1; i <= centerV+1; i++ {
		for j := centerH - 1; j <= centerH+1; j++ {
			if i >= 0 && i < w.rows && j >= 0 && j < w.cols {
				log.Println("Using existing chunk")
				chunks[idx] = w.world[i][j]
			} else {
				chunks[idx] = invalidChunk
			}
			idx += 1
		}
	}

	return chunks, nil
}

// private functions
func (w *ChunkedWorld) getChunkByCellCoordinate(row int, col int) (*WorldChunk, error) {
	err := w.checkCellCoordinate(row, col)
	if err != nil {
		return nil, err
	}

	if w.chunkLen == 0 {
		return nil, errors.New("chunkLen is 0")
	}

	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen

	if chunkIndexRow >= len(w.world) || chunkIndexCol >= len(w.world[chunkIndexRow]) {
		fmt.Println(w.rows, w.cols)
		fmt.Println(len(w.world[chunkIndexRow]))
		return nil, errors.New(
			fmt.Sprintf("chunk index is out of range: row=%d, col=%d, chunkIndexRow: %d chunkIndexCol: %d",
				row, col, chunkIndexRow, chunkIndexCol))
	}
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	return &chunk, nil
}

// Returns the current chunk based on the row and column of world chunk not to be confused with cell coordinates.
func (w *ChunkedWorld) getChunkByChunkCoordinate(row int, col int) (WorldChunk, error) {
	err := w.checkChunkCoordinate(row, col)
	if err != nil {
		return WorldChunk{}, err
	}
	return w.world[row][col], nil
}

// returns an error if row and col are out of bounds
func (w *ChunkedWorld) checkCellCoordinate(row int, col int) error {
	if row < 0 || row >= w.rows || col < 0 || col >= w.cols {
		return errors.New(fmt.Sprintf("Invalid coordinate: %d:%d, row:%d cols:%d", row, col,
			w.rows, w.cols))
	}
	return nil
}

func (w *ChunkedWorld) checkChunkCoordinate(row int, col int) error {
	chunkLenV := w.rows / w.chunkLen
	chunkLenH := w.cols / w.chunkLen
	if row < 0 || row > chunkLenV || col < 0 || col > chunkLenH {
		return fmt.Errorf(
			"coordinates out of bounds: row=%d, col=%d (valid row: 0-%d, valid col: 0-%d)",
			row,
			col,
			chunkLenV,
			chunkLenH)
	}
	return nil
}
