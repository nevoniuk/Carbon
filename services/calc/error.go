package calcapi

import (
	"context"
	"errors"
	"fmt"
	"goa.design/clue/log"
	goa "goa.design/goa/v3/pkg"
	"github.com/crossnokaye/carbon/services/calc/clients/storage"
	"github.com/crossnokaye/carbon/services/calc/clients/power"
	"github.com/crossnokaye/carbon/services/calc/clients/facilityconfig"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
)
// mapAndLogErrorf maps client errors to a coordinator service error responses using
// the format string.
func mapAndLogErrorf(ctx context.Context, format string, a ...interface{}) error {
	return mapAndLogError(ctx, fmt.Errorf(format, a...))
}

// mapAndLogError maps client errors to a coordinator service error responses.
// It logs the error using the context if it does not map to a design error
// (i.e. is unexpected).
func mapAndLogError(ctx context.Context, err error) error {
	var gerr *goa.ServiceError
	var carbonreportsNotFound storage.ErrReportsNotFound
	var fNotFound facilityconfig.ErrFacilityNotFound
	var lNotFound facilityconfig.ErrLocationNotFound
	var powerreportsNotFound power.ErrPowerReportsNotFound
	if errors.As(err, &carbonreportsNotFound) {
		gerr = gencalc.MakeReportsNotFound(carbonreportsNotFound)
	}
	if errors.As(err, &fNotFound) {
		gerr = gencalc.MakeFacilityNotFound(fNotFound)
	}
	if errors.As(err, &lNotFound) {
		gerr = gencalc.MakeFacilityNotFound(lNotFound)
	}
	if errors.As(err, &powerreportsNotFound) {
		gerr = gencalc.MakeReportsNotFound(powerreportsNotFound)
	} else {
		log.Errorf(ctx, err, "Error not found: %w", err)
		gerr = gencalc.MakeReportsNotFound(err)
	}
	log.Error(ctx, gerr)
	return err
}