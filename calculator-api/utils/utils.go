package utils

func SumMatrix(matrix1 [3][3]int, matrix2 [3][3]int) [3][3]int {
	type RowResult struct {
		RowIndex int
		RowData  [3]int
	}

	var res [3][3]int
	rows := 3

	/*
		The sum of each row is executed cocurrently in a goroutine.
		The result of each sum is sent to a buffered channel.
		The main goroutine consumes (receives) the channel results to rebuild the final matrix in correct order, using RowIndex.
	*/
	resultsCh := make(chan RowResult, rows)

	for i := range rows {
		row := i

		go func(rowIdx int) {
			rowRes := RowResult{RowIndex: rowIdx}

			for j := 0; j < 3; j++ {
				rowRes.RowData[j] = matrix1[rowIdx][j] + matrix2[rowIdx][j]
			}

			resultsCh <- rowRes
		}(row)
	}

	for i := 0; i < rows; i++ {
		result := <-resultsCh

		res[result.RowIndex] = result.RowData
	}

	return res
}
