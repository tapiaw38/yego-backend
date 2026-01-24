package order

import (
	"context"
	"time"

	orderRepo "wappi/internal/adapters/datasources/repositories/order"
	ordertokenRepo "wappi/internal/adapters/datasources/repositories/ordertoken"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
)

// CreateWithLinkInput represents the input for creating an order with a claim link
type CreateWithLinkInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	ETA         string `json:"eta"`
}

// CreateWithLinkOutput represents the output after creating an order with link
type CreateWithLinkOutput struct {
	OrderID   string `json:"order_id"`
	Token     string `json:"token"`
	ClaimURL  string `json:"claim_url"`
	Status    string `json:"status"`
	ETA       string `json:"eta"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

// CreateWithLinkUsecase defines the interface for creating orders with claim links
type CreateWithLinkUsecase interface {
	Execute(ctx context.Context, input CreateWithLinkInput, baseURL string) (*CreateWithLinkOutput, apperrors.ApplicationError)
}

type createWithLinkUsecase struct {
	orderRepo      orderRepo.Repository
	orderTokenRepo ordertokenRepo.Repository
}

// NewCreateWithLinkUsecase creates a new instance of CreateWithLinkUsecase
func NewCreateWithLinkUsecase(orderRepo orderRepo.Repository, orderTokenRepo ordertokenRepo.Repository) CreateWithLinkUsecase {
	return &createWithLinkUsecase{
		orderRepo:      orderRepo,
		orderTokenRepo: orderTokenRepo,
	}
}

// Execute creates a new order and generates a claim link
func (u *createWithLinkUsecase) Execute(ctx context.Context, input CreateWithLinkInput, baseURL string) (*CreateWithLinkOutput, apperrors.ApplicationError) {
	// Create order without user assignment
	newOrder := &domain.Order{
		ETA: input.ETA,
	}

	created, err := u.orderRepo.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	// Create order token for claiming
	expiresAt := time.Now().Add(24 * time.Hour)
	orderToken := &domain.OrderToken{
		OrderID:     created.ID,
		PhoneNumber: &input.PhoneNumber,
		ExpiresAt:   expiresAt,
	}

	tokenCreated, err := u.orderTokenRepo.Create(ctx, orderToken)
	if err != nil {
		return nil, err
	}

	// Generate claim URL
	claimURL := baseURL + "/order/claim/" + tokenCreated.Token

	return &CreateWithLinkOutput{
		OrderID:   created.ID,
		Token:     tokenCreated.Token,
		ClaimURL:  claimURL,
		Status:    string(created.Status),
		ETA:       created.ETA,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z"),
		CreatedAt: created.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
