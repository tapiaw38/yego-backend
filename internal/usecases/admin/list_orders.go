package admin

import (
	"context"

	orderRepo "wappi/internal/adapters/datasources/repositories/order"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
)

// OrderOutput represents an order in the admin list
type OrderOutput struct {
	ID            string           `json:"id"`
	ProfileID     *string          `json:"profile_id,omitempty"`
	UserID        *string          `json:"user_id,omitempty"`
	Status        string           `json:"status"`
	StatusMessage *string          `json:"status_message,omitempty"`
	StatusIndex   int              `json:"status_index"`
	ETA           string           `json:"eta"`
	Data          *domain.OrderData `json:"data,omitempty"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
	AllStatuses   []string         `json:"all_statuses"`
}

// ListOrdersOutput represents the output for listing orders
type ListOrdersOutput struct {
	Orders []OrderOutput `json:"orders"`
	Total  int           `json:"total"`
}

// ListOrdersUsecase defines the interface for listing orders
type ListOrdersUsecase interface {
	Execute(ctx context.Context) (*ListOrdersOutput, apperrors.ApplicationError)
}

type listOrdersUsecase struct {
	repo orderRepo.Repository
}

// NewListOrdersUsecase creates a new instance of ListOrdersUsecase
func NewListOrdersUsecase(repo orderRepo.Repository) ListOrdersUsecase {
	return &listOrdersUsecase{repo: repo}
}

// Execute lists all orders
func (u *listOrdersUsecase) Execute(ctx context.Context) (*ListOrdersOutput, apperrors.ApplicationError) {
	orders, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	allStatuses := make([]string, len(domain.ValidStatuses))
	for i, s := range domain.ValidStatuses {
		allStatuses[i] = string(s)
	}

	output := &ListOrdersOutput{
		Orders: make([]OrderOutput, 0, len(orders)),
		Total:  len(orders),
	}

	for _, o := range orders {
		output.Orders = append(output.Orders, OrderOutput{
			ID:            o.ID,
			ProfileID:     o.ProfileID,
			UserID:        o.UserID,
			Status:        string(o.Status),
			StatusMessage: o.StatusMessage,
			StatusIndex:   o.StatusIndex(),
			ETA:           o.ETA,
			Data:          o.Data,
			CreatedAt:     o.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			AllStatuses:   allStatuses,
		})
	}

	return output, nil
}
