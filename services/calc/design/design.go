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
		Description("This endpoint is used by a front end service to return carbon emission reports")
		Payload(RequestPayload)
		Result(AllReports)
		GRPC(func() {})
	})
	Method("get_carbon_report", func() {
		Description("Make reports available to external/R&D clients")
		GRPC(func() {})
	})

})

var EmissionsPayload = Type("EmissionsPayload", func() {
	Description("Period is the range to get data for, interval is the time unit at which to parse the data by.")
	Field(1, "Period", Period, "Period")
	Field(2, "Interval", String, "Interval", func() {
		Example("hours, days, weeks, months, years")
	})
})

var AllReports = Type("AllReports", func() {
	Description("CO2 intensity reports, power reports, and CO2 emission reports")
	Field(1, "CarbonIntensityReports", ArrayOf(CarbonReport), "CarbonIntensityReports")
	Field(2, "PowerReports", ArrayOf(ElectricalReport), "PowerReports")
	Field(3, "TotalEmissionReports", ArrayOf(EmissionsReport), "TotalEmissionReports")
	Required("CarbonIntensityReports", "PowerReports", "TotalEmissionReports")
})

var RequestPayload = Type("RequestPayload", func() {
	Description("Payload wraps the payload for past-values GetValues() and carbon poller service")
	Field(1, "Org", UUID, "Org")
	Field(2, "Duration", Period, "Duration")
	Field(3, "Agent", String, "Agent")
	Field(4, "Interval", String, "Interval", func() {
		Example("hours, days, weeks, months, years")
	})
	Required("Org", "Period", "Agent", "Interval")
})

var EmissionsReport = Type("EmissionsReport", func() {
	Description("Carbon/Energy Generation Report")
	Field(1, "Duration", Period, "Duration")
	Field(2, "DurationType", String, "DurationType")
	Field(3, "Points", ArrayOf(DataPoint), "Points")
	Field(4, "Org", UUID, "Org")
	Field(5, "Agent", String, "Agent")
	
	Required("Duration", "Points", "Org", "DurationType", "Agent")
})

var CarbonReport = Type("CarbonReport", func() {
	Description("Carbon Report from clickhouse")
	
	Field(1, "GeneratedRate", Float64, "GeneratedRate")
	Field(2, "Duration", Period, "Duration")
	Field(3, "DurationType", String, "DurationType")
	Field(4, "Region", String, "Region", func() {
		Example("MISO, ISO...")
	})
	Required("GeneratedRate", "Region", "Duration", "DurationType")
})

var DataPoint = Type("DataPoint", func() {
	Description("Contains a time stamp with its respective x&y coordinates")
	Field(1, "Time", String, "Time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Field(2, "CarbonFootprint", Float64, "CarbonFootprint", func() {
		Example(37.8267)
		Description("carbon footprint is the lbs of CO2 emissions")
	})

	Required("Time", "CarbonFootprint")
})

var ElectricalReport = Type("ElectricalReport", func() {
	Description("Energy Generation Report from the Past values function GetValues")
	Field(1, "Duration", Period, "Duration")
	Field(2, "Org", UUID, "Org")
	Field(3, "Agent", String, "Agent")
	Field(4, "Stamp", PowerStamp, "Stamp")
	Field(5, "IntervalType", String, "IntervalType")

	Required("Org", "Duration", "Agent", "Stamp", "IntervalType")
})

var PowerStamp = Type("PowerStamp", func() {
	Description("Used by Electrical Report to store power meter data from GetValues()")
	Field(1, "Time", String, "Time", func() {
		Format(FormatDateTime)
	})
	Field(2, "GeneratedRate", Float64, "GeneratedRate", func() {
		Description("power stamp in KW")
	})
	Required("GeneratedRate", "Time")
})

var Period = Type("Period", func() {
	Description("Period of time from start to end for any report type")
	Field(1, "StartTime", String, "Start time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Field(2, "EndTime", String, "End time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Required("StartTime", "EndTime")
})

var UUID = Type("UUID", String, func() {
	Description("Universally unique identifier")
	Format(FormatUUID)
})


