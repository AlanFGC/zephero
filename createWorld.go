package main

import (
	"errors"
	world "zephero/shared"
	utils "zephero/utils"
)


func main() {
	dao := world.NewSqliteDAO("world")
	err := dao.OpenDb("world.db")
	if err != nil {
		errors.New("Failed to open DB")
		return
	}
	var w world.World
	w, err = world.NewChunkedWorld(100, 100, 16)
	if err != nil {
		errors.New("	Failed to create new world")
	}
	setRandomUUIDs(w)
	w.Save(dao)
	dao.CloseDb()
}

func setRandomUUIDs(w world.World) {
	rows, cols := w.GetSize()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if utils.Chance(0.30) {
				id := utils.GenerateTimeBasedID()
				w.SetSpace(id, 0, i, j)
			}
		}
	}
}
