//+build wireinject

package server

import (
	"futuagro.com/pkg/config"
	"futuagro.com/pkg/domain/services"
	"futuagro.com/pkg/http/rest"
	"futuagro.com/pkg/store"
	"github.com/google/wire"
)

func provideServer(supplierRepository *mongodb.SupplierRepository, ) *rest.Server {
	wire.Build(
		config.NewDefaultConfig,
		services.NewSupplierService,
	)
	return nil // These return values are ignored.
}
