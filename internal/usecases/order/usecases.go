package order

import (
	orderRepo "wappi/internal/adapters/datasources/repositories/order"
	ordertokenRepo "wappi/internal/adapters/datasources/repositories/ordertoken"
	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
)

// Usecases aggregates all order-related use cases
type Usecases struct {
	Create         CreateUsecase
	CreateWithLink CreateWithLinkUsecase
	Claim          ClaimUsecase
	Get            GetUsecase
	UpdateStatus   UpdateStatusUsecase
}

// NewUsecases creates all order use cases
func NewUsecases(repo orderRepo.Repository, tokenRepo ordertokenRepo.Repository, profileRepo profileRepo.Repository) *Usecases {
	return &Usecases{
		Create:         NewCreateUsecase(repo),
		CreateWithLink: NewCreateWithLinkUsecase(repo, tokenRepo),
		Claim:          NewClaimUsecase(repo, tokenRepo, profileRepo),
		Get:            NewGetUsecase(repo),
		UpdateStatus:   NewUpdateStatusUsecase(repo),
	}
}
