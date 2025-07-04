package dto

type TransactionCreateRequest struct {
	Transaction struct {
		CustomerId string  `json:"customerId" binding:"required,uuid"`
		Notes      string  `json:"notes"`
	} `json:"transaction" binding:"required"`
	Details []struct {
		ProductId      string  `json:"productId" binding:"required,uuid"`
		ProductPrice   float64 `json:"productPrice" binding:"required,gt=0"`
		Qty            int     `json:"qty" binding:"required,gt=0"`
		DiscountAmount float64 `json:"discountAmount"`
		Notes          string  `json:"notes"`
	} `json:"details" binding:"required,min=1"`
	PaymentDetails map[string]interface{} `json:"payment_details"`
}

type TransactionDetailResponse struct {
	ID             string  `json:"id"`
	ProductId      string  `json:"productId"`
	ProductPrice   float64 `json:"productPrice"`
	Qty            int     `json:"qty"`
	DiscountAmount float64 `json:"discountAmount"`
	Subtotal       float64 `json:"subtotal"`
	Status         string  `json:"status"`
	Notes          string  `json:"notes,omitempty"`
}

type TransactionResponse struct {
	ID            string                      `json:"id"`
	CustomerId    string                      `json:"customerId"`
	TotalAmount   float64                     `json:"totalAmount"`
	PaymentStatus string                      `json:"paymentStatus"`
	PaymentMethod string                      `json:"paymentMethod,omitempty"`
	PaymentURL    string                      `json:"paymentUrl,omitempty"`
	Status        string                      `json:"status"`
	Notes         string                      `json:"notes,omitempty"`
	Details       []TransactionDetailResponse `json:"details"`
}

type MidtransReq struct {
	
}
