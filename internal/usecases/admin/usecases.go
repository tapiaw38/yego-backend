package admin

import (
	orderRepo "wappi/internal/adapters/datasources/repositories/order"
	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
)

// Usecases aggregates all admin-related use cases
type Usecases struct {
	ListProfiles ListProfilesUsecase
	ListOrders   ListOrdersUsecase
	UpdateOrder  UpdateOrderUsecase
}

// NewUsecases creates all admin use cases
func NewUsecases(profileRepository profileRepo.Repository, orderRepository orderRepo.Repository) *Usecases {
	return &Usecases{
		ListProfiles: NewListProfilesUsecase(profileRepository),
		ListOrders:   NewListOrdersUsecase(orderRepository),
		UpdateOrder:  NewUpdateOrderUsecase(orderRepository),
	}
}
