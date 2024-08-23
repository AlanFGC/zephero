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

func (w *chunkedWorld) Save(dao *SqliteDAO) error {
	for i, _ := range w.world {
		for j, _ := range w.world[i] {
			chunk := w.world[i][j]
			byteData, err := serializeChunk(&chunk)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err = dao.SaveWorldChunk(0, i, j, byteData, false)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}

func (w *chunkedWorld) Load(dao *SqliteDAO) error {
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
			data, err := dao.FetchWorldChunk(0, i, j)
			if err != nil {
				fmt.Println("err")
				continue
			}
			matrix, err := deserializeChunk(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			w.world[i][j] = { data: matrix}
		}
	}
	return nil
}
