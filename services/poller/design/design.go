package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("Poller", func() {
	Title("Poller")
	Description("The Poller Service will query the Singularity Carbonara API daily and store in Clickhouse")
})


var _ = Service("Poller", func() {
	Description("Service that provides forecasts to clickhouse from Carbonara API")

	Method("carbon_emissions", func() {
		Description("query api getting search data for carbon_intensity event. Return reports in 5 minute intervals")
		Payload(CarbonPayload)
		Result(ArrayOf(CarbonForecast))
		//payload is region and start time
		//result is minute reports
		
		//Error("data_not_available", ErrorResult, "The data is not available or server error")
		//Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//Response("data_not_available", CodeDataLoss)
			//Response("missing-required-parameter", CodeNotFound)
		})
	})
	Method("aggregate_data", func() {
		Description("convert 5 minute reports into hourly, daily, monthly, yearly reports using clickhouse aggregate queries")
		Payload(AggregatePayload)
		//region and dates to get reports for

		//Error("data_not_available", ErrorResult, "The data is not available or server error")
		//Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//Response("data_not_available", CodeDataLoss)
			//Response("missing-required-parameter", CodeNotFound)
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
	Field(4, "Duration", Period, "Duration")

	Field(5, "duration_type", String, "duration_type")

	Field(5, "generated_source", String, "generated_source", func() {
		Example("EGRID_2019")
	})
	Field(6, "region", String, "region", func() {
		Example("MISO, ISO...")
	})
	Required("generated_rate", "marginal_rate", "consumed_rate", "generated_source", "region", "Duration", "duration_type")
})

/**
var AggregateData = Type("aggregateData", func() {

	Field(1, "average", Float64, "average", func() {
		Example(37.8267)
	})
	Field(2, "min", Float64, "min", func() {
		Example(37.8267)
	})
	Field(3, "max", Float64, "max", func() {
		Example(37.8267)
	})
	Field(4, "sum", Float64, "sum", func() {
		Example(37.8267)
	})
	Field(5, "count", Int, "count", func() {
		Example(50)
	})
	Field(6, "duration", Period, "duration", func() {
		//Example(37.8267)
	})
	Field(7, "report_type", String, "report_type", func() {
		Example("hourly")
	})

	Required("average", "count", "max", "min", "sum", "duration", "report_type")
})
*/

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
})

var AggregatePayload = Type("AggregatePayload", func() {
	Field(1, "region", String, "region", func() {
	})
	Field(2, "periods", ArrayOf(Period), "periods", func() {
	
	})
	Field(3, "duration", String, "duration", func() {
	
	})
})
