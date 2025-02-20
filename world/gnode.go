package world

type Node interface {
	SetId(id uint64)
	GetId() uint64
	SetTerrainId(id uint64)
	GetTerrainId() uint64
}

type GNode struct {
	Eid uint64 `json:"entity_id"`
	Tid uint64 `json:"terrain_id"`
}

func NewGNode(id uint64) *GNode {
	n := new(GNode)
	n.Eid = id
	n.Tid = 0
	return n
}

func (n *GNode) SetId(id uint64) {
	n.Eid = id
}

func (n *GNode) GetId() uint64 {
	return n.Eid
}

func (n *GNode) GetTerrainId() uint64 {
	return n.Tid
}

func (n *GNode) SetTerrainId(id uint64) {
	n.Tid = id
}
