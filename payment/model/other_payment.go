package model

type Qris struct {
	Acquirer string `json:"acquirer"` 
}

type BankTransfer struct {
	Bank string `json:"bank"`
}