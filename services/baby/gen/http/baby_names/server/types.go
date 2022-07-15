// Code generated by goa v3.7.6, DO NOT EDIT.
//
// BabyNames HTTP server types
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/baby/design

package server

import (
	babynames "github.com/crossnokaye/carbon/services/baby/gen/baby_names"
)

// GetNameResponseBody is the type of the "BabyNames" service "GetName"
// endpoint HTTP response body.
type GetNameResponseBody struct {
	// name
	Name string `form:"name" json:"name" xml:"name"`
}

// NewGetNameResponseBody builds the HTTP response body from the result of the
// "GetName" endpoint of the "BabyNames" service.
func NewGetNameResponseBody(res *babynames.Name) *GetNameResponseBody {
	body := &GetNameResponseBody{
		Name: res.Name,
	}
	return body
}

// NewGetNamePayload builds a BabyNames service GetName endpoint payload.
func NewGetNamePayload(year string) *babynames.GetNamePayload {
	v := &babynames.GetNamePayload{}
	v.Year = &year

	return v
}
