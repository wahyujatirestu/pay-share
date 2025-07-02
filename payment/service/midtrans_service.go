package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/wahyujatirestu/payshare/payment/model"
)

type MidtransService interface {
	Pay(payload *model.MidtransRequest) (*model.MidtransResponse, error)
}

type midtransService struct {
		Client		*resty.Client
		Url			string
}

func NewMidtransService() MidtransService {
	client := resty.New()
	url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	return &midtransService{Client: client, Url: url}
}

func (m *midtransService) Pay(payload *model.MidtransRequest) (*model.MidtransResponse, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return &model.MidtransResponse{}, errors.New("Servey key is empty")
	}

	encodedKey := base64.StdEncoding.EncodeToString([]byte(serverKey))

	payloadJSON, _ := json.MarshalIndent(payload, "", "  ")
    fmt.Println("🔍 Midtrans Request Payload:")
    fmt.Println(string(payloadJSON))

	res, err := m.Client.R().SetHeader("Authorization", "Basic " + encodedKey).SetHeader("Content-Type", "application/json").SetBody(payload).Post(m.Url)
	if err != nil {
		return &model.MidtransResponse{}, err
	}

	if res.StatusCode() != http.StatusCreated && res.StatusCode() != http.StatusOK {
		return &model.MidtransResponse{}, errors.New("Midtrans payment failed with status: " + res.Status())
	}

	var response model.MidtransResponse
	if err = json.Unmarshal(res.Body(), &response); err != nil {
		return nil, err
	}

	response.RedirectUrl = "https://app.sandbox.midtrans.com/snap/v2/vtweb/" + response.Token

	return &response, nil
}