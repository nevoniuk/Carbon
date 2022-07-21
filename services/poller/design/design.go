package design

import (
	. "goa.design/goa/v3/dsl"
	"github.com/crossnokaye/carbon/types/design"
)

var _ = API("Poller", func() {
	Title("Poller")
	Description("The Poller Service will query the Singularity Carbonara API daily and store in Clickhouse")
	Docs(func() {
		Description("Additional documentation on Singularity and Region Enums")
		URL("https://docs.google.com/document/d/1t-_9GNZLyI98LujRzXwbjMVE6mVNBeiV7O2pwIecd9I/edit#")
	})
})


var _ = Service("Poller", func() {
	Description("Service that provides forecasts to clickhouse from Carbonara API")
	Method("update", func() {
		Description("query Singularity's search endpoint and convert 5 min interval reports into averages")
		Error("server_error", ErrorResult, "Error with Singularity Server.")
		GRPC(func() {
			Response("server_error", CodeNotFound)
		})
	})
	Method("get_emissions_for_region", func() {
		Description("query search endpoint for a region.")
		Payload(CarbonPayload)
		Result(ArrayOf(CarbonForecast))
		Error("server_error", ErrorResult, "Error with Singularity Server.")
		Error("no_data", ErrorResult, "No new data available for any region")
		Error("region_not_found", ErrorResult, "The given region is not represented by Singularity")
		GRPC(func() {
			Response("no_data", CodeNotFound)
			Response("region_not_found", CodeNotFound)
			Response("server_error", CodeNotFound)
		})
	})
	
})

var CarbonForecast = Type("CarbonForecast", func() {
	Description("Emissions Forecast")
	
	Field(1, "generated_rate", Float64, "generated_rate", func() {
		Example(37.8267)
	})
	Field(2, "marginal_rate", Float64, "marginal_rate", func() {
		Example(37.8267)
	})
	Field(3, "consumed_rate", Float64, "consumed_rate", func() {
		Example(37.8267)
	})
	Field(4, "duration", Period, "duration")

	Field(4, "interval", design.IntervalType, "interval")

	Field(7, "region", String, "region", func() {
		Example("MISO, ISO...")
	})
	Required("generated_rate", "marginal_rate", "consumed_rate", "region", "duration", "interval")
})


var Period = Type("Period", func() {
	Description("Period of time from start to end of Forecast")
	Field(1, "startTime", String, "Start time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Field(2, "endTime", String, "End time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Required("startTime", "endTime")
})

var CarbonPayload = Type("CarbonPayload", func() {
	Field(1, "region", String, "region", func() {
	})
	Field(2, "start", String, "start", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Field(3, "end", String, "end", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
})
