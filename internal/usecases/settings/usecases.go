package settings

import (
	"context"
	"math"

	settingsRepo "wappi/internal/adapters/datasources/repositories/settings"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
)

// Usecases contains all settings-related use cases
type Usecases struct {
	Get                   GetUsecase
	Update                UpdateUsecase
	CalculateDeliveryFee  CalculateDeliveryFeeUsecase
}

// NewUsecases creates all settings usecases
func NewUsecases(repo settingsRepo.Repository) *Usecases {
	return &Usecases{
		Get:                   NewGetUsecase(repo),
		Update:                NewUpdateUsecase(repo),
		CalculateDeliveryFee:  NewCalculateDeliveryFeeUsecase(repo),
	}
}

// --- Get Usecase ---

type GetOutput struct {
	Settings *domain.Settings `json:"settings"`
}

type GetUsecase interface {
	Execute(ctx context.Context) (*GetOutput, apperrors.ApplicationError)
}

type getUsecase struct {
	repo settingsRepo.Repository
}

func NewGetUsecase(repo settingsRepo.Repository) GetUsecase {
	return &getUsecase{repo: repo}
}

func (u *getUsecase) Execute(ctx context.Context) (*GetOutput, apperrors.ApplicationError) {
	settings, err := u.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &GetOutput{Settings: settings}, nil
}

// --- Update Usecase ---

type UpdateInput struct {
	BusinessName        *string  `json:"business_name,omitempty"`
	BusinessLatitude    *float64 `json:"business_latitude,omitempty"`
	BusinessLongitude   *float64 `json:"business_longitude,omitempty"`
	DefaultMapLatitude  *float64 `json:"default_map_latitude,omitempty"`
	DefaultMapLongitude *float64 `json:"default_map_longitude,omitempty"`
	DefaultMapZoom      *int     `json:"default_map_zoom,omitempty"`
	DefaultItemWeight   *int     `json:"default_item_weight,omitempty"`
	DeliveryBasePrice   *float64 `json:"delivery_base_price,omitempty"`
	DeliveryPricePerKm  *float64 `json:"delivery_price_per_km,omitempty"`
	DeliveryPricePerKg  *float64 `json:"delivery_price_per_kg,omitempty"`
}

type UpdateOutput struct {
	Settings *domain.Settings `json:"settings"`
}

type UpdateUsecase interface {
	Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
}

type updateUsecase struct {
	repo settingsRepo.Repository
}

func NewUpdateUsecase(repo settingsRepo.Repository) UpdateUsecase {
	return &updateUsecase{repo: repo}
}

func (u *updateUsecase) Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	// Get current settings
	current, err := u.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.BusinessName != nil {
		current.BusinessName = *input.BusinessName
	}
	if input.BusinessLatitude != nil {
		current.BusinessLatitude = *input.BusinessLatitude
	}
	if input.BusinessLongitude != nil {
		current.BusinessLongitude = *input.BusinessLongitude
	}
	if input.DefaultMapLatitude != nil {
		current.DefaultMapLatitude = *input.DefaultMapLatitude
	}
	if input.DefaultMapLongitude != nil {
		current.DefaultMapLongitude = *input.DefaultMapLongitude
	}
	if input.DefaultMapZoom != nil {
		current.DefaultMapZoom = *input.DefaultMapZoom
	}
	if input.DefaultItemWeight != nil {
		current.DefaultItemWeight = *input.DefaultItemWeight
	}
	if input.DeliveryBasePrice != nil {
		current.DeliveryBasePrice = *input.DeliveryBasePrice
	}
	if input.DeliveryPricePerKm != nil {
		current.DeliveryPricePerKm = *input.DeliveryPricePerKm
	}
	if input.DeliveryPricePerKg != nil {
		current.DeliveryPricePerKg = *input.DeliveryPricePerKg
	}

	// Save
	updated, err := u.repo.Upsert(ctx, current)
	if err != nil {
		return nil, err
	}

	return &UpdateOutput{Settings: updated}, nil
}

// --- Calculate Delivery Fee Usecase ---

type CalculateDeliveryFeeInput struct {
	UserLatitude  float64 `json:"user_latitude" binding:"required"`
	UserLongitude float64 `json:"user_longitude" binding:"required"`
	Items         []struct {
		Quantity int  `json:"quantity"`
		Weight   *int `json:"weight,omitempty"` // in grams, optional
	} `json:"items"`
}

type CalculateDeliveryFeeOutput struct {
	DistanceKm   float64 `json:"distance_km"`
	TotalWeightG int     `json:"total_weight_g"`
	TotalWeightKg float64 `json:"total_weight_kg"`
	BasePrice    float64 `json:"base_price"`
	DistancePrice float64 `json:"distance_price"`
	WeightPrice  float64 `json:"weight_price"`
	TotalPrice   float64 `json:"total_price"`
}

type CalculateDeliveryFeeUsecase interface {
	Execute(ctx context.Context, input CalculateDeliveryFeeInput) (*CalculateDeliveryFeeOutput, apperrors.ApplicationError)
}

type calculateDeliveryFeeUsecase struct {
	repo settingsRepo.Repository
}

func NewCalculateDeliveryFeeUsecase(repo settingsRepo.Repository) CalculateDeliveryFeeUsecase {
	return &calculateDeliveryFeeUsecase{repo: repo}
}

func (u *calculateDeliveryFeeUsecase) Execute(ctx context.Context, input CalculateDeliveryFeeInput) (*CalculateDeliveryFeeOutput, apperrors.ApplicationError) {
	settings, err := u.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate distance using Haversine formula
	distanceKm := haversineDistance(
		settings.BusinessLatitude, settings.BusinessLongitude,
		input.UserLatitude, input.UserLongitude,
	)

	// Calculate total weight
	totalWeightG := 0
	for _, item := range input.Items {
		weight := settings.DefaultItemWeight
		if item.Weight != nil {
			weight = *item.Weight
		}
		totalWeightG += weight * item.Quantity
	}
	totalWeightKg := float64(totalWeightG) / 1000.0

	// Calculate prices
	basePrice := settings.DeliveryBasePrice
	distancePrice := distanceKm * settings.DeliveryPricePerKm
	weightPrice := totalWeightKg * settings.DeliveryPricePerKg
	totalPrice := basePrice + distancePrice + weightPrice

	return &CalculateDeliveryFeeOutput{
		DistanceKm:    math.Round(distanceKm*100) / 100,
		TotalWeightG:  totalWeightG,
		TotalWeightKg: math.Round(totalWeightKg*100) / 100,
		BasePrice:     basePrice,
		DistancePrice: math.Round(distancePrice*100) / 100,
		WeightPrice:   math.Round(weightPrice*100) / 100,
		TotalPrice:    math.Round(totalPrice*100) / 100,
	}, nil
}

// haversineDistance calculates the distance between two points on Earth in km
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	lat1Rad := degreesToRadians(lat1)
	lat2Rad := degreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
