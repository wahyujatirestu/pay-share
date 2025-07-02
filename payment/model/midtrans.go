package model


type MidtransRequest struct {
	TransactionDetails struct {
		OrderId		string		`json:"order_id"`
		GrossAmt	*float64	`json:"gross_amount"`
	}	`json:"transaction_details"`
	CustomerDetails struct {
		Name 		string		`json:"first_name"`
		Email 		string		`json:"email"`
		Phone		string		`json:"phone"`
	}	`json:"customer_details"`			
}

type MidtransResponse struct {
	Token		string	`json:"token"`
	RedirectUrl	string	`json:"redirectUrl"`
	StatusCode	string	`json:"statusCode"`
	StatusMsg	string	`json:"statusMsg"`
}