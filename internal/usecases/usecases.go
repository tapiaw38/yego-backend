package usecases

import (
	"yego/internal/adapters/web/websocket"
	"yego/internal/platform/appcontext"
	"yego/internal/usecases/admin"
	"yego/internal/usecases/order"
	"yego/internal/usecases/profile"
	"yego/internal/usecases/settings"
)

type Usecases struct {
	Order    Order
	Profile  Profile
	Admin    Admin
	Settings Settings
}

type Order struct {
	CreateUsecase               order.CreateUsecase
	CreateWithLinkUsecase       order.CreateWithLinkUsecase
	ClaimUsecase                order.ClaimUsecase
	GetUsecase                  order.GetUsecase
	GetClaimInfoUsecase         order.GetClaimInfoUsecase
	PayForOrderUsecase          order.PayForOrderUsecase
	CreatePaymentLinkUsecase    order.CreatePaymentLinkUsecase
	HandlePaymentWebhookUsecase order.HandlePaymentWebhookUsecase
	UpdateStatusUsecase         order.UpdateStatusUsecase
	ListMyOrdersUsecase         order.ListMyOrdersUsecase
}

type Profile struct {
	GenerateLinkUsecase    profile.GenerateLinkUsecase
	ValidateTokenUsecase   profile.ValidateTokenUsecase
	CompleteProfileUsecase profile.CompleteProfileUsecase
	GetUsecase             profile.GetProfileUsecase
	UpdateUsecase          profile.UpdateProfileUsecase
	UpsertUsecase          profile.UpsertProfileUsecase
	CheckCompletedUsecase  profile.CheckCompletedUsecase
}

type Admin struct {
	ListProfilesUsecase     admin.ListProfilesUsecase
	ListOrdersUsecase       admin.ListOrdersUsecase
	ListTransactionsUsecase admin.ListTransactionsUsecase
	UpdateOrderUsecase      admin.UpdateOrderUsecase
}

type Settings struct {
	GetUsecase                  settings.GetUsecase
	UpdateUsecase               settings.UpdateUsecase
	CalculateDeliveryFeeUsecase settings.CalculateDeliveryFeeUsecase
}

func CreateUsecases(contextFactory appcontext.Factory) *Usecases {
	app := contextFactory()
	hub := app.Integrations.WebSocket.GetHub()
	notifier := websocket.NewNotifier(hub)

	settingsUsecases := Settings{
		GetUsecase:                  settings.NewGetUsecase(contextFactory),
		UpdateUsecase:               settings.NewUpdateUsecase(contextFactory),
		CalculateDeliveryFeeUsecase: settings.NewCalculateDeliveryFeeUsecase(contextFactory),
	}

	return &Usecases{
		Order: Order{
			CreateUsecase:               order.NewCreateUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase),
			CreateWithLinkUsecase:       order.NewCreateWithLinkUsecase(contextFactory),
			ClaimUsecase:                order.NewClaimUsecase(contextFactory, notifier, settingsUsecases.CalculateDeliveryFeeUsecase),
			GetUsecase:                  order.NewGetUsecase(contextFactory),
			GetClaimInfoUsecase:         order.NewGetClaimInfoUsecase(contextFactory),
			PayForOrderUsecase:          order.NewPayForOrderUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase),
			CreatePaymentLinkUsecase:    order.NewCreatePaymentLinkUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase),
			HandlePaymentWebhookUsecase: order.NewHandlePaymentWebhookUsecase(contextFactory),
			UpdateStatusUsecase:         order.NewUpdateStatusUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase),
			ListMyOrdersUsecase:         order.NewListMyOrdersUsecase(contextFactory),
		},
		Profile: Profile{
			GenerateLinkUsecase:    profile.NewGenerateLinkUsecase(contextFactory),
			ValidateTokenUsecase:   profile.NewValidateTokenUsecase(contextFactory),
			CompleteProfileUsecase: profile.NewCompleteProfileUsecase(contextFactory),
			GetUsecase:             profile.NewGetProfileUsecase(contextFactory),
			UpdateUsecase:          profile.NewUpdateProfileUsecase(contextFactory),
			UpsertUsecase:          profile.NewUpsertProfileUsecase(contextFactory),
			CheckCompletedUsecase:  profile.NewCheckCompletedUsecase(contextFactory),
		},
		Admin: Admin{
			ListProfilesUsecase:     admin.NewListProfilesUsecase(contextFactory),
			ListOrdersUsecase:       admin.NewListOrdersUsecase(contextFactory),
			ListTransactionsUsecase: admin.NewListTransactionsUsecase(contextFactory),
			UpdateOrderUsecase:      admin.NewUpdateOrderUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase),
		},
		Settings: settingsUsecases,
	}
}
