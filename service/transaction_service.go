package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/model"
	payMod "github.com/wahyujatirestu/payshare/payment/model"
	"github.com/wahyujatirestu/payshare/payment/service"
	"github.com/wahyujatirestu/payshare/repository"
	"github.com/wahyujatirestu/payshare/utils"
)

type TransactionService interface {
	Create(transaction *model.Transactions, details []*model.TransactionDetails) (*payMod.MidtransResponse, error)
	GetById(id string)(*model.Transactions, error)
	GetAll(filters map[string]interface{})([]*model.Transactions, error)
	GetDetails(transactionId string)([]*model.TransactionDetails, error)
	Update(transaction *model.Transactions) error
	Delete(id string) error
	UpdateStatusFromWebhook(orderID string, status string) error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	transactionDetailsRepo repository.TransactionDetailsRespository
	userRepo	repository.UserRepository
	midtransService service.MidtransService
}

func NewTransactionService(tr repository.TransactionRepository, td repository.TransactionDetailsRespository, ur repository.UserRepository, midtransSvc service.MidtransService) TransactionService {
	return &transactionService{
		transactionRepo: tr,
		transactionDetailsRepo: td,
		userRepo: ur,
		midtransService: midtransSvc,
	}
}

func (s *transactionService) Create(ts *model.Transactions, details []*model.TransactionDetails) (*payMod.MidtransResponse, error) {
	if ts.CustomerId == uuid.Nil {
		return nil, errors.New("customer id is required")
	}
	if len(details) == 0 {
		return nil, errors.New("transaction details is required")
	}

	var total float64
	for _, d := range details {
		d.Subtotal = d.ProductPrice*float64(d.Qty) - d.DiscountAmount
		if d.Subtotal < 0 {
			d.Subtotal = 0
		}
		d.Status = "pending"
		total += d.Subtotal
	}

	ts.TotalAmount = total
	ts.PaymentStatus = "pending"
	ts.Status = "created"
	now := time.Now()
	ts.EntryDate = &now
	ts.ID = uuid.New()
	ts.Created_At = now
	ts.Updated_At = now
	ts.PaymentMethod = utils.PtrString("midtrans_snap")

	err := s.transactionRepo.Create(ts)
	if err != nil {
		return nil, err
	}

	for _, d := range details {
		d.TransactionId = ts.ID
		d.ID = uuid.New()
		d.CreatedAt = now
		d.UpdatedAt = now

		err := s.transactionDetailsRepo.Create(d)
		if err != nil {
			return nil, err
		}
	}

	customer, err := s.userRepo.GetById(ts.CustomerId.String())
	if err != nil || customer == nil {
		return nil, errors.New("Failed to get customer data for midtrans")
	}

	payload := &payMod.MidtransRequest{}
	payload.TransactionDetails.OrderId = ts.ID.String()
	payload.TransactionDetails.GrossAmt = &ts.TotalAmount
	payload.CustomerDetails.Name = customer.Name
	payload.CustomerDetails.Email = customer.Email
	payload.CustomerDetails.Phone = customer.Phone

	
	chargeRes, err := 	s.midtransService.Pay(payload)
	if err != nil {
		return nil, err
	} 

	ts.PaymentURL = &chargeRes.RedirectUrl
	ts.Status = "pending_payment"
	ts.Updated_At = time.Now()

	if err := s.transactionRepo.Update(ts); err != nil {
		return nil, err
	}

	return chargeRes, nil
}

func (s *transactionService) GetById(id string) (*model.Transactions, error) {
	return s.transactionRepo.GetById(id)
}

func (s *transactionService) GetAll(filters map[string]interface{}) ([]*model.Transactions, error) {
	return s.transactionRepo.GetAll(filters)
}

func (s *transactionService) Delete(id string) error {
	return s.transactionRepo.Delete(id)
}

func (s *transactionService) Update(ts *model.Transactions) error {
	ts.Updated_At = time.Now()
	return s.transactionRepo.Update(ts)

}

func (s *transactionService) GetDetails(transactionId string) ([]*model.TransactionDetails, error) {
	return s.transactionDetailsRepo.GetByTransactionId(transactionId)
}

func (s *transactionService) UpdateStatusFromWebhook(orderID string, status string) error {
	tx, err := s.transactionRepo.GetById(orderID)
	if err != nil || tx == nil {
		return fmt.Errorf("Transaction not found: %v", orderID)
	}

	tx.PaymentStatus = status
	tx.Updated_At = time.Now()

	return s.transactionRepo.Update(tx)
}