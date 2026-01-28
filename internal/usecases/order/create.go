package order

import (
	"context"

	"wappi/internal/adapters/datasources/repositories/order"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
)

// CreateInput represents the input for creating an order
type CreateInput struct {
	ProfileID string `json:"profile_id" binding:"required"`
	ETA       string `json:"eta"`
}

// CreateOutput represents the output after creating an order
type CreateOutput struct {
	ID        string  `json:"id"`
	ProfileID *string `json:"profile_id,omitempty"`
	UserID    *string `json:"user_id,omitempty"`
	Status    string  `json:"status"`
	ETA       string  `json:"eta"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// CreateUsecase defines the interface for creating orders
type CreateUsecase interface {
	Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError)
}

type createUsecase struct {
	repo order.Repository
}

// NewCreateUsecase creates a new instance of CreateUsecase
func NewCreateUsecase(repo order.Repository) CreateUsecase {
	return &createUsecase{repo: repo}
}

// Execute creates a new order
func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	newOrder := &domain.Order{
		ProfileID: &input.ProfileID,
		ETA:       input.ETA,
	}

	created, err := u.repo.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	return &CreateOutput{
		ID:        created.ID,
		ProfileID: created.ProfileID,
		UserID:    created.UserID,
		Status:    string(created.Status),
		ETA:       created.ETA,
		CreatedAt: created.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: created.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
