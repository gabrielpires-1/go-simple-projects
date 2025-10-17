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

type sumPayload []int64
