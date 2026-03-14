package order

import (
	"yego/internal/platform/appcontext"
	"yego/internal/usecases/notification"
	settingsUsecase "yego/internal/usecases/settings"
)

// Usecases aggregates all order-related use cases
type Usecases struct {
	Create         CreateUsecase
	CreateWithLink CreateWithLinkUsecase
	Claim          ClaimUsecase
	Get            GetUsecase
	UpdateStatus   UpdateStatusUsecase
	ListMyOrders   ListMyOrdersUsecase
}

// NewUsecases creates all order use cases
func NewUsecases(contextFactory appcontext.Factory, notificationSvc notification.Service, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase) *Usecases {
	return &Usecases{
		Create:         NewCreateUsecase(contextFactory, calculateDeliveryFeeUse, notificationSvc),
		CreateWithLink: NewCreateWithLinkUsecase(contextFactory),
		Claim:          NewClaimUsecase(contextFactory, notificationSvc, calculateDeliveryFeeUse),
		Get:            NewGetUsecase(contextFactory),
		UpdateStatus:   NewUpdateStatusUsecase(contextFactory, calculateDeliveryFeeUse, notificationSvc),
		ListMyOrders:   NewListMyOrdersUsecase(contextFactory),
	}
}
