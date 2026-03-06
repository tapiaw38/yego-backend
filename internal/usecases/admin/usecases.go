package admin

import (
	s3service "yego/internal/services/s3"
	"yego/internal/platform/appcontext"
	settingsUsecase "yego/internal/usecases/settings"
)

// Usecases aggregates all admin-related use cases
type Usecases struct {
	ListProfiles     ListProfilesUsecase
	ListOrders       ListOrdersUsecase
	ListTransactions ListTransactionsUsecase
	UpdateOrder      UpdateOrderUsecase
	UploadImport     UploadImportUsecase
	ListImports      ListImportsUsecase
	CreateImport     CreateImportUsecase
	UpdateImport     UpdateImportUsecase
	DeleteImport     DeleteImportUsecase
	ClearImports     ClearImportsUsecase
	PresignUpload    PresignUploadUsecase
}

// NewUsecases creates all admin use cases
func NewUsecases(contextFactory appcontext.Factory, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase, s3Client *s3service.Client) *Usecases {
	return &Usecases{
		ListProfiles:     NewListProfilesUsecase(contextFactory),
		ListOrders:       NewListOrdersUsecase(contextFactory),
		ListTransactions: NewListTransactionsUsecase(contextFactory),
		UpdateOrder:      NewUpdateOrderUsecase(contextFactory, calculateDeliveryFeeUse),
		UploadImport:     NewUploadImportUsecase(contextFactory),
		ListImports:      NewListImportsUsecase(contextFactory),
		CreateImport:     NewCreateImportUsecase(contextFactory),
		UpdateImport:     NewUpdateImportUsecase(contextFactory),
		DeleteImport:     NewDeleteImportUsecase(contextFactory),
		ClearImports:     NewClearImportsUsecase(contextFactory),
		PresignUpload:    NewPresignUploadUsecase(s3Client),
	}
}
