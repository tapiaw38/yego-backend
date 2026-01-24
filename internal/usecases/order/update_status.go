package order

import (
	"context"

	"wappi/internal/adapters/datasources/repositories/order"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"

	"github.com/google/uuid"
)

// UpdateStatusInput represents the input for updating order status
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatusOutput represents the output after updating order status
type UpdateStatusOutput struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	StatusIndex int    `json:"status_index"`
	ETA         string `json:"eta"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UpdateStatusUsecase defines the interface for updating order status
type UpdateStatusUsecase interface {
	Execute(ctx context.Context, id string, input UpdateStatusInput) (*UpdateStatusOutput, apperrors.ApplicationError)
}

type updateStatusUsecase struct {
	repo order.Repository
}

// NewUpdateStatusUsecase creates a new instance of UpdateStatusUsecase
func NewUpdateStatusUsecase(repo order.Repository) UpdateStatusUsecase {
	return &updateStatusUsecase{repo: repo}
}

// Execute updates the status of an order
func (u *updateStatusUsecase) Execute(ctx context.Context, id string, input UpdateStatusInput) (*UpdateStatusOutput, apperrors.ApplicationError) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderInvalidIDError, err)
	}

	if !domain.IsValidStatus(input.Status) {
		return nil, apperrors.NewApplicationError(mappings.OrderInvalidStatusError, nil)
	}

	updated, err := u.repo.UpdateStatus(ctx, id, domain.OrderStatus(input.Status))
	if err != nil {
		return nil, err
	}

	return &UpdateStatusOutput{
		ID:          updated.ID,
		Status:      string(updated.Status),
		StatusIndex: updated.StatusIndex(),
		ETA:         updated.ETA,
		CreatedAt:   updated.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   updated.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
