package shared

type Node interface {
	SetId(id uint64)
	GetId() uint64
	SetTerrainId(id uint64)
	GetTerrainId() uint64
}

type GNode struct {
	EntityID  uint64
	TerrainID uint64
}

func NewGNode(id uint64) *GNode {
	n := new(GNode)
	n.EntityID = id
	n.TerrainID = 0
	return n
}

func (n *GNode) SetId(id uint64) {
	n.EntityID = id
}

func (n *GNode) GetId() uint64 {
	return n.EntityID
}

func (n *GNode) GetTerrainId() uint64 {
	return n.TerrainID
}

func (n *GNode) SetTerrainId(id uint64) {
	n.TerrainID = id
}
