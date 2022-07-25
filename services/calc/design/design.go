package design

import (
	. "goa.design/goa/v3/dsl"
	"github.com/crossnokaye/carbon/model"
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
	Field(1, "carbon_intensity_reports", ArrayOf(CarbonReport), "CarbonIntensityReports", func() {
		MinLength(1)
	})
	Field(2, "power_reports", ArrayOf(ElectricalReport), "PowerReports", func() {
		MinLength(1)
	})
	Field(3, "total_emission_report", EmissionsReport, "TotalEmissionReport")
	Required("carbon_intensity_reports", "power_reports", "total_emission_report")
})

var RequestPayload = Type("RequestPayload", func() {
	Description("Payload wraps the payload for to use the facility config client and poller client")
	Field(1, "orgID", UUID, "OrgID")
	Field(2, "duration", Period, "Duration")
	Field(3, "facilityID", UUID, "FacilityID")
	Field(4, "interval", String, IntervalFunc)
	Field(5, "locationID", UUID, "LocationID")
	Required("orgID", "duration", "interval", "facilityID", "locationID")
})

var PastValPayload = Type("PastValPayload", func() {
	Description("Payload wraps the payload for past-values GetValues() and carbon poller service")
	Field(1, "orgID", String, "OrgID")
	Field(2, "duration", Period, "Duration")
	Field(3, "past_val_interval", Int64, "PastValInterval")
	Field(4, "interval", String, IntervalFunc)
	Field(5, "control_point", String, "ControlPoint", func() {
		MinLength(1)
	})
	Field(6, "formula", String, "Formula", func() {
		MinLength(1)
	})
	Field(7, "agent_name", String, "AgentName", func() {
		MinLength(1)
	})
	Required("orgID", "duration", "interval", "past_val_interval", "control_point", "agent_name")
})

var EmissionsReport = Type("EmissionsReport", func() {
	Description("Carbon/Energy Generation Report")
	Field(1, "duration", Period, "Duration")
	Field(2, "interval", String, IntervalFunc)
	Field(3, "points", ArrayOf(DataPoint), "Points")
	Field(4, "orgID", UUID, "OrgID")
	Field(5, "facilityID", UUID, "FacilityID")
	Field(6, "locationID", UUID, "LocationID")
	Field(7, "region", String, RegionFunc)
	Required("duration", "points", "orgID", "interval", "facilityID", "locationID", "region")
})

var CarbonReport = Type("CarbonReport", func() {
	Description("Carbon Report from clickhouse")
	Field(1, "generated_rate", Float64, "GeneratedRate", func() {
		Description("This is in units of (lbs of CO2/MWh)")
	})
	Field(2, "duration", Period, "Duration")
	Field(3, "interval", String, IntervalFunc)
	Field(4, "region", String, RegionFunc)
	Required("generated_rate", "region", "duration", "interval")
})

var DataPoint = Type("DataPoint", func() {
	Description("Contains carbon emissions in terms of DataPoints, which can be used as points for a time/CO2 emissions graph")
	Field(1, "time", String, "Time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Field(2, "carbon_footprint", Float64, "CarbonFootprint", func() {
		Example(37.8267)
		Description("carbon footprint is the lbs of CO2 emissions")
	})

	Required("time", "carbon_footprint")
})

var ElectricalReport = Type("ElectricalReport", func() {
	Description("Energy Generation Report from the Past values function GetValues")
	Field(1, "duration", Period, "Duration")
	Field(2, "power", Float64, "Power", func() {
		Description("Power meter data in KWh")
	})
	Field(3, "interval", String, IntervalFunc)
	Field(4, "payload", PastValPayload, "Payload")
	Required("duration", "power", "interval", "payload")
})

var Period = Type("Period", func() {
	Description("Period of time from start to end for any report type")
	Field(1, "start_time", String, "Start time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Field(2, "end_time", String, "End time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Required("start_time", "end_time")
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

