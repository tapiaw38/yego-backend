package order

import (
	"context"
	"errors"
	"time"

	"wappi/internal/platform/appcontext"
	"wappi/internal/usecases/notification"
	settingsUsecase "wappi/internal/usecases/settings"
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
	contextFactory          appcontext.Factory
	notificationSvc         notification.Service
	calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase
}

// NewClaimUsecase creates a new instance of ClaimUsecase
func NewClaimUsecase(contextFactory appcontext.Factory, notificationSvc notification.Service, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase) ClaimUsecase {
	return &claimUsecase{
		contextFactory:          contextFactory,
		notificationSvc:         notificationSvc,
		calculateDeliveryFeeUse: calculateDeliveryFeeUse,
	}
}

// Execute claims an order for a user
func (u *claimUsecase) Execute(ctx context.Context, input ClaimInput) (*ClaimOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Get the order token
	orderToken, err := app.Repositories.OrderToken.GetByToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}

	// Check if token has expired
	if time.Now().After(orderToken.ExpiresAt) {
		return nil, apperrors.NewApplicationError(mappings.OrderTokenExpiredError, errors.New("token expired"))
	}

	// Check if already claimed
	if orderToken.ClaimedAt != nil {
		// Return the order info so frontend can redirect to it
		order, _ := app.Repositories.Order.GetByID(ctx, orderToken.OrderID)
		if order != nil {
			return &ClaimOutput{
				OrderID:   order.ID,
				UserID:    *order.UserID,
				ProfileID: order.ProfileID,
				Status:    string(order.Status),
				ETA:       order.ETA,
				ClaimedAt: orderToken.ClaimedAt.Format("2006-01-02T15:04:05Z"),
			}, nil
		}
		return nil, apperrors.NewApplicationError(mappings.OrderTokenAlreadyClaimedError, errors.New("token already claimed"))
	}

	// Get the order to verify it's not already assigned
	order, err := app.Repositories.Order.GetByID(ctx, orderToken.OrderID)
	if err != nil {
		return nil, err
	}

	// Check if order is already assigned to another user
	if order.UserID != nil && *order.UserID != input.UserID {
		return nil, apperrors.NewApplicationError(mappings.OrderAlreadyAssignedError, errors.New("order already assigned to another user"))
	}

	// Assign user to order
	if assignErr := app.Repositories.Order.AssignUser(ctx, orderToken.OrderID, input.UserID); assignErr != nil {
		return nil, assignErr
	}

	// Try to get user's profile and assign it to the order
	profile, _ := app.Repositories.Profile.GetByUserID(ctx, input.UserID)
	if profile != nil {
		// User has a profile, assign it to the order
		_ = app.Repositories.Order.AssignProfile(ctx, orderToken.OrderID, profile.ID)
	}

	// Mark token as claimed
	if markErr := app.Repositories.OrderToken.MarkAsClaimed(ctx, input.Token, input.UserID); markErr != nil {
		return nil, markErr
	}

	// Get the updated order
	updatedOrder, err := app.Repositories.Order.GetByID(ctx, orderToken.OrderID)
	if err != nil {
		return nil, err
	}

	// Send notification to managers
	if u.notificationSvc != nil {
		payload := notification.OrderClaimedPayload{
			OrderID:   updatedOrder.ID,
			UserID:    input.UserID,
			Status:    string(updatedOrder.Status),
			ETA:       updatedOrder.ETA,
			ClaimedAt: time.Now().Format("2006-01-02T15:04:05Z"),
		}
		if updatedOrder.ProfileID != nil {
			payload.ProfileID = *updatedOrder.ProfileID
		}
		// Send notification asynchronously to not block the response
		go func() {
			if notifyErr := u.notificationSvc.NotifyOrderClaimed(payload); notifyErr != nil {
				// Log error but don't fail the request
				// In production, you might want to use a proper logger
				_ = notifyErr
			}
		}()
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
