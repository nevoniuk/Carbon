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

	Method("handle_requests", func() {
		Description("This endpoint is used by a front end service to return energy usage information")
		Payload(RequestPayload)
		Result(EmissionsReport)
		GRPC(func() {})
	})

	Method("carbon_report", func() {
		Description("Make reports available to external/R&D clients")
		GRPC(func() {})
	})

})

//payloads
var EmissionsPayload = Type("EmissionsPayload", func() {
	Description("Payload for the get_emissions function")
	Field(1, "Period", Period, "Period")
	Field(2, "interval", String, "interval", func() {
		Example("hours, days, weeks, months, years")
	})
})

var AllReports = Type("AllReports", func() {
	Description("CO2 intensity reports, power reports, and CO2 emission reports")
	Field(1, "carbon_intensity_reports", ArrayOf(CarbonReport), "carbon_intensity_reports")
	Field(2, "power_reports", ArrayOf(ElectricalReport), "power_reports")
	Field(3, "total_emission_reports", ArrayOf(EmissionsReport), "total_emission_reports")
	Required("carbon_intensity_reports", "power_reports", "total_emission_reports")
})

var RequestPayload = Type("RequestPayload", func() {
	Description("Payload for the handle_requests function")

	Field(1, "org", UUID, "org")

	Field(2, "Period", Period, "Period")

	Field(3, "building", UUID, "building")

	Field(4, "interval", String, "interval", func() {
		Example("hours, days, weeks, months, years")
	})
	Required("org", "Period", "building", "interval")
})

var EmissionsReport = Type("EmissionsReport", func() {
	Description("Carbon/Energy Generation Report")
	
	Field(1, "Duration", Period, "Duration")

	Field(2, "duration_type", String, "duration_type")

	Field(3, "point", ArrayOf(DataPoint), "point", func() {
		
	})
	Field(4, "facility", UUID, "facility")
	
	Required("Duration", "point", "facility", "duration_type")
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

	Required("time", "carbon_rate")
})

//reports read from past-values
var ElectricalReport = Type("ElectricalReport", func() {
	Description("Energy Generation Report")

	Field(1, "period", Period, "period", func() {

	})
	Field(2, "postalcode", String, "postalcode", func() {
		
	})
	Field(3, "facility", UUID, "facility")

	Field(4, "building", UUID, "building")

	Field(5, "stamp", ArrayOf(PowerStamp), "stamp", func() {
		
	})

	Field(6, "intervalType", String, "intervalType")
	
	Required("postalcode", "facility", "stamp", "building", "intervalType")
})

var PowerStamp = Type("PowerStamp", func() {

	Field(1, "time", String, "time", func() {
		Format(FormatDateTime)
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

var UUID = Type("UUID", String, func() {
	Description("Universally unique identifier")
	Format(FormatUUID)
})


