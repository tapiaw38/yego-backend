package order

import (
	"context"

	orderRepo "yego/internal/adapters/datasources/repositories/order"
	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
)

type ListDeliveryOrdersOutput struct {
	Orders []OrderOutputData `json:"orders"`
	Total  int               `json:"total"`
}

type ListDeliveryOrdersUsecase interface {
	Execute(ctx context.Context, deliveryUserID string) (*ListDeliveryOrdersOutput, apperrors.ApplicationError)
}

type listDeliveryOrdersUsecase struct {
	contextFactory appcontext.Factory
}

func NewListDeliveryOrdersUsecase(contextFactory appcontext.Factory) ListDeliveryOrdersUsecase {
	return &listDeliveryOrdersUsecase{contextFactory: contextFactory}
}

func (u *listDeliveryOrdersUsecase) Execute(ctx context.Context, deliveryUserID string) (*ListDeliveryOrdersOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	rows, err := app.Repositories.Order.GetByDeliveryUserID(ctx, deliveryUserID)
	if err != nil {
		return nil, err
	}

	output := &ListDeliveryOrdersOutput{
		Orders: make([]OrderOutputData, 0, len(rows)),
		Total:  len(rows),
	}

	for _, row := range rows {
		out := toOrderOutputData(&row.Order, true)
		out.ProfileInfo = toProfileDeliveryInfo(row)
		output.Orders = append(output.Orders, out)
	}

	return output, nil
}

func toProfileDeliveryInfo(row *orderRepo.DeliveryOrderRow) *ProfileDeliveryInfo {
	info := &ProfileDeliveryInfo{
		PhoneNumber: row.PhoneNumber,
	}
	if row.Address.Valid {
		info.Address = row.Address.String
	}
	if row.Latitude.Valid {
		lat := row.Latitude.Float64
		info.Latitude = &lat
	}
	if row.Longitude.Valid {
		lng := row.Longitude.Float64
		info.Longitude = &lng
	}
	return info
}
