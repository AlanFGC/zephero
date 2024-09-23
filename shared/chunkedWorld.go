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

type ChunkedWorld struct {
	id        int
	world     [][]WorldChunk
	rows      int
	cols      int
	chunkSize int
	chunkLen  int
}

func newChunkedWorld(chunkLenV int, chunkLenH int, chunkLen int) (*ChunkedWorld, error) {
	chunkSize := int(math.Pow(float64(chunkLen), 2))
	fmt.Println("chunksize: ", chunkSize)

	if chunkSize&3 != 0 {
		return nil, errors.New("chunk size needs to be a multiple of 4")
	}

	w := ChunkedWorld{
		world:     make([][]WorldChunk, chunkLenV),
		rows:      chunkLenV * chunkSize,
		cols:      chunkLenH * chunkSize,
		chunkSize: chunkSize,
		chunkLen:  chunkLen,
	}

	for i := 0; i < chunkLenV; i++ {
		w.world[i] = make([]WorldChunk, chunkLenH)
	}
	var idx = 0
	for rowIndex, rowChunk := range w.world {
		for colIndex := range rowChunk {
			w.world[rowIndex][colIndex] = WorldChunk{
				chunkId: idx,
				data:    newChunkData(chunkLen),
			}
			idx += 1
		}
	}
	chunksCreated := len(w.world) * len(w.world[0])
	fmt.Println("Total spaces created:", chunksCreated*chunkSize)
	fmt.Println("Total chunks created:", chunksCreated)
	return &w, nil
}

func newChunkData(chunkLen int) [][]GNode {
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

	if indexRow < 0 || indexRow >= w.chunkLen || indexCol < 0 || indexCol >= w.chunkLen {
		return errors.New(fmt.Sprintf("invalid chunk coordinate row %d, col %d", indexRow, indexCol))
	}
	chunk.data[indexRow][indexCol] = GNode{
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
	node := chunk.data[row%w.chunkLen][col%w.chunkLen]
	return &node, nil
}

func (w *ChunkedWorld) GetSize() (int, int) {
	return w.rows, w.cols
}

func (w *ChunkedWorld) SaveWorld(dao *SqliteDAO) error {
	if dao == nil {
		return errors.New("dao is nil")
	}
	for i := range w.world {
		for j := range w.world[i] {
			chunk := w.world[i][j]
			byteData, err := SerializeChunkData(&chunk)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err = dao.SaveWorldChunk(w.id, i, j, byteData, false)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}

func (w *ChunkedWorld) LoadWorld(dao *SqliteDAO) error {
	if dao == nil {
		return errors.New("dao is nil")
	}
	if len(w.world) != w.rows {
		w.world = make([][]WorldChunk, w.rows)
	}

	for z := 0; z < w.rows; z++ {
		if len(w.world[z]) != w.cols {
			w.world[z] = make([]WorldChunk, w.cols)
		}
	}

	for i := 0; i < w.rows; i++ {
		for j := 0; j < w.cols; j++ {
			data, err := dao.FetchWorldChunk(w.id, i, j)
			if err != nil {
				fmt.Println("err")
				continue
			}
			matrix, err := DeserializeChunkData(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			w.world[i][j].chunkId = i*w.chunkLen + j
			w.world[i][j].data = matrix
		}
	}
	return nil
}

func SerializeChunkData(chunk *WorldChunk) ([]byte, error) {
	data := &chunk.data
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
	err := w.checkCellCoordinate(row, col)
	chunks := make([]WorldChunk, 9)
	if err != nil {
		return chunks, err
	}

	numbChunkV := w.rows / w.chunkLen
	numbChunkH := w.cols / w.chunkLen
	startV := numbChunkV/(row/w.chunkLen) - 1
	startH := numbChunkH/(col/w.chunkLen) - 1
	idx := 0
	for i := startV; i < startV+3; i++ {
		for j := startH; j < startH+3; j++ {
			if i < 0 || j < 0 {
				// chunkId -5 mean out of bounds
				chunks[idx] = WorldChunk{chunkId: -5}
				continue
			}
			chunk, err := w.getChunkByChunkCoordinate(i, j)
			if err != nil {
				fmt.Println(fmt.Sprintf("err: %v", err))
				return nil, err
			}
			chunks[idx] = chunk
		}
		idx += 1
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

	var chunkIndexRow = 0
	if row/w.chunkLen > 0 {
		chunkIndexRow = row / w.chunkLen
	}

	var chunkIndexCol = 0
	if col/w.chunkLen > 0 {
		chunkIndexCol = col / w.chunkLen
	}

	if chunkIndexRow >= len(w.world) || chunkIndexCol >= len(w.world[chunkIndexRow]) {
		return nil, errors.New(fmt.Sprintf("chunk index is out of range, chunkRow: %d, chunkCol: %d",
			chunkIndexRow, chunkIndexCol))
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
	if row < 0 || row > w.rows || col < 0 || col > w.cols {
		return errors.New(fmt.Sprintf("Invalid coordinate: %d:%d", row, col))
	}
	return nil
}

func (w *ChunkedWorld) checkChunkCoordinate(row, col int) error {
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
