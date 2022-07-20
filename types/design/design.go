package design

import (
	//"github.com/crossnokaye/carbon/types"
	. "goa.design/goa/v3/dsl"
)

var IntervalType = Type("IntervalType", func() {
	Description("Defines the interval or 'granularity' types by which to obtain Carbon Intensity reports and Power reports")
	Field(1, "Kind", String, "Interval type kind", func() {
		Enum("minute", "hourly", "daily", "weekly", "monthly")
		Required("minute", "hourly", "daily", "weekly", "monthly")
	})
	
	
})

var RegionName = Type("RegionName", func() {
	Description("Defines the region by which to obtain Carbon Intensity reports and Power reports")
	Field(1, "Region", String, "Acceptable region name", func() {
		Enum("CAISO", "AESO", "BPA", "ERCO", "IESO", "ISONE", "MISO","NYISO", "NYISO.NYCW", "NYISO.NYLI", "NYISO.NYUP", "PJM", "SPP")
		Required("CAISO", "AESO", "BPA", "ERCO", "IESO", "ISONE", "MISO","NYISO", "NYISO.NYCW", "NYISO.NYLI", "NYISO.NYUP", "PJM", "SPP")
	})

	
})