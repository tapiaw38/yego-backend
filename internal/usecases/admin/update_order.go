package admin

import (
	"context"

	orderRepo "wappi/internal/adapters/datasources/repositories/order"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"

	"github.com/google/uuid"
)

// UpdateOrderInputItem represents an item in the order data
type UpdateOrderInputItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// UpdateOrderInputData represents the order data structure
type UpdateOrderInputData struct {
	Items []UpdateOrderInputItem `json:"items"`
}

// UpdateOrderInput represents the input for updating an order
type UpdateOrderInput struct {
	Status        *string               `json:"status,omitempty"`
	StatusMessage *string               `json:"status_message,omitempty"`
	ETA           *string               `json:"eta,omitempty"`
	Data          *UpdateOrderInputData `json:"data,omitempty"`
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

	if input.StatusMessage != nil {
		order.StatusMessage = input.StatusMessage
	}

	if input.ETA != nil {
		order.ETA = *input.ETA
	}

	if input.Data != nil {
		items := make([]domain.OrderItem, len(input.Data.Items))
		for i, item := range input.Data.Items {
			items[i] = domain.OrderItem{
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			}
		}
		order.Data = &domain.OrderData{Items: items}
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
		ID:            updatedOrder.ID,
		ProfileID:     updatedOrder.ProfileID,
		UserID:        updatedOrder.UserID,
		Status:        string(updatedOrder.Status),
		StatusMessage:  updatedOrder.StatusMessage,
		StatusIndex:    updatedOrder.StatusIndex(),
		ETA:            updatedOrder.ETA,
		CreatedAt:      updatedOrder.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      updatedOrder.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		AllStatuses:    allStatuses,
	}, nil
}
