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
	"strings"
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
		// GetCarbonConfig obtains the above carbon configuration for the given input
		GetCarbonConfig(ctx context.Context, orgID string, facilityID string, locationID string) (*Carbon, error)
	}
	client struct {
		env string
	}
	LocationConfig struct {
		ID string `yaml: id`
		Carbon *CarbonConfig `yaml: "carbon"`
	}
	CarbonConfig struct {
		ControlPoint string `yaml: "controlpoint"`
		Multiplier string `yaml: "multiplier"`
		SingularityRegion string `yaml: "singularityregion"`
	}

	Carbon struct {
		OrgID string
		FacilityID string
		BuildingID string
		ControlPointName string 
		Formula string
		Region string 
		AgentName string
	}

	// ErrNotFound is returned when a facility config is not found.
	ErrFacilityNotFound struct{ Err error }

	// ErrNotFound is returned when a location config is not found.
	ErrLocationNotFound struct{ Err error }

	// ErrConfigNotFound is returned when a Carbon config is not found or invalid
	ErrConfigNotFound struct { Err error }

)

var FacilityDataFilePath = "deploy/facility_data"

// New returns a new client for the facility config data.
func New(env string) Client {
	return &client{env: env}
}

// GetCarbonConfig will load the data from a location.yaml file into a Carbon struct
func (c *client) GetCarbonConfig(ctx context.Context, orgID string, facilityID string, locationID string) (*Carbon, error) {
	path, config, err := loadLocationConfig(ctx, c.env, orgID, facilityID, locationID)
	if err != nil {
		return nil, err
	}
	agentName, err := findAgentNameFromLocation(ctx, c.env, orgID, facilityID, path)
	if err != nil || agentName == "" {
		return nil, ErrConfigNotFound{fmt.Errorf("could not find the agent name for location %s with err: %w", path, err)}
	}
	if c.env == "office" {
		agentName = strings.Join([]string{"office ", agentName}, "")
	}
	if config == nil {
		return nil, ErrConfigNotFound{fmt.Errorf("could not find the carbon config for orgID: %s, agent %s, location %s, facility %s", orgID, agentName, locationID, facilityID)}
	}
	log.Info(ctx, log.KV{K: "carbon config", V: config.Carbon})
	carbon := &Carbon{OrgID: orgID, FacilityID: facilityID, BuildingID: locationID, ControlPointName: config.Carbon.ControlPoint, Formula: config.Carbon.Multiplier, 
	Region: config.Carbon.SingularityRegion, AgentName: agentName}
	err = validate(carbon)
	if err != nil {
		return nil, ErrConfigNotFound{fmt.Errorf("could not validate carbon config with err: %w\n", err)}
	}
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
	return "", err
}

// findFacility returns the path to the facility config for the given org and facility IDs.
func findFacility(ctx context.Context, env, orgID string, facilityID string) (string, error) {
	path, err := findOrg(ctx, env, orgID)
	if err != nil {
		return "", err
	}
	facilities, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
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
		return "", err
	}
	return facilityPath, nil
}
// findLocation will find the location path from location/building ID instead of the agentID
func findLocation(ctx context.Context, env string, orgID string, facilityID string, locationID string) (string, error) {
	path, err := findFacility(ctx, env, orgID, facilityID)
	if err != nil {
		return "", &ErrFacilityNotFound{fmt.Errorf("facility not found for org %s facility %s: %w", orgID, facilityID, err)}
	}
	buildings, err := ioutil.ReadDir(filepath.Dir(path))
	if err != nil {
		return "", &ErrLocationNotFound{fmt.Errorf("failed to list buildings in path %s: %w", path, err)}
	}

	var locationPath string
	for _, b := range buildings {
		if !b.IsDir() {
			continue
		}
		locationPath = filepath.Join(filepath.Dir(path), b.Name(), "location.yaml")
		read := readID(ctx, locationPath)
		if env != "production" {
			read = mapIDToNonProd(read, facilityID)
			if read != locationID {
				locationPath = ""
			} else {
				break
			}
		}
	}
	if locationPath == "" {
		return "", &ErrLocationNotFound{fmt.Errorf("location config not found for org %s and facility %s\n", orgID, facilityID)}
	}

	return locationPath, nil

}

// findAgentIDFromLocation is a separate function to take in the location path and return the agent ID for a location
func findAgentNameFromLocation(ctx context.Context, env, orgID, facilityID, locationPath string) (string, error) {
	if locationPath == "" {
		return "", fmt.Errorf("No location path")
	}
	var agentPath = filepath.Join(filepath.Dir(locationPath), "agent.yaml")
	read := readName(ctx, agentPath)
	return read, nil
}


// loadLocationConfig returns the building config for the given org, facility, and agent IDs.
// it will also return the buildingpath in order to avoid an extra function call to use findAgentIDFromLocation
func loadLocationConfig(ctx context.Context, env, orgID, facilityID, locationID string) (string, *LocationConfig, error) {
	buildingPath, err := findLocation(ctx, env, orgID, facilityID, locationID)
	if err != nil {
		return "", nil, err
	}
	cfg, err := ioutil.ReadFile(buildingPath)
	if err != nil {
		return "", nil, &ErrLocationNotFound{fmt.Errorf("failed to read building config file %s: %w", buildingPath, err)}
	}
	var config LocationConfig
	if err := yaml.Unmarshal(cfg, &config); err != nil {
		return "", nil, &ErrLocationNotFound{fmt.Errorf("failed to unmarshal into location config %s: %w", buildingPath, err)}
	}
	log.Info(ctx, log.KV{K: "location config", V: config.ID})
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

// readID returns the ID field for a given path
func readID(ctx context.Context, path string) string {
	return readField(ctx, path, "id")
}

// readName returns the Name field for given path
func readName(ctx context.Context, path string) string {
	return readField(ctx, path, "name")
}

// readField returns the given field for the given pth
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

// validate will ensure the Carbon struct is not null
func validate(fc *Carbon) error {
	if fc.AgentName == "" {
		return &ErrConfigNotFound{errors.New("agent name not specified")}
	}
	if fc.BuildingID == "" {
		return &ErrConfigNotFound{errors.New("no building id")}
	}
	if fc.ControlPointName == "" {
		return &ErrConfigNotFound{errors.New("control point name not specified")}
	}
	return nil
}


// timeRegex is a regular expression for parsing 24 hour time strings.
var timeRegex = regexp.MustCompile(`^(2[0-3]|[01]?[0-9]):([0-5]?[0-9])$`)
func (err ErrFacilityNotFound) Error() string { return err.Err.Error() }
func (err ErrLocationNotFound) Error() string { return err.Err.Error() }
func (err ErrConfigNotFound) Error() string { return err.Err.Error() }
