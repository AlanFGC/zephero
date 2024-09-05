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
				data:    newChunkData(chunkLen),
			}
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
	if row < 0 || col < 0 {
		return errors.New("invalid coordinate")
	}
	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	chunk.data[row%w.chunkLen][col%w.chunkLen] = GNode{
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
	return w.rows * w.chunkLen, w.cols * w.chunkLen
}

func (w *ChunkedWorld) Save(dao *SqliteDAO) error {
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

func (w *ChunkedWorld) Load(dao *SqliteDAO) error {
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

// private functions
func (w *ChunkedWorld) getChunkFromRowCol(row int, col int) (*WorldChunk, error) {
	if row < 0 || col < 0 {
		return nil, errors.New("invalid coordinate")
	}
	chunkIndexRow := row / w.chunkLen
	chunkIndexCol := col / w.chunkLen
	chunk := w.world[chunkIndexRow][chunkIndexCol]
	return &chunk, nil
}
