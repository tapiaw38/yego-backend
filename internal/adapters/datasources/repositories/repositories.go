package repositories

import (
	"yego/internal/adapters/datasources"
	"yego/internal/adapters/datasources/repositories/coupon"
	"yego/internal/adapters/datasources/repositories/importrecord"
	"yego/internal/adapters/datasources/repositories/order"
	"yego/internal/adapters/datasources/repositories/ordertoken"
	"yego/internal/adapters/datasources/repositories/profile"
	"yego/internal/adapters/datasources/repositories/settings"
	"yego/internal/adapters/datasources/repositories/transaction"
)

type Repositories struct {
	Coupon       coupon.Repository
	ImportRecord importrecord.Repository
	Order        order.Repository
	OrderToken   ordertoken.Repository
	Profile      profile.Repository
	Settings     settings.Repository
	Transaction  transaction.Repository
}

type Factory func() *Repositories

func NewFactory(datasources *datasources.Datasources) func() *Repositories {
	return func() *Repositories {
		return &Repositories{
			Coupon:       coupon.NewRepository(datasources.DB),
			ImportRecord: importrecord.NewRepository(datasources.DB),
			Order:        order.NewRepository(datasources.DB),
			OrderToken:   ordertoken.NewRepository(datasources.DB),
			Profile:      profile.NewRepository(datasources.DB),
			Settings:     settings.NewRepository(datasources.DB),
			Transaction:  transaction.NewRepository(datasources.DB),
		}
	}
}
