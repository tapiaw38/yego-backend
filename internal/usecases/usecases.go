package usecases

import (
	"yego/internal/adapters/web/websocket"
	"yego/internal/platform/appcontext"
	s3service "yego/internal/services/s3"
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
	UploadImport            admin.UploadImportUsecase
	ListImports             admin.ListImportsUsecase
	CreateImport            admin.CreateImportUsecase
	UpdateImport            admin.UpdateImportUsecase
	DeleteImport            admin.DeleteImportUsecase
	ClearImports            admin.ClearImportsUsecase
	PresignUpload           admin.PresignUploadUsecase
	DeleteUpload            admin.DeleteUploadUsecase
	ListCoupons             admin.ListCouponsUsecase
	CreateCoupon            admin.CreateCouponUsecase
	UpdateCoupon            admin.UpdateCouponUsecase
	DeleteCoupon            admin.DeleteCouponUsecase
}

type Settings struct {
	GetUsecase                  settings.GetUsecase
	UpdateUsecase               settings.UpdateUsecase
	CalculateDeliveryFeeUsecase settings.CalculateDeliveryFeeUsecase
}

func CreateUsecases(contextFactory appcontext.Factory, s3Client *s3service.Client) *Usecases {
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
			CreateUsecase:               order.NewCreateUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase, notifier),
			CreateWithLinkUsecase:       order.NewCreateWithLinkUsecase(contextFactory),
			ClaimUsecase:                order.NewClaimUsecase(contextFactory, notifier, settingsUsecases.CalculateDeliveryFeeUsecase),
			GetUsecase:                  order.NewGetUsecase(contextFactory),
			GetClaimInfoUsecase:         order.NewGetClaimInfoUsecase(contextFactory),
			PayForOrderUsecase:          order.NewPayForOrderUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase),
			CreatePaymentLinkUsecase:    order.NewCreatePaymentLinkUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase),
			HandlePaymentWebhookUsecase: order.NewHandlePaymentWebhookUsecase(contextFactory),
			UpdateStatusUsecase:         order.NewUpdateStatusUsecase(contextFactory, settingsUsecases.CalculateDeliveryFeeUsecase, notifier),
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
			UploadImport:            admin.NewUploadImportUsecase(contextFactory),
			ListImports:             admin.NewListImportsUsecase(contextFactory),
			CreateImport:            admin.NewCreateImportUsecase(contextFactory),
			UpdateImport:            admin.NewUpdateImportUsecase(contextFactory),
			DeleteImport:            admin.NewDeleteImportUsecase(contextFactory),
			ClearImports:            admin.NewClearImportsUsecase(contextFactory),
			PresignUpload:           admin.NewPresignUploadUsecase(s3Client),
			DeleteUpload:            admin.NewDeleteUploadUsecase(s3Client),
			ListCoupons:             admin.NewListCouponsUsecase(contextFactory),
			CreateCoupon:            admin.NewCreateCouponUsecase(contextFactory),
			UpdateCoupon:            admin.NewUpdateCouponUsecase(contextFactory),
			DeleteCoupon:            admin.NewDeleteCouponUsecase(contextFactory),
		},
		Settings: settingsUsecases,
	}
}
