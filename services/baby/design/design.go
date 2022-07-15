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
	
	HTTP(func(){
		Path("/BabyNames")
	})

	Method("GetName", func() {
		Description("get most popular baby name")
		//Error("server_error", ErrorResult, "Error with Singularity Server.")
		Payload(func() {
			Attribute("year", String, "year") 
		})
        Result(name)
		HTTP(func(){
			GET("/{year}")
		})
	})
})

var name = Type("name", func() {
    Attribute("name", String, "name", func() {
    })
    Required("name")
})
