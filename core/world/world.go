package world

type WorldChunk struct {
	ChunkId int       `json:"chunk_id"`
	Data    [][]GNode `json:"data"`
	Row     int
	Col     int
}

type World interface {
	SetSpace(id uint64, child uint64, row, col int) error
	GetSpace(row, col int) (Node, error)
	GetSize() (int, int)
	GetChunkData() ([][]WorldChunk, error)
	SetChunk(chunkRowId int, chunkColId int, chunk [][]GNode) error
}

func NewWorld(chunkLenVertical int, chunkLenHorizontal int, chunkSideLen int) (World, error) {
	return NewChunkedWorld(chunkLenVertical, chunkLenHorizontal, chunkSideLen)
}
