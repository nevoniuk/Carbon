package pollerapi

import (
	"context"
	"errors"
	"fmt"
	"goa.design/clue/log"
	"github.com/crossnokaye/carbon/services/poller/clients/carbonara"
	"github.com/crossnokaye/carbon/services/poller/clients/storage"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

// mapAndLogError maps client errors to a coordinator service error responses.
// It logs the error using the context if it does not map to a design error
// (i.e. is unexpected).
func mapAndLogError(ctx context.Context, err error) error {
	var serverError carbonara.ServerError
	var noDataError carbonara.NoDataError
	var noReports storage.NoReportsError
	var badReports storage.IncorrectReportsError
	if errors.As(err, &serverError) {
		log.Error(ctx, err)
		return genpoller.MakeServerError(serverError)
	}
	if errors.As(err, &noDataError) {
		log.Error(ctx, err)
		return genpoller.MakeNoData(noDataError)
	}
	if errors.As(err, &noReports) {
		log.Error(ctx, err)
		return genpoller.MakeClickhouseError(noReports)
	}
	if errors.As(err, &badReports) {
		log.Error(ctx, err)
		return genpoller.MakeClickhouseError(badReports)
	}
	log.Error(ctx, err)
	return err
}

// mapAndLogErrorf maps client errors to a coordinator service error responses using
// the format string.
func mapAndLogErrorf(ctx context.Context, format string, a ...interface{}) error {
	return mapAndLogError(ctx, fmt.Errorf(format, a...))
}