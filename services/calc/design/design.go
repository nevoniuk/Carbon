package design

import (
    . "goa.design/goa/v3/dsl"
    "github.com/crossnokaye/carbon/model"
)

var _ = API("Calc", func() {
    Title("Calc")
	Description("Service to interpret CO2 emissions through KW and carbon intensity data")
    /**
    Server("calc", func() {
        Host("localhost", func() {
            URI("http://localhost:8080")
        })
    })
    */
})

var _ = Service("Calc", func() {
    Description("Service to interpret CO2 emissions through KW and carbon intensity data. Offers the endpoint Historical Carbon Emissions")
    
	Method("historical_carbon_emissions", func() {
        Description("This endpoint is used by a front end service to return carbon emission reports")
        Payload(RequestPayload)
        Result(AllReports)
        Error("reports_not_found", ErrorResult, "Carbon reports not found")
        Error("facility_not_found", ErrorResult, "facilty or location not found")
        GRPC(func() {
			Response(CodeOK)
            Response("facility_not_found", CodeNotFound)
            Response("reports_not_found", CodeNotFound)
        })
    })
})

var AllReports = Type("AllReports", func() {
    Description("CO2 intensity reports, power reports, and CO2 emission reports")
	Field(1, "total_emission_report", EmissionsReport, "TotalEmissionReport")
    Field(2, "carbon_intensity_reports", CarbonReport, "CarbonIntensityReports")
    Field(3, "power_reports", ElectricalReport, "PowerReports")
    Required("total_emission_report")
})

var RequestPayload = Type("RequestPayload", func() {
    Description("Payload wraps the payload for to use the facility config client and poller client")
    Field(1, "org_id", UUID, "OrgID")
    Field(2, "duration", Period, "Duration")
    Field(3, "facility_id", UUID, "FacilityID")
    Field(4, "interval", String, IntervalFunc)
    Field(5, "location_id", UUID, "LocationID")
    Required("org_id", "duration", "interval", "facility_id", "location_id")
})
/**
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
*/
var EmissionsReport = Type("EmissionsReport", func() {
    Description("Carbon/Energy Generation Report")
    Field(1, "duration", Period, "Duration")
    Field(2, "interval", String, IntervalFunc)
    Field(3, "points", ArrayOf(DataPoint), "Points")
    Field(4, "org_id", UUID, "OrgID")
	Field(5, "facility_id", UUID, "FacilityID")
	Field(6, "location_id", UUID, "LocationID")
    Field(7, "region", String, RegionFunc)
    Required("duration", "points", "org_id", "interval", "facility_id", "location_id", "region")
})

var CarbonReport = Type("CarbonReport", func() {
    Description("Carbon Report from clickhouse")
	Field(1, "intensity_points", ArrayOf(DataPoint), "Carbon Intensity points", func() {
		Description("Values are in units of (lbs of CO2/MWh)")
	})
    Field(2, "duration", Period, "Duration")
    Field(3, "interval", String, IntervalFunc)
    Field(4, "region", String, RegionFunc)
    Required("intensity_points", "region", "duration", "interval")
})

var DataPoint = Type("DataPoint", func() {
    Description("Contains carbon emissions in terms of DataPoints, which can be used as points for a time/CO2 emissions graph")
    Field(1, "time", String, "Time", func() {
        Format(FormatDateTime)
        Example("2020-01-01T00:00:00Z")
    })
    Field(2, "value", Float64, "value", func() {
        Example(37.8267)
        Description("either a carbon footprint(lbs of Co2) in a CarbonEmissions struct or power stamp(KW) in an Electrical Report")
    })
    Required("time", "value")
})


var ElectricalReport = Type("ElectricalReport", func() {
    Description("Energy Generation Report from the Past values function GetValues")
    Field(1, "duration", Period, "Duration")
    Field(2, "power_stamps", ArrayOf(DataPoint), "Power Stamps", func() {
        Description("Power meter data in KWh")
    })
    Field(3, "interval", String, IntervalFunc)

    Required("duration", "power_stamps", "interval")
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
