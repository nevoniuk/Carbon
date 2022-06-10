package design

import (
	//"github.com/dimfeld/httppath"
	. "goa.design/goa/v3/dsl"
)

var _ = API("Data Service API", func() {
	Title("The Smart Service API")
	Description("The Smart Service will query the Singularity Carbonara API daily and store in Clickhouse")
})

//always use == instead of = for clickhouse queries

var _ = Service("Data", func() {
	Description("Service that provides forecasts to clickhouse from Carbonara API")

	Method("carbon_emissions", func() {
		Description("query api getting search data for carbon_intensity event")
		Payload(ArrayOf(String), Period)
		Result(CarbonForecast)
		Error("data_not_available", ErrorResult, "The data is not available or server error")
		Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//response to errors
			Response("data_not_available", CodeDataLoss)
			Response("missing-required-parameter", CodeNotFound)
		})
	})

	Method("fuels", func() {
		Description("query api using a search call for a fuel event")
		Payload(ArrayOf(String), Period, String)
		Result(FuelsForecast)
		Error("data_not_available", ErrorResult, "The data is not available or server error")
		Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//response to errors
			Response("data_not_available", CodeDataLoss)
			Response("missing-required-parameter", CodeNotFound)
		})
	})

	Method("get_aggregate_data", func() {
		Description("get the aggregate data for an event")
		Payload(String, Period, String)
		Result(aggregateData)
		Error("data_not_available", ErrorResult, "The data is not available or server error")
		Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//response to errors
			Response("data_not_available", CodeDataLoss)
			Response("missing-required-parameter", CodeNotFound)
		})
	})
	/**
	Method("list_emission_factors", func() {
		Description("lists all emission factors used to calculate carbon intensity")
		Payload(Region)
		Result(ArrayOf(Emissionfactor))
		Error("Forbidden", ErrorResult, "Forbidden")
		Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//response to errors
			Response("data_not_available", CodeDataLoss)
			Response("missing-required-parameter", CodeNotFound)
		})
	})

	Method("get_nearest_ba", func() {
		Description("retrieves the nearest ba within a region given a postal code")
		Payload(Region)
		Result()
		Error("Forbidden", ErrorResult, "Forbidden")
		Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//response to errors
			Response("Forbidden", CodeNotFound)
			Response("missing-required-parameter", CodeNotFound)
		})
	})
	*/
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

	Field(5, "marginal_source", String, "marginal_source", func() {
		Example("EGRID_2019")
	})
	Field(6, "consumed_source", String, "consumed_source", func() {
		Example("EGRID_2019")
	})
	Field(7, "generated_source", String, "generated_source", func() {
		Example("EGRID_2019")
	})
	Field(8, "emission_factor", String, "emission_factor", func() {
		Example("EGRID_2019")
	})

	Required("generated_rate", "marginal_rate", "consumed_rate",
		"duration", "marginal_source", "consumed_source", "generated_source", "emission_factor")
})

var FuelsForecast = Type("FuelsForecast", func() {
	Description("Emissions Forecast")

	Field(1, "fuels", FuelMix, "fuels")
	Field(2, "duration", Period, "duration")
	Field(3, "marginal_source", String, "marginal_source", func() {
		Example("EGRID_2019")
	})
	Field(4, "generated_source", String, "generated_source", func() {
		Example("EGRID_2019")
	})

	Required("fuels", "duration", "marginal_source", "generated_source")
})

var aggregateData = Type("aggregateData", func() {
	Description("aggregate data")

	Field(1, "average", Float64, "average", func() {
		Example(37.8267)
	})
	Field(2, "count", Int, "count", func() {
		Example(1000)
	})
	Field(3, "max", Float64, "max", func() {
		Example(37.8267)
	})
	Field(4, "min", Float64, "min", func() {
		Example(37.8267)
	})
	Field(5, "sum", Float64, "sum", func() {
		Example(37.8267)
	})

	Required("average", "count", "max", "min", "sum")
})

/**
var DailyForecast = Type("DailyForecast", func() {
	Description("Daily Emissions Forecast")

	Field(1, "generated_rate", Float64, "generated_rate", func() {
		Example(37.8267)
	})
	Field(2, "marginal_rate", Float64, "marginal_rate", func() {
		Example(37.8267)
	})
	Field(3, "consumed_rate", Float64, "consumed_rate", func() {
		Example(37.8267)
	})
	Field(4, "duration", Period, "duration", func() {
		//Example("2022-06-07T00:20:00+00:00")
	})
	Field(5, "marginal_source", String, "marginal_source", func() {
		Example("EGRID_2019")
	})
	Field(6, "consumed_source", String, "consumed_source", func() {
		Example("EGRID_2019")
	})
	Field(7, "generated_source", String, "generated_source", func() {
		Example("EGRID_2019")
	})

	Required("generated_rate", "marginal_rate", "consumed_rate",
	"duration", "marginal_source", "consumed_source", "generated_source")
})

var MonthlyForecast = Type("MonthlyForecast", func() {
	Description("MonthlyEmissions Forecast")

	Field(1, "generated_rate", Float64, "generated_rate", func() {
		Example(37.8267)
	})
	Field(2, "marginal_rate", Float64, "marginal_rate", func() {
		Example(37.8267)
	})
	Field(3, "consumed_rate", Float64, "consumed_rate", func() {
		Example(37.8267)
	})
	Field(4, "duration", Period, "duration", func() {
		//Example("2022-06-07T00:20:00+00:00")
	})
	Field(5, "marginal_source", String, "marginal_source", func() {
		Example("EGRID_2019")
	})
	Field(6, "consumed_source", String, "consumed_source", func() {
		Example("EGRID_2019")
	})
	Field(7, "generated_source", String, "generated_source", func() {
		Example("EGRID_2019")
	})

var HourlyForecast = Type("HourlyForecast", func() {
	Description("Hourly Emissions Forecast")
	Field(1, "generated_rate", Float64, "generated_rate", func() {
		Example(37.8267)
	})
	Field(2, "marginal_rate", Float64, "marginal_rate", func() {
		Example(37.8267)
	})
	Field(3, "start_time", Float64, "start_time", func() {
		Example()
	})
	Field(4, "FuelMix", FuelMix, "FuelMix", func() {
		Example("coal_mw, nuclear_mw...")
	})
	Required("generated_rate", "marginal_rate", "start_time", "FuelMix")
})
*/
var FuelMix = Type("FuelMix", func() {
	Description("Generated Fuel Mix")
	Field(1, "Fuels", ArrayOf(Fuel), "Fuels")
	//could add forecast horizon in mins/hours
	Required("Fuels")
})

var Fuel = Type("Fuel", func() {
	Description("Generated Fuel Mix")
	Field(1, "mw", Float64, "MW", func() {
		Example(101)
	})
	Required("mw")
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

/**
var Region = Type("Region", func() {
	Description("region")
	Field(1, "Region", String, "region", func() {
		Example("MISO")
	})

	Required("Region")
})

var Event = Type("Event", func() {
	Description("event_type")
	Field(1, "Event", String, "event_type", func() {
		Example("carbon_intensity")
	})
	Required("Event")
})
*/
//observation structure to hold data from API that clickhouse will unerstand
