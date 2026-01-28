package order

import (
	"context"

	"wappi/internal/adapters/datasources/repositories/order"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"

	"github.com/google/uuid"
)

// GetOutputItem represents a single item in the order output
type GetOutputItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// GetOutputData represents the order data in the output
type GetOutputData struct {
	Items []GetOutputItem `json:"items"`
}

// GetOutput represents the output for getting an order
type GetOutput struct {
	ID          string         `json:"id"`
	ProfileID   *string        `json:"profile_id,omitempty"`
	UserID      *string        `json:"user_id,omitempty"`
	Status      string         `json:"status"`
	StatusIndex int            `json:"status_index"`
	ETA         string         `json:"eta"`
	Data        *GetOutputData `json:"data,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	AllStatuses []string       `json:"all_statuses"`
}

// GetUsecase defines the interface for getting orders
type GetUsecase interface {
	Execute(ctx context.Context, id string) (*GetOutput, apperrors.ApplicationError)
}

type getUsecase struct {
	repo order.Repository
}

// NewGetUsecase creates a new instance of GetUsecase
func NewGetUsecase(repo order.Repository) GetUsecase {
	return &getUsecase{repo: repo}
}

// Execute retrieves an order by ID
func (u *getUsecase) Execute(ctx context.Context, id string) (*GetOutput, apperrors.ApplicationError) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderInvalidIDError, err)
	}

	orderData, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allStatuses := make([]string, len(domain.ValidStatuses))
	for i, s := range domain.ValidStatuses {
		allStatuses[i] = string(s)
	}

	// Convert domain data to output data
	var outputData *GetOutputData
	if orderData.Data != nil && len(orderData.Data.Items) > 0 {
		items := make([]GetOutputItem, len(orderData.Data.Items))
		for i, item := range orderData.Data.Items {
			items[i] = GetOutputItem{
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			}
		}
		outputData = &GetOutputData{Items: items}
	}

	return &GetOutput{
		ID:          orderData.ID,
		ProfileID:   orderData.ProfileID,
		UserID:      orderData.UserID,
		Status:      string(orderData.Status),
		StatusIndex: orderData.StatusIndex(),
		ETA:         orderData.ETA,
		Data:        outputData,
		CreatedAt:   orderData.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   orderData.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		AllStatuses: allStatuses,
	}, nil
}
