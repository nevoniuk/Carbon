package design

import . "goa.design/goa/v3/dsl"

var _ = API("Calc", func() {
	Title("Service to interpret CO2 emissions through power and carbon intensity data")
	Server("calc", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

//need some kidn of auth?
var _ = Service("calc", func() {
	
	Method("calculate_reports", func() {
		
		Description("helper method to make kW/lbs of Co2 report")
		GRPC(func() {})
	})

	Method("get_control_points", func() {
		//talks to power client
		Description("This endpoint will retrieve the control points for a facility")
		GRPC(func() {})
	})

	Method("get_power", func() {
		//talks to power client
		Description("This endpoint will retrieve the power data using control points from the past-values service")
		GRPC(func() {})
	})

	Method("get_emissions", func() {
		//talks to power client
		Description("This endpoint will retrieve the emissions data for a facility")
		GRPC(func() {})
	})

	
	Method("handle_requests", func() {
		//gets:
		//1.a specific time peripd
		//2.an auth
		//3.a specific time interval
		//4.a specific facility ID/name
		//5.a specific building
		//6. then calls the above functions to get the power/emission reports
		Description("This endpoint is used by a front end service to return energy usage information")
		GRPC(func() {})
	})

	Method("carbonreport", func() {
		//R&D client
		Description("Make reports available to external/R&D clients")
		GRPC(func() {})
	})

})
//1.make electrical reports
//2.make carbon reports
//3.make total reports from both
var TotalReport = Type("TotalReport", func() {
	Description("Carbon/Energy Generation Report")
	
	Field(1, "Interval", Period, "Interval", func() {
		
	})
	Field(2, "point", ArrayOf(DataPoint), "point", func() {
		
	})
	Field(3, "facility", Int, "facility", func() {
		
	})
	
	Required("Intervals", "point", "facility")
})

var CarbonReport = Type("CarbonReport", func() {
	Description("Carbon Report from clickhouse")
	
	Field(1, "generated_rate", Float64, "generated_rate", func() {
		Example(37.8267)
	})

	Field(2, "Duration", Period, "Duration")

	Field(3, "duration_type", String, "duration_type")

	Field(4, "region", String, "region", func() {
		Example("MISO, ISO...")
	})
	Required("generated_rate", "region", "Duration", "duration_type")
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
	Field(3, "building", Int, "building", func() {
		
	})
	Field(4, "stamp", ArrayOf(PowerStamp), "stamp", func() {
		
	})
	
	Required("postalcode", "facility", "stamp", "building")
})

var PowerStamp = Type("PowerStamp", func() {
	Field(1, "period", Period, "period", func() {

	})
	Field(2, "genRate", Float64, "genRate", func() {
		Description("power stamp in KW")
	})
})

