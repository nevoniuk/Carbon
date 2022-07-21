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
		region_not_found
		no_data
		server_error
		if gerr.Name == "region_not_found" {
			return genpoller.MakeRegionNotFound()
		}
	}
	var dbNotFound dynamo.ErrNotFound
	if errors.As(err, &dbNotFound) {
		return gencoordinator.MakeNotFound(dbNotFound)
	}
	var fNotFound *facilityconfig.ErrNotFound
	if errors.As(err, &fNotFound) {
		return gencoordinator.MakeNotFound(fNotFound)
	}
	var syncerNotFound syncer.ErrNotFound
	if errors.As(err, &syncerNotFound) {
		return gencoordinator.MakeNotFound(syncerNotFound)
	}
	log.Error(ctx, err)
	return err
}

// mapAndLogErrorf maps client errors to a coordinator service error responses using
// the format string.
func mapAndLogErrorf(ctx context.Context, format string, a ...interface{}) error {
	return mapAndLogError(ctx, fmt.Errorf(format, a...))
}