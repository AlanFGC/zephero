package main

func PrintWorld(w game.World) {
	rows, cols := w.GetSize()
	t := table.NewWriter(os.Stdout)
	strMatrix := make([][]string, rows)
	for i := 0; i < rows; i++ {
		strMatrix[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			node, _ := w.GetSpace(i, j)
			if node != nil {
				str := strconv.FormatUint(node.GetId(), 10)
				strMatrix[i][j] = string(str[0])
			} else {
				strMatrix[i][j] = " "
			}
		}
	}
	for _, row := range strMatrix {
		t.Append(row)
	}
	t.Render()
}

func generateTimeBasedID() uint64 {
	timestamp := uint64(time.Now().UnixNano())
	counterValue := atomic.AddUint64(&counter, 1)
	return (timestamp << 16) | (counterValue & 0xFFFF)
}

func chance(probability float64) bool {
	return rand.Float64() < probability
}