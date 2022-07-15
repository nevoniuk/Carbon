// Code generated by goa v3.7.6, DO NOT EDIT.
//
// BabyNames HTTP client CLI support package
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/baby/design

package client

import (
	babynames "github.com/crossnokaye/carbon/services/baby/gen/baby_names"
)

// BuildGetNamePayload builds the payload for the BabyNames GetName endpoint
// from CLI flags.
func BuildGetNamePayload(babyNamesGetNameYear string) (*babynames.GetNamePayload, error) {
	var year string
	{
		year = babyNamesGetNameYear
	}
	v := &babynames.GetNamePayload{}
	v.Year = &year

	return v, nil
}
