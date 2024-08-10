package shared

type Node interface {
	SetId(id uint64)
	GetId() uint64
	GetChild() Node
}

type World interface {
	SetSpace(node Node, row, col int) error
	GetSpace(row, col int) (Node, error)
	GetSize() (int, int)
}

func NewSimpleWorld(cols int, rows int) World {
	return newGameWorld(cols, rows)
}

func NewChunkedWorld(chunkLenVertical int, chunkLenHorizontal int, chunkSideLen int) (World, error) {
	return newChunkedWorld(chunkLenVertical, chunkLenHorizontal, chunkSideLen)
}
