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
	FindControlPointIDsByName(orgID uuid.UUID, clientName, pointName string) ([]uuid.UUID, error)
}

func (pc *powerConfig) getOrgByID(orgID uuid.UUID) (*facilityconfig.Org, error) {
	for _, org := range pc.f.Orgs {
		if org.ID == orgID {
			return org, nil
		}
	}
	return nil, fmt.Errorf("organisation %v does not exist in the config", orgID)
}

func (pc *powerConfig) getAgentByName(org *facilityconfig.Org, name string) (*facilityconfig.Agent, error) {
	for _, agent := range org.AgentByID {
		if strings.EqualFold(agent.Name, name) {
			return agent, nil
		}
	}
	return nil, fmt.Errorf("agent %s cannot be found in config", name)
}

func (pc *powerConfig) FindControlPointIDsByName(orgID uuid.UUID, clientName string, pointName string) ([]uuid.UUID, error) {
	org, err := pc.getOrgByID(orgID)
	var nullid []uuid.UUID
	if err != nil {
		return nullid, err
	}

	agent, err := pc.getAgentByName(org, clientName)
	
	if err != nil {
		return nullid, err
	}

	var points []uuid.UUID
	cp, ok := agent.ControlPointByAliasName[pointName]
	var err1 = fmt.Errorf("error given pointName")
	if !ok {
		return nullid, err1
	}
	points = append(points, cp.ID)
	return points, nil
}
//function to obtain pointname: configuration file to load in
//schedlure
