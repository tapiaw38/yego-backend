package order

import (
	"context"
	"errors"
	"time"

	orderRepo "wappi/internal/adapters/datasources/repositories/order"
	ordertokenRepo "wappi/internal/adapters/datasources/repositories/ordertoken"
	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// ClaimInput represents the input for claiming an order
type ClaimInput struct {
	Token  string `json:"token" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

// ClaimOutput represents the output after claiming an order
type ClaimOutput struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	ProfileID *string `json:"profile_id,omitempty"`
	Status    string  `json:"status"`
	ETA       string  `json:"eta"`
	ClaimedAt string  `json:"claimed_at"`
}

// ClaimUsecase defines the interface for claiming orders
type ClaimUsecase interface {
	Execute(ctx context.Context, input ClaimInput) (*ClaimOutput, apperrors.ApplicationError)
}

type claimUsecase struct {
	orderRepo      orderRepo.Repository
	orderTokenRepo ordertokenRepo.Repository
	profileRepo    profileRepo.Repository
}

// NewClaimUsecase creates a new instance of ClaimUsecase
func NewClaimUsecase(orderRepo orderRepo.Repository, orderTokenRepo ordertokenRepo.Repository, profileRepo profileRepo.Repository) ClaimUsecase {
	return &claimUsecase{
		orderRepo:      orderRepo,
		orderTokenRepo: orderTokenRepo,
		profileRepo:    profileRepo,
	}
}

// Execute claims an order for a user
func (u *claimUsecase) Execute(ctx context.Context, input ClaimInput) (*ClaimOutput, apperrors.ApplicationError) {
	// Get the order token
	orderToken, err := u.orderTokenRepo.GetByToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}

	// Check if token has expired
	if time.Now().After(orderToken.ExpiresAt) {
		return nil, apperrors.NewApplicationError(mappings.OrderTokenExpiredError, errors.New("token expired"))
	}

	// Check if already claimed
	if orderToken.ClaimedAt != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderTokenAlreadyClaimedError, errors.New("token already claimed"))
	}

	// Get the order to verify it's not already assigned
	order, err := u.orderRepo.GetByID(ctx, orderToken.OrderID)
	if err != nil {
		return nil, err
	}

	// Check if order is already assigned to another user
	if order.UserID != nil && *order.UserID != input.UserID {
		return nil, apperrors.NewApplicationError(mappings.OrderAlreadyAssignedError, errors.New("order already assigned to another user"))
	}

	// Assign user to order
	if assignErr := u.orderRepo.AssignUser(ctx, orderToken.OrderID, input.UserID); assignErr != nil {
		return nil, assignErr
	}

	// Try to get user's profile and assign it to the order
	profile, _ := u.profileRepo.GetByUserID(ctx, input.UserID)
	if profile != nil {
		// User has a profile, assign it to the order
		_ = u.orderRepo.AssignProfile(ctx, orderToken.OrderID, profile.ID)
	}

	// Mark token as claimed
	if markErr := u.orderTokenRepo.MarkAsClaimed(ctx, input.Token, input.UserID); markErr != nil {
		return nil, markErr
	}

	// Get the updated order
	updatedOrder, err := u.orderRepo.GetByID(ctx, orderToken.OrderID)
	if err != nil {
		return nil, err
	}

	return &ClaimOutput{
		OrderID:   updatedOrder.ID,
		UserID:    input.UserID,
		ProfileID: updatedOrder.ProfileID,
		Status:    string(updatedOrder.Status),
		ETA:       updatedOrder.ETA,
		ClaimedAt: time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}
