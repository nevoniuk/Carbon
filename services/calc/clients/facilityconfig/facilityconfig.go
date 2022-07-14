package facilityconfig
import (
	"fmt"
	"github.com/crossnokaye/facilityconfig"
	"github.com/google/uuid"
	"io/ioutil"
	"path/filepath"
	"errors"
	"regexp"
	"context"
	"goa.design/clue/log"
	"gopkg.in/yaml.v3"
)

//steps to make this work
//1. define interface and methods necessar to get control point name from FILE
//2. define shape and write the FILE for any configurations
//3. probably add this file to another YAML file
//4. load that file using loader
//5. Add region identifier to facility.yaml


//changes to facility.yaml file
//carbon:
	/*
	id:
	....
	carbon:
		controlpointname:
		scale: ->this may change depending on how to make power formula generic for all facilities
		region:
	# ...
	timezone: 'America/Los_Angeles' # required for threshold algo with tou rates
	carbon:
		controlPointAliasName: whatever that is for the pulse val
		scale: .6(in Oxnard's example)
		formula:(may or may not need this depending on how generic Riverside data ends up being)
	 */
type(

	Client interface {
		LoadFacilityConfig(ctx context.Context, env, orgID, facilityID string) (*facilityConfig, error)
	}
	ControlPoint struct {
		ID        uuid.UUID
		Name      string
		Units     string
		Scaling   string
		Unscaling string
	}
	client struct {
		env string
	}
	//need to make a configstruct for the above data
	facilityConfig struct {
		ID string `yaml: id`
		Carbon *CarbonConfig `yaml: carbon`
	}

	CarbonConfig struct {
		ControlPointName string `yaml:"controlPointAliasName"`
		scale float64 `yaml: "scale"`
		region string `yaml: "region"`
	}

	// ErrNotFound is returned when a facility config is not found.
	ErrNotFound struct{ Err error }

	// ErrNoConfig is returned when no schedule type is specified.
	ErrNoConfig struct{ Err error }
)

var (
	FacilityDataFilePath = "deploy/facility_data"
)

// New returns a new client for the facility config data.
func New(env string) Client {
	return &client{env: env}
}

//GetCarbonConfig obtains the above carbon configuration for the given input
//called in load facility config to get carbonconfig for facility config struct
func GetCarbonConfig(ctx context.Context, orgID string, agentName string, facilityID string) (*CarbonConfig, error) {
	return nil, nil
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
	return "", &ErrNotFound{fmt.Errorf("org %s not found", orgID)}
}

// findFacility returns the path to the facility config for the given org and facility IDs.
func findFacility(ctx context.Context, env, orgID, facilityID string) (string, error) {
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
		return "", &ErrNotFound{fmt.Errorf("facility config not found for org path %s facility %s", path, facilityID)}
	}
	return facilityPath, nil
}

// findLocation returns the path to the building config for the given org, facility, and agent IDs.
func findLocation(ctx context.Context, env, orgID, facilityID, agentID string) (string, error) {
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
		agentPath := filepath.Join(filepath.Dir(path), b.Name(), "agent.yaml")
		read := readID(ctx, agentPath)
		if env != "production" {
			read = mapIDToNonProd(read, facilityID)
			if env == "office" {
				name := readName(ctx, agentPath)
				id := mapAgentToNonProd(name)
				if id != "" {
					read = id
				}
			}
		}
		if read == agentID {
			locationPath = filepath.Join(filepath.Dir(path), b.Name(), "location.yaml")
			break
		}
	}
	if locationPath == "" {
		return "", &ErrNotFound{errors.New("location config not found")}
	}
	return locationPath, nil
}


// loadFacilityConfig returns the facility config for the given org and facility IDs.
func (c *client) LoadFacilityConfig(ctx context.Context, env, orgID, facilityID string) (*facilityConfig, error) {
	facilityPath, err := findFacility(ctx, env, orgID, facilityID)
	if err != nil {
		return nil, err
	}
	cfg, err := ioutil.ReadFile(facilityPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read facility config file %s: %w", facilityPath, err)
	}
	//call to carbonconfig
	var config facilityConfig
	if err := yaml.Unmarshal(cfg, &config); err != nil {
		return nil, fmt.Errorf("failed to parse facility config file %s: %w", facilityPath, err)
	}
	return &config, nil
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

// validate validates the facility config.
func validate(fc *facilityConfig) error {
	if fc.ID== "" {
		return &ErrNoConfig{errors.New("ID not specified")}
	}
	if fc.Carbon.region == "" {
		return fmt.Errorf("region not specified")
	}
	if fc.Carbon.ControlPointName== "" {
		return fmt.Errorf("control point name not specified")
	}
	if fc.Carbon.scale == 0 {
		return fmt.Errorf("formula scale name not specified")
	}
	return nil
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
func (err ErrNotFound) Error() string { return err.Err.Error() }
func (err ErrNoConfig) Error() string { return err.Err.Error() }