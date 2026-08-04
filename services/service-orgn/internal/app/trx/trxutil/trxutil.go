package trxutil

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
)

func RollbackOrLog(ctx context.Context, logger log.Logger, additionalInfo string) {
	if err := trx.Rollback(ctx); err != nil {
		logger.Log(log.Error).
			Set("message", "trx rollback error").
			Set("error", err.Error()).
			Set("additional_info", additionalInfo).
			Write()
	}
}
