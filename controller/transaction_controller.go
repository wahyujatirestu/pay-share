package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/dto"
	"github.com/wahyujatirestu/payshare/model"
	// payMod "github.com/wahyujatirestu/payshare/payment/model"
	"github.com/wahyujatirestu/payshare/service"
)

type TransactionsController struct {
	tsService	service.TransactionService
}

func NewTransactionsController(tsService service.TransactionService) *TransactionsController {
	return &TransactionsController{tsService: tsService}
}

func (c *TransactionsController) Create(ctx *gin.Context)  {
	var req dto.TransactionCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ts := model.Transactions{
		CustomerId: uuid.MustParse(req.Transaction.CustomerId),
		Notes: 		&req.Transaction.Notes,
	}

	var details []*model.TransactionDetails
	for _, d := range req.Details {
		d := &model.TransactionDetails{
			ProductId: 		uuid.MustParse(d.ProductId),
			ProductPrice: 	d.ProductPrice,
			Qty: 			d.Qty,
			DiscountAmount: d.DiscountAmount,
			Notes: 			&d.Notes,
		}
		details = append(details, d)
	}


	chargeRes, err := c.tsService.Create(&ts, details)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res := dto.TransactionResponse{
		ID: ts.ID.String(),
		CustomerId: ts.CustomerId.String(),
		TotalAmount: ts.TotalAmount,
		PaymentStatus: ts.PaymentStatus,
		PaymentMethod: "",
		PaymentURL: "",
		Status: ts.Status,
		Notes: "",
		Details: []dto.TransactionDetailResponse{},
	}

	if ts.PaymentMethod != nil {
		res.PaymentMethod = *ts.PaymentMethod
	}
	if ts.PaymentURL != nil {
		res.PaymentURL = *ts.PaymentURL
	}
	if ts.Notes != nil {
		res.Notes = *ts.Notes
	}

	detailsModel, err := c.tsService.GetDetails(ts.ID.String())
	if err == nil {
		for _, d := range detailsModel{
			res.Details = append(res.Details, dto.TransactionDetailResponse{
				ID: 	d.ID.String(),
				ProductId: d.ProductId.String(),
				ProductPrice: d.ProductPrice,
				Qty: d.Qty,
				DiscountAmount: d.DiscountAmount,
				Subtotal: d.Subtotal,
				Status: d.Status,
				Notes: "",
			})
		}
	}

	ctx.JSON(201, gin.H{
		"transaction": res,
		"paymentURL" : chargeRes.RedirectUrl,
		"message": "Transaction created successfully",
	})
}

func (c *TransactionsController) GetById(ctx *gin.Context)  {
	id := ctx.Param("id")
	ts, err := c.tsService.GetById(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if ts == nil {
		ctx.JSON(404, gin.H{"error": "Transaction not found"})
	}

	details, err := c.tsService.GetDetails(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res := dto.TransactionResponse{
		ID: ts.ID.String(),
		CustomerId: ts.CustomerId.String(),
		TotalAmount: ts.TotalAmount,
		PaymentStatus: ts.PaymentStatus,
		PaymentMethod: "",
		PaymentURL: "",
		Status: ts.Status,
		Notes: "",
		Details: []dto.TransactionDetailResponse{},
	}

	if ts.PaymentMethod != nil {
		res.PaymentMethod = *ts.PaymentMethod
	}
	if ts.PaymentURL != nil {
		res.PaymentURL = *ts.PaymentURL
	}
	if ts.Notes != nil {
		res.Notes = *ts.Notes
	}

	for _, d := range details{
		res.Details = append(res.Details, dto.TransactionDetailResponse{
			ID: d.ID.String(),
			ProductId: d.ProductId.String(),
			ProductPrice: d.ProductPrice,
			Qty: d.Qty,
			DiscountAmount: d.DiscountAmount,
			Subtotal: d.Subtotal,
			Status: d.Status,
			Notes: "",
		})
	}

	ctx.JSON(200, gin.H{"transaction": res})
}

func (c *TransactionsController) GetAll(ctx *gin.Context) {
	filters := make(map[string]interface{})

	if customerId := ctx.Query("customer_id"); customerId != "" {
		filters["customer_id"] = customerId
	}
	if paymentStatus := ctx.Query("payment_status"); paymentStatus != "" {
		filters["payment_status"] = paymentStatus
	}
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}

	ts, err := c.tsService.GetAll(filters)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var res []dto.TransactionResponse
	for _, t := range ts {
		res = append(res, dto.TransactionResponse{
			ID: t.ID.String(),
			CustomerId: t.CustomerId.String(),
			TotalAmount: t.TotalAmount,
			PaymentStatus: t.PaymentStatus,
			PaymentMethod: "",
			PaymentURL: "",
			Status: t.Status,
			Notes: "",
		})
	}
	ctx.JSON(200, gin.H{"transactions": res})
}

func (c *TransactionsController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.TransactionResponse

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transaction := &model.Transactions{
		ID:        id,
		PaymentStatus: req.PaymentStatus,
		Status: req.Status,
		Notes: &req.Notes,
		Updated_At: time.Now(),
	}

	if err := c.tsService.Update(transaction); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Transaction updated successfully"})
}

func (c *TransactionsController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.tsService.Delete(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Transaction deleted Successfully"})
}


func (c *TransactionsController) WebHook(ctx *gin.Context)  {
	var body map[string]interface{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	orderID := body["order_id"].(string)
	transactionStatus := body["transaction_status"].(string)

	c.tsService.UpdateStatusFromWebhook(orderID, transactionStatus)

	ctx.JSON(200, gin.H{"message": "Webhook received successfully"})
}