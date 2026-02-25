package payment

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"yego/internal/adapters/web/middlewares"
	"yego/internal/platform/config"
)

type PaymentMethodsHandler struct {
	serviceURL string
	apiKey     string
	client     *http.Client
}

func NewPaymentMethodsHandler(cfg *config.ConfigurationService) *PaymentMethodsHandler {
	return &PaymentMethodsHandler{
		serviceURL: cfg.PaymentServiceURL,
		apiKey:     cfg.PaymentAPIKey,
		client:     &http.Client{},
	}
}

func (h *PaymentMethodsHandler) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if h.apiKey != "" {
		req.Header.Set("X-API-Key", h.apiKey)
	}
	return h.client.Do(req)
}

func (h *PaymentMethodsHandler) forward(c *gin.Context, resp *http.Response) {
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	c.Data(resp.StatusCode, "application/json", body)
}

func (h *PaymentMethodsHandler) userIDFromContext(c *gin.Context) (string, bool) {
	userID, ok := middlewares.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	return userID, ok
}

// List GET /api/payment-methods
func (h *PaymentMethodsHandler) List(c *gin.Context) {
	userID, ok := h.userIDFromContext(c)
	if !ok {
		return
	}
	url := fmt.Sprintf("%s/api/v1/payment-methods/user/%s", h.serviceURL, userID)
	resp, err := h.doRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payment service unavailable"})
		return
	}
	h.forward(c, resp)
}

// GetDefault GET /api/payment-methods/default
func (h *PaymentMethodsHandler) GetDefault(c *gin.Context) {
	userID, ok := h.userIDFromContext(c)
	if !ok {
		return
	}
	url := fmt.Sprintf("%s/api/v1/payment-methods/user/%s/default", h.serviceURL, userID)
	resp, err := h.doRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payment service unavailable"})
		return
	}
	h.forward(c, resp)
}

// Create POST /api/payment-methods
func (h *PaymentMethodsHandler) Create(c *gin.Context) {
	userID, ok := h.userIDFromContext(c)
	if !ok {
		return
	}
	url := fmt.Sprintf("%s/api/v1/payment-methods/?user_id=%s", h.serviceURL, userID)
	resp, err := h.doRequest("POST", url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payment service unavailable"})
		return
	}
	h.forward(c, resp)
}

// Update PUT /api/payment-methods/:id
func (h *PaymentMethodsHandler) Update(c *gin.Context) {
	userID, ok := h.userIDFromContext(c)
	if !ok {
		return
	}
	id := c.Param("id")
	url := fmt.Sprintf("%s/api/v1/payment-methods/%s?user_id=%s", h.serviceURL, id, userID)
	resp, err := h.doRequest("PUT", url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payment service unavailable"})
		return
	}
	h.forward(c, resp)
}

// Delete DELETE /api/payment-methods/:id
func (h *PaymentMethodsHandler) Delete(c *gin.Context) {
	userID, ok := h.userIDFromContext(c)
	if !ok {
		return
	}
	id := c.Param("id")
	url := fmt.Sprintf("%s/api/v1/payment-methods/%s?user_id=%s", h.serviceURL, id, userID)
	resp, err := h.doRequest("DELETE", url, nil)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payment service unavailable"})
		return
	}
	h.forward(c, resp)
}

// MercadoPago proxy — CreateToken POST /api/mercadopago/token
func (h *PaymentMethodsHandler) CreateToken(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/mercadopago/token", h.serviceURL)
	resp, err := h.doRequest("POST", url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payment service unavailable"})
		return
	}
	h.forward(c, resp)
}

// MercadoPago proxy — GetPaymentMethod GET /api/mercadopago/payment_method
func (h *PaymentMethodsHandler) GetPaymentMethod(c *gin.Context) {
	bin := c.Query("bin")
	url := fmt.Sprintf("%s/api/v1/mercadopago/payment_method?bin=%s", h.serviceURL, bin)
	resp, err := h.doRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payment service unavailable"})
		return
	}
	h.forward(c, resp)
}
