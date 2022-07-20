package design

import (
	. "goa.design/goa/v3/dsl"
	"github.com/crossnokaye/carbon/model"
	//"github.com/crossnokaye/carbon/types/design"
)

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
	Error("reports_not_found", func() {
		Description("Carbon reports not found")
	})
	Error("not_found", func() {
		Description("facilty or location not found")
	})
	GRPC(func() {
		Response("not_found", CodeNotFound)
	})
	Method("historical_carbon_emissions", func() {
		Description("This endpoint is used by a front end service to return carbon emission reports")
		Payload(RequestPayload)
		Result(AllReports)
		GRPC(func() {
			Response("not_found", CodeNotFound)
		})
	})
})

var AllReports = Type("AllReports", func() {
	Description("CO2 intensity reports, power reports, and CO2 emission reports")
	Field(1, "CarbonIntensityReports", ArrayOf(CarbonReport), "CarbonIntensityReports")
	Field(2, "PowerReports", ArrayOf(ElectricalReport), "PowerReports")
	Field(3, "TotalEmissionReport", EmissionsReport, "TotalEmissionReport")
	Required("CarbonIntensityReports", "PowerReports", "TotalEmissionReport")
})

var RequestPayload = Type("RequestPayload", func() {
	Description("Payload wraps the payload for past-values GetValues() and carbon poller service")
	Field(1, "OrgID", UUID, "OrgID")
	Field(2, "Duration", Period, "Duration")
	Field(3, "FacilityID", UUID, "FacilityID")
	Field(4, "Interval", String, IntervalFunc)
	Field(5, "LocationID", UUID, "LocationID")
	Required("OrgID", "Duration", "Interval", "FacilityID", "LocationID")
})

var EmissionsReport = Type("EmissionsReport", func() {
	Description("Carbon/Energy Generation Report")
	Field(1, "Duration", Period, "Duration")
	Field(2, "Interval", String, IntervalFunc)
	Field(3, "Points", ArrayOf(DataPoint), "Points")
	Field(4, "OrgID", UUID, "OrgID")
	Field(5, "FacilityID", UUID, "FacilityID")
	Field(6, "LocationID", UUID, "LocationID")
	Field(7, "Region", String, RegionFunc)
	Required("Duration", "Points", "OrgID", "Interval", "FacilityID", "LocationID", "Region")
})

var CarbonReport = Type("CarbonReport", func() {
	Description("Carbon Report from clickhouse")
	Field(1, "GeneratedRate", Float64, "GeneratedRate", func() {
		Description("This is in units of (lbs of CO2/MWh)")
	})
	Field(2, "Duration", Period, "Duration")
	Field(3, "Interval", String, IntervalFunc)
	Field(4, "Region", String, RegionFunc)
	Required("GeneratedRate", "Region", "Duration", "Interval")
})

var DataPoint = Type("DataPoint", func() {
	Description("Contains carbon emissions in terms of DataPoints, which can be used as points for a time/CO2 emissions graph")
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
	Field(2, "Power", Float64, "Power", func() {
		Description("Power meter data in KWh")
	})
	Field(3, "Interval", String, IntervalFunc)
	Required("Duration", "Power", "Interval")
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

var IntervalFunc =  func() {
	Enum(model.Minute, model.Hourly, model.Daily, model.Weekly, model.Monthly)
}

var RegionFunc = func() {
	Enum(model.Caiso, model.Aeso, model.Bpa, model.Erco, model.Ieso, model.Isone, model.Miso, model.Nyiso, model.Nyiso_nycw,
		 model.Nyiso_nyli, model.Nyiso_nyup, model.Pjm, model.Spp)
}

