package model

type Qris struct {
	Acquirer string `json:"acquirer"` 
}

type BankTransfer struct {
	Bank string `json:"bank"`
}

type Gopay struct{}

type ShopeePay struct{}

type EChannel struct {
	BillInfo1 string `json:"bill_info1"`
	BillInfo2 string `json:"bill_info2"`
}

type CStore struct {
	Store   string `json:"store"`
	Message string `json:"message"`
}