package order

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"yego/internal/domain"
	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
)

// HandlePaymentWebhookInput represents the MercadoPago webhook notification
type HandlePaymentWebhookInput struct {
	// From query params (IPN format)
	PaymentID string
	Topic     string
	// From JSON body (Webhooks format)
	Type string
	Data struct {
		ID string
	}
}

type HandlePaymentWebhookUsecase interface {
	Execute(ctx context.Context, resourceID string, topic string) apperrors.ApplicationError
}

type handlePaymentWebhookUsecase struct {
	contextFactory appcontext.Factory
}

func NewHandlePaymentWebhookUsecase(contextFactory appcontext.Factory) HandlePaymentWebhookUsecase {
	return &handlePaymentWebhookUsecase{contextFactory: contextFactory}
}

func (u *handlePaymentWebhookUsecase) Execute(ctx context.Context, resourceID string, topic string) apperrors.ApplicationError {
	app := u.contextFactory()
	mpAccessToken := app.ConfigService.MPAccessToken
	if mpAccessToken == "" {
		log.Printf("Webhook: MP_ACCESS_TOKEN not configured, skipping")
		return nil
	}

	type paymentInfo struct {
		orderID     string
		amount      float64
		mpPaymentID string
	}

	var info paymentInfo

	switch topic {
	case "payment":
		pi, err := u.getPaymentInfo(resourceID, mpAccessToken)
		if err != nil {
			log.Printf("Webhook: error getting payment %s: %v", resourceID, err)
			return nil
		}
		info = paymentInfo{orderID: pi.orderID, amount: pi.amount, mpPaymentID: resourceID}

	case "merchant_order":
		pi, err := u.getMerchantOrderInfo(resourceID, mpAccessToken)
		if err != nil {
			log.Printf("Webhook: error getting merchant_order %s: %v", resourceID, err)
			return nil
		}
		info = paymentInfo{orderID: pi.orderID, amount: pi.amount, mpPaymentID: pi.mpPaymentID}

	default:
		log.Printf("Webhook: unsupported topic %s, skipping", topic)
		return nil
	}

	if info.orderID == "" {
		log.Printf("Webhook: no external_reference found for %s/%s, skipping", topic, resourceID)
		return nil
	}

	return u.confirmOrder(ctx, info.orderID, info.mpPaymentID, info.amount)
}

func (u *handlePaymentWebhookUsecase) mpGet(path string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.mercadopago.com"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MP API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

type mpPaymentResult struct {
	orderID     string
	amount      float64
	mpPaymentID string
}

func (u *handlePaymentWebhookUsecase) getPaymentInfo(paymentID string, token string) (*mpPaymentResult, error) {
	body, err := u.mpGet("/v1/payments/"+paymentID, token)
	if err != nil {
		return nil, err
	}
	var p struct {
		Status            string  `json:"status"`
		ExternalReference string  `json:"external_reference"`
		TransactionAmount float64 `json:"transaction_amount"`
	}
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	log.Printf("Webhook: payment %s status=%s external_reference=%s amount=%.2f", paymentID, p.Status, p.ExternalReference, p.TransactionAmount)
	if p.Status != "approved" {
		return &mpPaymentResult{}, nil
	}
	return &mpPaymentResult{orderID: p.ExternalReference, amount: p.TransactionAmount, mpPaymentID: paymentID}, nil
}

func (u *handlePaymentWebhookUsecase) getMerchantOrderInfo(orderID string, token string) (*mpPaymentResult, error) {
	body, err := u.mpGet("/merchant_orders/"+orderID, token)
	if err != nil {
		return nil, err
	}
	var mo struct {
		Status            string  `json:"status"`
		ExternalReference string  `json:"external_reference"`
		TotalAmount       float64 `json:"total_amount"`
		Payments          []struct {
			ID     int64   `json:"id"`
			Status string  `json:"status"`
			Amount float64 `json:"transaction_amount"`
		} `json:"payments"`
	}
	if err := json.Unmarshal(body, &mo); err != nil {
		return nil, err
	}
	log.Printf("Webhook: merchant_order %s status=%s external_reference=%s", orderID, mo.Status, mo.ExternalReference)

	var approvedPaymentID string
	var approvedAmount float64
	for _, p := range mo.Payments {
		if p.Status == "approved" {
			approvedPaymentID = fmt.Sprintf("%d", p.ID)
			approvedAmount = p.Amount
			break
		}
	}
	if approvedPaymentID == "" && mo.Status != "closed" {
		return &mpPaymentResult{}, nil
	}
	amount := approvedAmount
	if amount == 0 {
		amount = mo.TotalAmount
	}
	return &mpPaymentResult{orderID: mo.ExternalReference, amount: amount, mpPaymentID: approvedPaymentID}, nil
}

func (u *handlePaymentWebhookUsecase) confirmOrder(ctx context.Context, orderID string, mpPaymentID string, amount float64) apperrors.ApplicationError {
	app := u.contextFactory()

	order, appErr := app.Repositories.Order.GetByID(ctx, orderID)
	if appErr != nil {
		log.Printf("Webhook: order %s not found: %v", orderID, appErr)
		return nil
	}

	if order.Status != domain.StatusCreated {
		log.Printf("Webhook: order %s already in status %s, skipping", orderID, order.Status)
		return nil
	}

	_, appErr = app.Repositories.Order.UpdateStatus(ctx, orderID, domain.StatusConfirmed)
	if appErr != nil {
		return appErr
	}

	userID := ""
	if order.UserID != nil {
		userID = *order.UserID
	}
	description := fmt.Sprintf("Pago por link para pedido %s", orderID)
	transaction := &domain.Transaction{
		OrderID:          orderID,
		UserID:           userID,
		ProfileID:        order.ProfileID,
		Amount:           amount,
		Currency:         "ARS",
		Status:           "approved",
		GatewayPaymentID: &mpPaymentID,
		Description:      &description,
	}
	if _, transErr := app.Repositories.Transaction.Create(ctx, transaction); transErr != nil {
		log.Printf("Webhook: warning: failed to create transaction for order %s: %v", orderID, transErr)
	}

	log.Printf("Webhook: order %s confirmed via payment link (mp_payment %s amount=%.2f)", orderID, mpPaymentID, amount)
	return nil
}
