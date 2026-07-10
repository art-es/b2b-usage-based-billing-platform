package trxutil

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
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
