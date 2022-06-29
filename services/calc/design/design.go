package design

import . "goa.design/goa/v3/dsl"

var _ = API("Trends Service", func() {
	Title("Service to interpret CO2 emissions through power and carbon intensity data")
	Server("design", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

//need some kidn of auth?
var _ = Service("calc", func() {
	/**
	Method("download_power", func() {
		//clickhouse client
		Description("Query clickhouse for power meter data")
		GRPC(func() {})
	})
	Method("download_carbon", func() {
		//clickhouse client
		Description("Query clickhouse for carbon intensity data")
		GRPC(func() {})
	})
	*/
	Method("calc_reports", func() {
		//helper method to make kW/lbs of Co2 report
		GRPC(func() {})
	})
	Method("display_carbon", func() {
		//talks to storage and power clients
		//by default returns co2 emissions for the past few months
		Description("This endpoint is used by a front end service to make carbon reports available")
		GRPC(func() {})
	})
	Method("handle_requests", func() {
		//frontend client
		//handles client requests for power usage from
		//1.a specific time peripd
		//2.an auth
		//3.a specific time interval
		Description("This endpoint is used by a front end service to return energy usage information")
		GRPC(func() {})
	})
	Method("carbonreport", func() {
		//R&D client
		Description("Make reports available to external/R&D clients")
		GRPC(func() {})
	})

})

var CarbonReport = Type("CarbonReport", func() {
	Description("Carbon/Energy Generation Report")
	
	Field(1, "Intervals", Period, "Intervals", func() {
		
	})
	Field(2, "point", DataPoint, "point", func() {
		Example(37.8267)
	})
	
	Required("Intervals", "point")
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

var DataPoint = Type("DataPoint", func() {
	Description("Contains a time stamp with its respective x-y coordinates")

	Field(1, "Time", String, "Time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})

	Field(2, "carbon_rate", Float64, "carbon_rate", func() {
		Example(37.8267)
	})

	Field(3, "power_rate", Float64, "power_rate", func() {
		Example(37.8267)
	})

	Required("Time", "carbon_rate", "power_rate")
})

var ElectricalReport = Type("ElectricalReport", func() {
	Description("Energy Generation Report")

	Field(1, "postalcode", Int, "postalcode", func() {
		
	})
	Field(2, "facility", Int, "facility", func() {
		
	})
	Field(2, "stamp", ArrayOf(PowerStamp), "stamp", func() {
		
	})
	
	Required("postalcode", "facility", "stamp")
})

var PowerStamp = Type("PowerStamp", func() {
	Field(1, "period", Period, "period", func() {

	})
	Field(2, "genRate", Float64, "genRate", func() {
		Description("power stamp in KW")
	})
})

