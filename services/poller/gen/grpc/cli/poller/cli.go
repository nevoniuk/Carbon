// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC client CLI support package
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package cli

import (
	"flag"
	"fmt"
	"os"

	pollerc "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/client"
	goa "goa.design/goa/v3/pkg"
	grpc "google.golang.org/grpc"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `poller (carbon-emissions|aggregate-data)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` poller carbon-emissions --message '{
      "field": [
         "Molestias eveniet doloribus quia ea.",
         "Minus dolores.",
         "Adipisci non rerum nisi quisquam.",
         "Aliquam pariatur sit iure debitis."
      ]
   }'` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(cc *grpc.ClientConn, opts ...grpc.CallOption) (goa.Endpoint, interface{}, error) {
	var (
		pollerFlags = flag.NewFlagSet("poller", flag.ContinueOnError)

		pollerCarbonEmissionsFlags       = flag.NewFlagSet("carbon-emissions", flag.ExitOnError)
		pollerCarbonEmissionsMessageFlag = pollerCarbonEmissionsFlags.String("message", "", "")

		pollerAggregateDataFlags = flag.NewFlagSet("aggregate-data", flag.ExitOnError)
	)
	pollerFlags.Usage = pollerUsage
	pollerCarbonEmissionsFlags.Usage = pollerCarbonEmissionsUsage
	pollerAggregateDataFlags.Usage = pollerAggregateDataUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "poller":
			svcf = pollerFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "poller":
			switch epn {
			case "carbon-emissions":
				epf = pollerCarbonEmissionsFlags

			case "aggregate-data":
				epf = pollerAggregateDataFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "poller":
			c := pollerc.NewClient(cc, opts...)
			switch epn {
			case "carbon-emissions":
				endpoint = c.CarbonEmissions()
				data, err = pollerc.BuildCarbonEmissionsPayload(*pollerCarbonEmissionsMessageFlag)
			case "aggregate-data":
				endpoint = c.AggregateDataEndpoint()
				data = nil
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// pollerUsage displays the usage of the poller command and its subcommands.
func pollerUsage() {
	fmt.Fprintf(os.Stderr, `Service that provides forecasts to clickhouse from Carbonara API
Usage:
    %[1]s [globalflags] poller COMMAND [flags]

COMMAND:
    carbon-emissions: query api getting search data for carbon_intensity event
    aggregate-data: get the aggregate data for an event from clickhouse

Additional help:
    %[1]s poller COMMAND --help
`, os.Args[0])
}
func pollerCarbonEmissionsUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] poller carbon-emissions -message JSON

query api getting search data for carbon_intensity event
    -message JSON: 

Example:
    %[1]s poller carbon-emissions --message '{
      "field": [
         "Molestias eveniet doloribus quia ea.",
         "Minus dolores.",
         "Adipisci non rerum nisi quisquam.",
         "Aliquam pariatur sit iure debitis."
      ]
   }'
`, os.Args[0])
}

func pollerAggregateDataUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] poller aggregate-data

get the aggregate data for an event from clickhouse

Example:
    %[1]s poller aggregate-data
`, os.Args[0])
}
