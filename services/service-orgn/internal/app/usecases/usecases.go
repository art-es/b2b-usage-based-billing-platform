package usecases

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/create"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/get"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/get_by_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/update"
)

var (
	NewGetUsecase     = get.NewUsecase
	NewGetByIDUsecase = get_by_id.NewUsecase
	NewCreateUsecase  = create.NewUsecase
	NewUpdateUsecase  = update.NewUsecase
)
