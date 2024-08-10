package shared

type GNode struct {
	Id    uint64
	Child Node
}

func NewGNode(id uint64) *GNode {
	n := new(GNode)
	n.Id = id
	return n
}

func (n *GNode) SetId(id uint64) {
	n.Id = id
}

func (n *GNode) GetId() uint64 {
	return n.Id
}

func (n *GNode) GetChild() Node {
	return n.Child
}
