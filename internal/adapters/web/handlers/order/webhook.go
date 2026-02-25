package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orderUsecase "yego/internal/usecases/order"
)

// NewPaymentWebhookHandler handles MercadoPago payment notifications (IPN + Webhooks format)
func NewPaymentWebhookHandler(usecase orderUsecase.HandlePaymentWebhookUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// MP sends different formats:
		// IPN v1:        ?id=123&topic=payment
		// IPN v2/Webhooks: ?data.id=123&type=payment  OR JSON body
		topic := c.Query("topic")
		if topic == "" {
			topic = c.Query("type")
		}
		paymentID := c.Query("id")
		if paymentID == "" {
			paymentID = c.Query("data.id")
		}

		// JSON body fallback: {"type": "payment", "data": {"id": "123"}}
		if paymentID == "" || topic == "" {
			var body struct {
				Type string `json:"type"`
				Data struct {
					ID string `json:"id"`
				} `json:"data"`
			}
			if err := c.ShouldBindJSON(&body); err == nil {
				if topic == "" {
					topic = body.Type
				}
				if paymentID == "" {
					paymentID = body.Data.ID
				}
			}
		}

		if (topic != "payment" && topic != "merchant_order") || paymentID == "" {
			c.JSON(http.StatusOK, gin.H{"ok": true})
			return
		}

		if appErr := usecase.Execute(c.Request.Context(), paymentID, topic); appErr != nil {
			appErr.Log(c)
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
