package pollerapi

import (
	"context"
	"errors"
	"fmt"
	"goa.design/clue/log"
	goa "goa.design/goa/v3/pkg"
	"github.com/crossnokaye/carbon/services/poller/clients/carbonara"
	"github.com/crossnokaye/carbon/services/poller/clients/storage"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

// mapAndLogError maps client errors to a coordinator service error responses.
// It logs the error using the context if it does not map to a design error
// (i.e. is unexpected).
func mapAndLogError(ctx context.Context, err error) error {
	var gerr *goa.ServiceError
	if errors.As(err, &gerr) {
		if gerr.Name == "server_error" {
			return genpoller.MakeServerError(gerr)
		}
		if gerr.Name == "no_data" {
			return genpoller.MakeServerError(gerr)
		}
	}
	var serverError carbonara.ServerError
	var noDataError carbonara.NoDataError
	if errors.As(err, &serverError) {
		return genpoller.MakeServerError(serverError)
	}
	if errors.As(err, &noDataError) {
		return genpoller.MakeNoData(noDataError)
	}
	var noReports storage.NoReportsError
	if errors.As(err, &noReports) {
		return genpoller.MakeNoData(noReports)
	}
	var badReports storage.IncorrectReportsError
	if errors.As(err, &badReports) {
		return genpoller.MakeNoData(badReports)
	}
	log.Error(ctx, err)
	return err
}

// mapAndLogErrorf maps client errors to a coordinator service error responses using
// the format string.
func mapAndLogErrorf(ctx context.Context, format string, a ...interface{}) error {
	return mapAndLogError(ctx, fmt.Errorf(format, a...))
}