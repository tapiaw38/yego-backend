package profile

import (
	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
)

// Usecases aggregates all profile-related use cases
type Usecases struct {
	GenerateLink       GenerateLinkUsecase
	ValidateToken      ValidateTokenUsecase
	CompleteProfile    CompleteProfileUsecase
	Get                GetProfileUsecase
	Update             UpdateProfileUsecase
	CreateOrUpdate     CreateOrUpdateProfileUsecase
	CheckCompleted     CheckCompletedUsecase
}

// NewUsecases creates all profile use cases
func NewUsecases(repo profileRepo.Repository) *Usecases {
	return &Usecases{
		GenerateLink:    NewGenerateLinkUsecase(repo),
		ValidateToken:   NewValidateTokenUsecase(repo),
		CompleteProfile: NewCompleteProfileUsecase(repo),
		Get:             NewGetProfileUsecase(repo),
		Update:          NewUpdateProfileUsecase(repo),
		CreateOrUpdate:  NewCreateOrUpdateProfileUsecase(repo),
		CheckCompleted:  NewCheckCompletedUsecase(repo),
	}
}
