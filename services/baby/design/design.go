package design

import . "goa.design/goa/v3/dsl"


var _ = API("BabyNames", func() {
	Title("BabyNames")
	Server("design", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})
var _ = Service("BabyNames", func() {
	Method("GetName", func() {
		Description("get most popular baby name")
		//Error("server_error", ErrorResult, "Error with Singularity Server.")
		Payload(ObtainListPayload)
        Result(Name)
	})
})
var ObtainListPayload = Type("ObtainListPayload", func() {
    Field(1, "year", String, "year", func() {
    })
    Required("year")
})
var Name = Type("Name", func() {
    Field(1, "name", String, "name", func() {
    })
    Required("name")
})
