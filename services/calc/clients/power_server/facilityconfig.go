package power_server

import (
	"fmt"
	"strings"
	"github.com/crossnokaye/facilityconfig"
	"github.com/google/uuid"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"time"
	"github.com/google/uuid"
	"goa.design/clue/log"
	"gopkg.in/yaml.v3"
)
//steps to make this work
//1. define interface and methods necessar to get control point name from FILE
//2. define shape and write the FILE for any configurations
//3. probably add this file to another YAML file
//4. load that file using loader


//changes to facility.yaml file
//carbon:
	/*
		id: a5746ffa-2073-455e-b811-322ad3c3c4b7
		name: Oxnard Lineage
		shortname: oxnard
	# ...
	timezone: 'America/Los_Angeles' # required for threshold algo with tou rates
	carbon:
		controlPointAliasName: whatever that is for the pulse val
		scale: .6(in Oxnard's example)
		formula:(may or may not need this depending on how generic Riverside data ends up being)
	 */
type ControlPoint struct {
	ID        uuid.UUID
	Name      string
	Units     string
	Scaling   string
	Unscaling string
}
var (
	FacilityDataFilePath = "deploy/facility_data"
)
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
	GetOrgByID(uuid.UUID) (*facilityconfig.Org, error)
	GetAgentByName(*facilityconfig.Org, string) (*facilityconfig.Agent, error)


}

//GetOrgByID returns the org for the given uuid
func (pc *powerConfig) GetOrgByID(orgID uuid.UUID) (*facilityconfig.Org, error) {
	for _, org := range pc.f.Orgs {
		if org.ID == orgID {
			return org, nil
		}
	}
	return nil, fmt.Errorf("organisation %v does not exist in the config", orgID)
}

//GetAgentByName returns the agent for the given agent name
func (pc *powerConfig) GetAgentByName(org *facilityconfig.Org, name string) (*facilityconfig.Agent, error) {
	for _, agent := range org.AgentByID {
		if strings.EqualFold(agent.Name, name) {
			return agent, nil
		}
	}
	return nil, fmt.Errorf("agent %s cannot be found in config", name)
}

//FindControlPointIDsByName returns the control point ID for the given pointName(alias name)
func (pc *powerConfig) FindControlPointIDsByName(orgID uuid.UUID, clientName string, pointName string) (uuid.UUID, error) {
	org, err := pc.GetOrgByID(orgID)
	var nullid uuid.UUID
	if err != nil {
		return nullid, err
	}

	agent, err := pc.GetAgentByName(org, clientName)
	
	if err != nil {
		return nullid, err
	}

	
	cp, ok := agent.ControlPointByAliasName[pointName]
	var err1 = fmt.Errorf("error given pointName")
	if !ok {
		return nullid, err1
	}
	return cp.ID, nil
}

// findFacility returns the path to the facility config for the given org and facility IDs.
func (pc *powerConfig) findFacility(ctx context.Context, env, orgID, facilityID string) (string, error) {

}



