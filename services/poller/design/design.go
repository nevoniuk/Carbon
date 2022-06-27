package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("Poller", func() {
	Title("Poller")
	Description("The Poller Service will query the Singularity Carbonara API daily and store in Clickhouse")
})

//always use == instead of = for clickhouse queries
//Fuel design may not be needed so not finished
var _ = Service("Poller", func() {
	Description("Service that provides forecasts to clickhouse from Carbonara API")

	Method("carbon_emissions", func() {
		Description("query api getting search data for carbon_intensity event")
		Result(ArrayOf(ArrayOf(CarbonForecast)))
		//Error("data_not_available", ErrorResult, "The data is not available or server error")
		//Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//Response("data_not_available", CodeDataLoss)
			//Response("missing-required-parameter", CodeNotFound)
		})
	})
/**
	Method("fuels", func() {
		Description("query api using a search call for a fuel event from Carbonara API")
		Result(ArrayOf(FuelsForecast))
		//Error("data_not_available", ErrorResult, "The data is not available or server error")
		//Error("missing-required-parameter", ErrorResult, "missing-required-parameter")
		GRPC(func() {
			//Response("data_not_available", CodeDataLoss)
			//Response("missing-required-parameter", CodeNotFound)
		})
	})
*/
	Method("aggregate_data", func() {
		Description("get the aggregate data for an event from clickhouse")
		Result(ArrayOf(ArrayOf(AggregateData)))
		//Result(ArrayOf(ArrayOf(AggregateData)))
		//Error("data_not_available", ErrorResult, "The data is not available or server error")
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

	Field(5, "generated_source", String, "generated_source", func() {
		Example("EGRID_2019")
	})
	Field(6, "region", String, "region", func() {
		Example("MISO, ISO...")
	})
	Required("generated_rate", "marginal_rate", "consumed_rate", "generated_source", "region", "Duration")
})
/**
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
	Field(5, "report_type", String, "report_type", func() {
		Example("hour, day, week, month")
	})

	Required("fuels", "duration", "marginal_source", "generated_source", "report_type")
})
*/
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
/**
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
var DailyCarbonReports = Type("DailyCarbonReports", func ()  {
	Description("Array of daily carbon reports for an area")
	Field(1, "DailyReports", ArrayOf(CarbonForecast), "DailyReports")
	Required("DailyReports")
})
var WeeklyCarbonReports = Type("WeeklyCarbonReports", func ()  {
	Description("Array of weekly carbon reports for an area")
	Field(1, "WeeklyReports", ArrayOf(CarbonForecast), "WeeklyReports")
	Required("WeeklyReports")
})
var MonthlyCarbonReports = Type("MonthlyCarbonReports", func ()  {
	Description("Array of monthly carbon reports for an area")
	Field(1, "MonthlyReports", ArrayOf(CarbonForecast), "MonthlyReports")
	Required("MonthlyReports")
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
