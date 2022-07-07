package calcapi
//something
import (
	"context"
	"fmt"
	//"fmt"
	//"sync"
	"time"

	"github.com/crossnokaye/carbon/services/calc/clients/power"
	"github.com/crossnokaye/carbon/services/calc/clients/power_server"
	"github.com/crossnokaye/carbon/services/calc/clients/storage"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
)

type calcSvc struct {
	psc power.Client
	dbc storage.Client
	psr power_server.Repository
	ctx context.Context
	cancel context.CancelFunc
	
}
var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"

func NewCalc(ctx context.Context, psc power.Client, dbc storage.Client, psr power_server.Repository) *calcSvc {
	ctx, cancel := context.WithCancel(ctx)
	s := &calcSvc{
		psc:				psc,
		dbc:				dbc,
		psr:                psr,
		ctx:                ctx,
		cancel: 			cancel,
	}
	
	return s
}