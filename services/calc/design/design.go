package design

import . "goa.design/goa/v3/dsl"

var _ = API("Calc", func() {
	Title("calc")
	Server("calc", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

var _ = Service("calc", func() {
	Description("Service to interpret CO2 emissions through power and carbon intensity data")
	Method("calculate_reports", func() {
		
		Description("helper method to make kW/lbs of Co2 report")
		GRPC(func() {})
	})

	Method("get_control_points", func() {
		//facility config stuff
		Description("wrapper for the power-service repo. gets the control points for the get_power function")
		Payload(PastValuesPayload)
		Result(ArrayOf(String)) //control points
		GRPC(func() {})
	})

	Method("get_power", func() {
		//talks to power client
		Payload(GetPowerPayload)
		Result(ElectricalReport)
		Description("This endpoint will retrieve the power data using control points from the get_control_points function")
		GRPC(func() {})
	})

	Method("get_emissions", func() {
		//talks to storage client
		Result(CarbonReport)
		Description("This endpoint will retrieve the emissions data for a facility")
		GRPC(func() {})
	})

	
	Method("handle_requests", func() {
		Payload(RequestPayload)
		//calls the above functions to get the power/emission reports
		Description("This endpoint is used by a front end service to return energy usage information")
		GRPC(func() {})
	})

	Method("carbonreport", func() {
		//R&D client
		Description("Make reports available to external/R&D clients")
		GRPC(func() {})
	})

})

//payloads
var RequestPayload = Type("RequestPayload", func() {
	Description("Payload for the handle_requests function")

	Field(1, "org", String, "org", func() {
		Format(FormatUUID)
	})
	Field(2, "Period", Period, "Period")

	Field(3, "building", String, "building", func() {
		Format(FormatUUID)
	})

	Field(4, "interval", String, "interval", func() {
		Example("hours, days, weeks, months, years")
	})
	Required("org", "Period", "building", "interval")
})

var PastValuesPayload = Type("PastValuesPayload", func() {
	Description("Payload for the past values get-values function")

	Field(1, "org", String, "org", func() {
		Format(FormatUUID)
	})
	Field(2, "Period", Period, "Period")

	Field(3, "building", String, "building", func() {
		Format(FormatUUID)
	})

	Field(4, "client", String, "client", func() {
		Format(FormatUUID)
	})
})

var GetPowerPayload = Type("GetPowerPayload", func() {
	Description("Payload for the past values get-values function")

	Field(1, "org", String, "org", func() {
		Format(FormatUUID)
	})

	Field(2, "Period", Period, "Period")

	Field(3, "cps", ArrayOf(String), "cps", func() {
		//Format(FormatUUID)
	})

	Field(4, "interval", Int64, "samping interval")

	Required("org", "Period", "cps", "interval")
})

var TotalReport = Type("TotalReport", func() {
	Description("Carbon/Energy Generation Report")
	
	Field(1, "Interval", Period, "Interval", func() {
		
	})
	Field(2, "point", ArrayOf(DataPoint), "point", func() {
		
	})
	Field(3, "facility", String, "facility", func() {
		
	})
	
	Required("Intervals", "point", "facility")
})

//reports read from clickhouse
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


var DataPoint = Type("DataPoint", func() {
	Description("Contains a time stamp with its respective x-y coordinates")

	Field(1, "time", String, "time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})

	Field(2, "carbon_rate", Float64, "carbon_rate", func() {
		Example(37.8267)
		//pounds of CO2
	})

	Required("Time", "carbon_rate",)
})

//reports read from past-values
var ElectricalReport = Type("ElectricalReport", func() {
	Description("Energy Generation Report")

	Field(1, "postalcode", String, "postalcode", func() {
		Format(FormatUUID)
	})
	Field(2, "facility", String, "facility", func() {
		Format(FormatUUID)
	})
	Field(3, "building", String, "building", func() {
		Format(FormatUUID)
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


