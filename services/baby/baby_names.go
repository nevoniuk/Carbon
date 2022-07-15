package babynamesapi

import (
	"context"
	"log"

	babynames "github.com/crossnokaye/carbon/services/baby/gen/baby_names"
)

// BabyNames service example implementation.
// The example methods log the requests and return zero values.
type babyNamessrvc struct {
	logger *log.Logger
}

// NewBabyNames returns the BabyNames service implementation.
func NewBabyNames(ctx context.Context) babynames.Service {
	return &babyNamessrvc{ctx}
}

// get most popular baby name
func (s *babyNamessrvc) GetName(ctx context.Context, year string) (res *babynames.Name, err error) {
	res = &babynames.Name{}
	s.logger.Print("babyNames.GetName")
	return
}
