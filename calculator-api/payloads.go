package main

type addPayload struct {
	N1 int `json:"number1"`
	N2 int `json:"number2"`
}

type multiplyPayload struct {
	N1 int `json:"number1"`
	N2 int `json:"number2"`
}

type divisionPayload struct {
	N1 int `json:"dividend"`
	N2 int `json:"divisor"`
}

type matrixSumPayload struct {
	Matrix1 [3][3]int `json:"matrix1"`
	Matrix2 [3][3]int `json:"matrix2"`
}

type matrixResult struct {
	Matrix [3][3]int `json:"result_matrix"`
}

type sumPayload []int64
