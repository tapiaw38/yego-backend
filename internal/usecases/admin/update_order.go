package admin

import (
	"context"

	orderRepo "wappi/internal/adapters/datasources/repositories/order"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"

	"github.com/google/uuid"
)

// UpdateOrderInput represents the input for updating an order
type UpdateOrderInput struct {
	Status *string `json:"status,omitempty"`
	ETA    *string `json:"eta,omitempty"`
}

// UpdateOrderUsecase defines the interface for updating orders
type UpdateOrderUsecase interface {
	Execute(ctx context.Context, id string, input UpdateOrderInput) (*OrderOutput, apperrors.ApplicationError)
}

type updateOrderUsecase struct {
	repo orderRepo.Repository
}

// NewUpdateOrderUsecase creates a new instance of UpdateOrderUsecase
func NewUpdateOrderUsecase(repo orderRepo.Repository) UpdateOrderUsecase {
	return &updateOrderUsecase{repo: repo}
}

// Execute updates an order
func (u *updateOrderUsecase) Execute(ctx context.Context, id string, input UpdateOrderInput) (*OrderOutput, apperrors.ApplicationError) {
	// Validate order ID
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderInvalidIDError, err)
	}

	// Get existing order
	order, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if input.Status != nil {
		if !domain.IsValidStatus(*input.Status) {
			return nil, apperrors.NewApplicationError(mappings.OrderInvalidStatusError, nil)
		}
		order.Status = domain.OrderStatus(*input.Status)
	}

	if input.ETA != nil {
		order.ETA = *input.ETA
	}

	// Save changes
	updatedOrder, err := u.repo.Update(ctx, order)
	if err != nil {
		return nil, err
	}

	allStatuses := make([]string, len(domain.ValidStatuses))
	for i, s := range domain.ValidStatuses {
		allStatuses[i] = string(s)
	}

	return &OrderOutput{
		ID:          updatedOrder.ID,
		ProfileID:   updatedOrder.ProfileID,
		UserID:      updatedOrder.UserID,
		Status:      string(updatedOrder.Status),
		StatusIndex: updatedOrder.StatusIndex(),
		ETA:         updatedOrder.ETA,
		CreatedAt:   updatedOrder.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   updatedOrder.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		AllStatuses: allStatuses,
	}, nil
}
