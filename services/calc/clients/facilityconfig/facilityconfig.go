package facilityconfig

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"github.com/google/uuid"
	"goa.design/clue/log"
	"gopkg.in/yaml.v3"
)

/*
	The following are changes to location.yaml file in order for the calc service to work
	carbon:
		controlPointAliasName: whatever that is for the pulse val
		formula:(This will be read and parsed using the GoMath)
		SingularityRegion:
*/

type(
	Client interface {
		GetCarbonConfig(ctx context.Context, orgID string, facilityID string, locationID string) (*Carbon, error)
	}
	client struct {
		env string
	}

	locationConfig struct {
		ID string `yaml: id`
		Carbon *carbonConfig `yaml: carbon`
	}

	carbonConfig struct {
		ControlPointName string `yaml: "controlPointAliasName"`
		Formula string `yaml: "formula"`
		SingularityRegion string `yaml: "region"`
	}

	Carbon struct {
		OrgID string
		FacilityID string
		BuildingID string
		ControlPointName string 
		Formula string
		Region string 
		AgentID string
	}

	// ErrNotFound is returned when a facility config is not found.
	ErrFacilityNotFound struct{ Err error }

	// ErrNotFound is returned when a location config is not found.
	ErrLocationNotFound struct{ Err error }

)
var (
	FacilityDataFilePath = "deploy/facility_data"
)
// New returns a new client for the facility config data.
func New(env string) Client {
	return &client{env: env}
}

//GetCarbonConfig obtains the above carbon configuration for the given input
func (c *client) GetCarbonConfig(ctx context.Context, orgID string, facilityID string, locationID string) (*Carbon, error) {
	path, config, err := loadLocationConfig(ctx, c.env, orgID, facilityID, locationID)
	if err != nil {
		return nil, err
	}
	id, err := findAgentIDFromLocation(ctx, c.env, orgID, facilityID, path)
	if err != nil {
		return nil, err
	}
	carbon := &Carbon{OrgID: orgID, FacilityID: facilityID, BuildingID: locationID, ControlPointName: config.Carbon.ControlPointName, Formula: config.Carbon.Formula, 
	Region: config.Carbon.SingularityRegion, AgentID: id}
	return carbon, nil
}

// findOrg finds the org file for the given org ID.
func findOrg(ctx context.Context, env, orgID string) (string, error) {
	files, err := ioutil.ReadDir(FacilityDataFilePath)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		path := filepath.Join(FacilityDataFilePath, f.Name(), "org.yaml")
		if readID(ctx, path) == orgID {
			return filepath.Dir(path), nil
		}
	}
	return "", ErrFacilityNotFound{Err: fmt.Errorf("org %s not found", orgID)}
}

// findFacility returns the path to the facility config for the given org and facility IDs.
func findFacility(ctx context.Context, env, orgID string, facilityID string) (string, error) {
	path, err := findOrg(ctx, env, orgID)
	if err != nil {
		return "", err
	}
	facilities, err := ioutil.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to list facilities in path %s: %w", path, err)
	}
	var facilityPath string
	for _, f := range facilities {
		if !f.IsDir() {
			continue
		}

		read := readID(ctx, filepath.Join(path, f.Name(), "facility.yaml"))
		if env != "production" {
			read = mapIDToNonProd(read, read)
		}
		if read == facilityID {
			facilityPath = filepath.Join(path, f.Name(), "facility.yaml")
			break
		}
	}
	if facilityPath == "" {
		return "", &ErrFacilityNotFound{fmt.Errorf("facility config not found for org path %s facility %s", path, facilityID)}
	}
	return facilityPath, nil
}
//findLocation will find the location path from location/building ID instead of the agentID
func findLocation(ctx context.Context, env string, orgID string, facilityID string, locationID string) (string, error) {
	path, err := findFacility(ctx, env, orgID, facilityID)
	if err != nil {
		return "", err
	}
	buildings, err := ioutil.ReadDir(filepath.Dir(path))
	if err != nil {
		return "", fmt.Errorf("failed to list buildings in path %s: %w", path, err)
	}
	var locationPath string
	for _, b := range buildings {
		if !b.IsDir() {
			continue
		}
		tempPath := filepath.Join(filepath.Dir(path), b.Name(), "location.yaml")
		read := readID(ctx, locationPath)
		if env != "production" {
			read = mapIDToNonProd(read, facilityID)
			if env == "office" {
				name := readName(ctx, tempPath)
				id := mapAgentToNonProd(name)
				if id != "" {
					read = id
				}
			}
		}
		if read == locationID {
			locationPath = tempPath
			break
		}
	}
	if locationPath == "" {
		return "", &ErrLocationNotFound{errors.New("location config not found")}
	}
	return locationPath, nil

}

