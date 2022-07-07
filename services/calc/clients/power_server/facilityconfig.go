package power_server

import (
	"fmt"
	//"log"
	"strings"

	"github.com/crossnokaye/facilityconfig"
	"github.com/google/uuid"
)

type ControlPoint struct {
	ID        uuid.UUID
	Name      string
	Units     string
	Scaling   string
	Unscaling string
}

func New(f *facilityconfig.Store) *powerConfig {
	return &powerConfig{f: f}
}

type powerConfig struct {
	f *facilityconfig.Store
}
type Repository interface {
	// FindControlPointConfigsByName searches control point description in the agent with specified name
	// if point is not found, no error and empty slice is returned.
	FindControlPointIDByName(orgID uuid.UUID, clientName, pointName string) (uuid.UUID, error)

	// GetPointScalingByIDs retrieves a map of scaling formulas for points with provided IDs
	// if scaling for provided IDs does not exist, empty map is returned
	GetPointScalingByIDs(ids []uuid.UUID) map[uuid.UUID]string
}

//goal is to get 
//1. org ID
var deviceName = "Power meter data acquisition server"


func (pc *powerConfig) getOrgByID(orgID uuid.UUID) (*facilityconfig.Org, error) {
	for _, org := range pc.f.Orgs {
		if org.ID == orgID {
			return org, nil
		}
	}
	return nil, fmt.Errorf("organisation %v does not exist in the config", orgID)
}


//helper
func (pc *powerConfig) getAgentByName(org *facilityconfig.Org, name string) (*facilityconfig.Agent, error) {
	for _, agent := range org.AgentByID {
		if strings.EqualFold(agent.Name, name) {
			return agent, nil
		}
	}
	return nil, fmt.Errorf("agent %s cannot be found in config", name)
}


func (pc *powerConfig) FindControlPointIDByName(orgID uuid.UUID, clientName string, pointName string) (uuid.UUID, error) {
	org, err := pc.getOrgByID(orgID)
	var ret uuid.UUID
	if err != nil {
		return ret, err
	}

	agent, err := pc.getAgentByName(org, clientName)
	
	if err != nil {
		return ret, err
	}

	//var points []*ControlPoint
	cp, ok := agent.ControlPointByAliasName[pointName]
	var err1 = fmt.Errorf("error given pointName")
	if !ok {
		return ret, err1
	}
	return cp.ID, nil
}

