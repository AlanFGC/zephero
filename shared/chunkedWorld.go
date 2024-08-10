package shared

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
)

type WorldChunk struct {
	chunkId int
	data    [][]Node
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

func newChunk(chunkLen int) [][]Node {
	c := make([][]Node, chunkLen)
	for i := 0; i < chunkLen; i++ {
		c[i] = make([]Node, chunkLen)
	}

	return c
}

func (w *chunkedWorld) SetSpace(node Node, row int, col int) error {
	if row < 0 || col < 0 {
		return errors.New("invalid coordinate")
	}
	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	chunk.data[row%w.chunkLen][col%w.chunkLen] = node
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

func (w *chunkedWorld) GetChunkFromRowCol(row int, col int) (*WorldChunk, error) {
	if row < 0 || col < 0 {
		return nil, errors.New("invalid coordinate")
	}
	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	return &chunk, nil
}

func (w *chunkedWorld) GetSize() (int, int) {
	return w.rows * w.chunkLen, w.cols * w.chunkLen
}

func (w *chunkedWorld) getEntityView(row int, col int) {

}

func (w *chunkedWorld) loadWorld(url string) {

}

func saveWorldChunks(db *sql.DB, world chunkedWorld) error {
	// Assuming the table and database connection are correctly set up
	for i := 0; i < world.rows; i++ {
		for j := 0; j < world.cols; j++ {
			chunk := world.world[i][j]
			serializedData := serializeChunk(chunk)
			if err := insertChunk(db, chunk.chunkId, i, j, serializedData); err != nil {
				return fmt.Errorf("error saving chunk at (%d, %d): %v", i, j, err)
			}
		}
	}
	return nil
}

func insertChunk(db *sql.DB, chunkId, row, col int, data string) error {
}

func serializeChunk(chunk WorldChunk) string {
}
