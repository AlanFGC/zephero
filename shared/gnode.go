package shared

type Node interface {
	SetId(id uint64)
	GetId() uint64
	SetChild(id uint64)
	GetChild() uint64
}

type GNode struct {
	Id    uint64
	Child uint64
}

func NewGNode(id uint64) *GNode {
	n := new(GNode)
	n.Id = id
	n.Child = 0
	return n
}

func (n *GNode) SetId(id uint64) {
	n.Id = id
}

func (n *GNode) GetId() uint64 {
	return n.Id
}

func (n *GNode) GetChild() uint64 {
	return n.Child
}

func (n *GNode) SetChild(id uint64) {
	n.Child = id
}
