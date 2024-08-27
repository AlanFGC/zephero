package shared

type World interface {
	SetSpace(id uint64, child uint64, row, col int) error
	GetSpace(row, col int) (Node, error)
	GetSize() (int, int)
	Save(*SqliteDAO) error
	Load(*SqliteDAO) error
}

func NewChunkedWorld(chunkLenVertical int, chunkLenHorizontal int, chunkSideLen int) (World, error) {
	return newChunkedWorld(chunkLenVertical, chunkLenHorizontal, chunkSideLen)
}
