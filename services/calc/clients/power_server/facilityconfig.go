package power_server

import(
	"fmt"
	"github.com/crossnokaye/facilityconfig"
	"github.com/google/uuid"
	"strings"
)


func New(f *facilityconfig.Store) *powerConfig {
	return &powerConfig{f: f}
}

type powerConfig struct {
	f *facilityconfig.Store
}
//control points
//

func (pc *powerConfig) getOrgByID(orgID uuid.UUID) (*facilityconfig.Org, error) {
	for _, org := range pc.f.Orgs {
		if org.ID == orgID {
			return org, nil
		}
	}
	return nil, fmt.Errorf("organisation %v does not exist in the config", orgID)
}

