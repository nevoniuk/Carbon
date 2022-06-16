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
		Error("data_not_available", ErrorResult, "The data is not available or server error")
		Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//response to errors
			Response("data_not_available", CodeDataLoss)
			Response("missing-required-parameter", CodeNotFound)
		})
	})

	Method("fuels", func() {
		Description("query api using a search call for a fuel event from Carbonara API")
		Error("data_not_available", ErrorResult, "The data is not available or server error")
		Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//response to errors
			Response("data_not_available", CodeDataLoss)
			Response("missing-required-parameter", CodeNotFound)
		})
	})

	Method("aggregate_data", func() {
		Description("get the aggregate data for an event from clickhouse")
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
	Field(9, "aggregate_data_generated", aggregateData, "aggregate_data_generated", func() {

	})
	Field(10, "aggregate_data_consumed", aggregateData, "aggregate_data_consumed", func() {

	})
	Field(11, "aggregate_data_marginal", aggregateData, "aggregate_data_marginal", func() {

	})
	Field(12, "report_duration", String, "report_duration", func() {
		Example("hour, day, week, month")
	})
	Field(13, "region", String, "region", func() {
		Example("MISO, ISO...")
	})
	Required("generated_rate", "marginal_rate", "consumed_rate",
		"duration", "marginal_source", "consumed_source", "generated_source", "emission_factor",
		 "aggregate_data_generated","aggregate_data_consumed", "aggregate_data_marginal", "report_duration", "region")
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
	Field(5, "report_duration", String, "report_duration", func() {
		Example("hour, day, week, month")
	})

	Required("fuels", "duration", "marginal_source", "generated_source", "report_duration")
})

var aggregateData = Type("aggregateData", func() {
	Description("aggregate data")

	Field(1, "min", Float64, "min", func() {
		Example(37.8267)
	})
	Field(2, "max", Float64, "max", func() {
		Example(37.8267)
	})
	Field(3, "sum", Float64, "sum", func() {
		Example(37.8267)
	})
	Field(4, "count", Float64, "count", func() {
		Example(37.8267)
	})
	Field(5, "period", Period, "period", func() {
		//Example(37.8267)
	})

	Required("count", "max", "min", "sum", "period")
})

var CarbonResponse = Type("CarbonResponse", func() {
	Field(1, "HourlyReports", HourlyCarbonReports, "HourlyReports")
	Field(2, "DailyReports", DailyCarbonReports, "DailyReports")
	Field(3, "WeeklyReports", WeeklyCarbonReports, "WeeklyReports")
	Field(4, "MonthlyReports", MonthlyCarbonReports, "MonthlyReports")
	Required("HourlyReports", "DailyReports", "WeeklyReports", "MonthlyReports")
})

var HourlyCarbonReports = Type("HourlyCarbonReports", func ()  {
	Description("Array of hourly carbon reports for an area")
	Field(1, "HourlyReports", ArrayOf(CarbonForecast), "HourlyReports")
	Required("HourlyReports")
})
var DailyCarbonReports = Type("HourlyCarbonReports", func ()  {
	Description("Array of daily carbon reports for an area")
	Field(1, "DailyReports", ArrayOf(CarbonForecast), "DailyReports")
	Required("DailyReports")
})
var WeeklyCarbonReports = Type("HourlyCarbonReports", func ()  {
	Description("Array of weekly carbon reports for an area")
	Field(1, "WeeklyReports", ArrayOf(CarbonForecast), "WeeklyReports")
	Required("DailyReports")
})
var MonthlyCarbonReports = Type("HourlyCarbonReports", func ()  {
	Description("Array of monthly carbon reports for an area")
	Field(1, "HourlyReports", ArrayOf(CarbonForecast), "MonhtlyReports")
	Required("HourlyReports")
})
var FuelMix = Type("FuelMix", func() {
	Description("Generated Fuel Mix")
	Field(1, "Fuels", ArrayOf(Fuel), "Fuels")
	//aggregate data per each fuel mix
	//for all fuels not each fuel
	Field(2, "aggregate_data", aggregateData, "aggregate_data")
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


//observation structure to hold data from API that clickhouse will unerstand