//findAgentIDFromLocation is a separate function to take in the location path and return the agent ID for a location
func findAgentIDFromLocation(ctx context.Context, env, orgID, facilityID, locationPath string) (string, error) {
	if locationPath == "" {
		return "", fmt.Errorf("No location path")
	}
	var agentPath = filepath.Join(filepath.Dir(locationPath), "agent.yaml")
	read := readID(ctx, agentPath)
	return read, nil
}


// loadLocationConfig returns the building config for the given org, facility, and agent IDs.
// it will also return the buildingpath in order to avoid an extra function call to use findAgentIDFromLocation
func loadLocationConfig(ctx context.Context, env, orgID, facilityID, locationID string) (string, *locationConfig, error) {
	buildingPath, err := findLocation(ctx, env, orgID, facilityID, locationID)
	if err != nil {
		return "", nil, err
	}
	cfg, err := ioutil.ReadFile(buildingPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read building config file %s: %w", buildingPath, err)
	}
	var config locationConfig
	if err := yaml.Unmarshal(cfg, &config); err != nil {
		return "", nil, fmt.Errorf("failed to parse building config file %s: %w", buildingPath, err)
	}
	return buildingPath, &config, nil
}


func mapIDToNonProd(id, facilityID string) string {
	return mapToNonProd(uuid.MustParse(id), uuid.MustParse(facilityID))
}

func mapToNonProd(u, fid uuid.UUID) string {
	if u == uuid.Nil {
		return u.String()
	}
	key, _ := fid.MarshalBinary()
	return uuid.NewSHA1(u, key).String()
}

// Keep this in sync with crossnokaye/pkg/facilityconfig/loader/overrides.go
func mapAgentToNonProd(name string) string {
	switch name {
	case "Lineage Oxnard Building 4":
		return "5b9f1afa-d921-41bd-a8ad-0efddfa918ba"
	case "Lineage Riverside Building 3":
		return "4c76a509-3367-4a0e-9511-55e53f4cdcdd"
	case "Nordic Foods Kansas City Building 1":
		return "6abf56c8-adec-4434-be4d-3f83c98b5eab"
	case "Lineage Oxnard Building 3": // was old Office 1
		return "6329d21c-71f5-4ea1-830a-083d219fcb80"
	}
	return ""
}

func readID(ctx context.Context, path string) string {
	return readField(ctx, path, "id")
}

func readName(ctx context.Context, path string) string {
	return readField(ctx, path, "name")
}

func readField(ctx context.Context, path, field string) string {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		log.Infof(ctx, "error reading file %s: %v, ignoring", path, err)
		return ""
	}
	var data map[string]interface{}
	if err := yaml.Unmarshal(cfg, &data); err != nil {
		log.Infof(ctx, "error parsing file %s: %v, ignoring", path, err)
		return ""
	}
	if val, ok := data[field]; ok {
		if s, ok := val.(string); ok {
			return s
		}
		err := fmt.Errorf("invalid %s type %T", field, val)
		log.Infof(ctx, "error parsing file %s: %v, ignoring", path, err)
		return ""
	}
	err = fmt.Errorf("%s not found in file %s", field, path)
	log.Infof(ctx, "error parsing file %s: %v, ignoring", path, err)
	return ""
}


// timeRegex is a regular expression for parsing 24 hour time strings.
var timeRegex = regexp.MustCompile(`^(2[0-3]|[01]?[0-9]):([0-5]?[0-9])$`)
func (err ErrFacilityNotFound) Error() string { return err.Err.Error() }
func (err ErrLocationNotFound) Error() string { return err.Err.Error() }