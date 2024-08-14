package shared

import (
	"errors"
)

type GameWorld struct {
	world [][]Node
	rows  int
	cols  int
}

func (w *GameWorld) Save(dao *SqliteDAO) error {
	//TODO implement me
	panic("implement me")
}

func (w *GameWorld) Load(dao *SqliteDAO) error {
	//TODO implement me
	panic("implement me")
}

func newGameWorld(rows int, cols int) *GameWorld {
	w := &GameWorld{
		world: make([][]Node, rows),
		cols:  cols,
		rows:  rows,
	}
	for i := range w.world {
		w.world[i] = make([]Node, cols)
	}
	return w
}
func (w *GameWorld) SetSpace(node Node, row, col int) error {
	if col < 0 || col >= w.cols || row < 0 || row >= w.rows {
		return errors.New("node position out of bounds")
	}
	w.world[row][col] = node
	return nil
}

func (w *GameWorld) GetSpace(row, col int) (Node, error) {
	if col < 0 || col >= w.cols || row < 0 || row >= w.rows {
		return nil, errors.New("node position out of bounds")
	}
	return w.world[row][col], nil
}

func (w *GameWorld) GetSize() (int, int) {
	return w.rows, w.cols
}
