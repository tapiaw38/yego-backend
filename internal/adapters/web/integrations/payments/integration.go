package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"wappi/internal/platform/config"
)

type Integration interface {
	HasPaymentMethod(userID string) (bool, error)
	GetDefaultPaymentMethod(userID string) (*PaymentMethod, error)
	ProcessPaymentWithSavedMethod(userID string, amount float64, description string, externalReference string, payerEmail string, collectorID string, securityCode string) (*ProcessPaymentResponse, error)
	CreatePreference(items []PreferenceItem, payerEmail string, externalReference string, backURLSuccess string, backURLFailure string, backURLPending string, notificationURL string) (*PreferenceResponse, error)
}

type ProcessPaymentResponse struct {
	PaymentID        int    `json:"id"`
	GatewayPaymentID string `json:"gateway_payment_id"`
	Status           string `json:"status"`
}

type PreferenceItem struct {
	Title      string  `json:"title"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	CurrencyID string  `json:"currency_id"`
}

type PreferenceResponse struct {
	PreferenceID    string `json:"preference_id"`
	InitPoint       string `json:"init_point"`
	SandboxInitPoint string `json:"sandbox_init_point"`
}

type integration struct {
	baseURL string
	client  *http.Client
}

type PaymentMethod struct {
	ID              int    `json:"id"`
	UserID          string `json:"user_id"`
	LastFourDigits  string `json:"last_four_digits"`
	PaymentMethodID string `json:"payment_method_id"`
	CardholderName  string `json:"cardholder_name"`
	IsDefault       bool   `json:"is_default"`
}

func NewIntegration(cfg *config.ConfigurationService) Integration {
	baseURL := cfg.PaymentServiceURL
	if baseURL == "" {
		baseURL = "http://payments-api:8008"
	}

	return &integration{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (i *integration) HasPaymentMethod(userID string) (bool, error) {
	// First try to get default payment method
	url := fmt.Sprintf("%s/api/v1/payment-methods/user/%s/default", i.baseURL, userID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var paymentMethod PaymentMethod
		if err := json.NewDecoder(resp.Body).Decode(&paymentMethod); err != nil {
			return false, err
		}
		return paymentMethod.ID > 0, nil
	}

	// If no default, check if user has any payment methods
	if resp.StatusCode == http.StatusNotFound {
		allMethodsURL := fmt.Sprintf("%s/api/v1/payment-methods/user/%s", i.baseURL, userID)
		req2, err := http.NewRequest("GET", allMethodsURL, nil)
		if err != nil {
			return false, err
		}

		resp2, err := i.client.Do(req2)
		if err != nil {
			return false, err
		}
		defer resp2.Body.Close()

		if resp2.StatusCode == http.StatusOK {
			var paymentMethods []PaymentMethod
			if err := json.NewDecoder(resp2.Body).Decode(&paymentMethods); err != nil {
				return false, err
			}
			return len(paymentMethods) > 0, nil
		}
		return false, nil
	}

	body, _ := io.ReadAll(resp.Body)
	return false, fmt.Errorf("payment service error: %s", string(body))
}

func (i *integration) GetDefaultPaymentMethod(userID string) (*PaymentMethod, error) {
	url := fmt.Sprintf("%s/api/v1/payment-methods/user/%s/default", i.baseURL, userID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("payment service error: %s", string(body))
	}

	var paymentMethod PaymentMethod
	if err := json.NewDecoder(resp.Body).Decode(&paymentMethod); err != nil {
		return nil, err
	}

	return &paymentMethod, nil
}

func (i *integration) CreatePreference(items []PreferenceItem, payerEmail string, externalReference string, backURLSuccess string, backURLFailure string, backURLPending string, notificationURL string) (*PreferenceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments/preferences", i.baseURL)

	type backURLsPayload struct {
		Success string `json:"success"`
		Failure string `json:"failure"`
		Pending string `json:"pending"`
	}
	type preferencePayload struct {
		Items             []PreferenceItem `json:"items"`
		PayerEmail        string          `json:"payer_email"`
		ExternalReference string          `json:"external_reference"`
		BackURLs          backURLsPayload `json:"back_urls"`
		NotificationURL   string          `json:"notification_url,omitempty"`
	}

	payload := preferencePayload{
		Items:             items,
		PayerEmail:        payerEmail,
		ExternalReference: externalReference,
		BackURLs: backURLsPayload{
			Success: backURLSuccess,
			Failure: backURLFailure,
			Pending: backURLPending,
		},
		NotificationURL: notificationURL,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("payment service error: %s", string(body))
	}

	var prefResp PreferenceResponse
	if err := json.NewDecoder(resp.Body).Decode(&prefResp); err != nil {
		return nil, err
	}

	return &prefResp, nil
}

func (i *integration) ProcessPaymentWithSavedMethod(userID string, amount float64, description string, externalReference string, payerEmail string, collectorID string, securityCode string) (*ProcessPaymentResponse, error) {
	paymentMethod, err := i.GetDefaultPaymentMethod(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	if paymentMethod == nil {
		return nil, fmt.Errorf("no payment method found for user")
	}

	url := fmt.Sprintf("%s/api/v1/payments/with-saved-method", i.baseURL)

	payload := map[string]interface{}{
		"transaction_amount": amount,
		"payment_method_id":  paymentMethod.ID,
		"payer": map[string]string{
			"email": payerEmail,
		},
		"installments":      1,
		"description":       description,
		"external_reference": externalReference,
		"user_id":           userID,
	}
	if collectorID != "" {
		payload["collector_id"] = collectorID
	}
	if securityCode != "" {
		payload["security_code"] = securityCode
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("payment service error: %s", string(body))
	}

	var paymentResponse ProcessPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		return nil, err
	}

	return &paymentResponse, nil
}
