package calcapi

import (
	"context"
	"errors"
	"fmt"
	"goa.design/clue/log"
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
	var carbonreportsNotFound storage.ErrReportsNotFound
	var fNotFound facilityconfig.ErrFacilityNotFound
	var lNotFound facilityconfig.ErrLocationNotFound
	var powerreportsNotFound power.ErrPowerReportsNotFound
	if errors.As(err, &carbonreportsNotFound) {
		log.Error(ctx, err)
		return gencalc.MakeReportsNotFound(carbonreportsNotFound)
	}
	if errors.As(err, &fNotFound) {
		log.Error(ctx, err)
		return gencalc.MakeFacilityNotFound(fNotFound)
	}
	if errors.As(err, &lNotFound) {
		log.Error(ctx, err)
		return gencalc.MakeFacilityNotFound(lNotFound)
	}
	if errors.As(err, &powerreportsNotFound) {
		log.Error(ctx, err)
		return gencalc.MakeReportsNotFound(powerreportsNotFound)
	}
	log.Error(ctx, err)
	return err
}