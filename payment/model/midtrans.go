package model


type MidtransRequest struct {
	TransactionDetails struct {
		OrderId		string		`json:"orderId"`
		GrossAmt	*float64	`json:"gross_amount"`
	}	`json:"transactionDetails"`
	CustomerDetails struct {
		Name 		string		`json:"name"`
		Email 		string		`json:"email"`
		Phone		string		`json:"phone"`
	}	`json:"customerDetails"`
	PaymentType		string			`json:"paymentType"`
	Qris			*Qris			`json:"qris,omitempty"`
	BankTransfer 	*BankTransfer	`json:"bankTransfer,omitempty"`	
}

type MidtransResponse struct {
	Token		string	`json:"token"`
	RedirectUrl	string	`json:"redirectUrl"`
	StatusCode	string	`json:"statusCode"`
	StatusMsg	string	`json:"statusMsg"`
}