// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Data Service API gRPC client CLI support package
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package cli

import (
	"flag"
	"fmt"
	"os"

	datac "github.com/crossnokaye/carbon/services/poller/gen/grpc/data/client"
	goa "goa.design/goa/v3/pkg"
	grpc "google.golang.org/grpc"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `data (carbon-emissions|fuels|aggregate-data)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` data carbon-emissions` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(cc *grpc.ClientConn, opts ...grpc.CallOption) (goa.Endpoint, interface{}, error) {
	var (
		dataFlags = flag.NewFlagSet("data", flag.ContinueOnError)

		dataCarbonEmissionsFlags = flag.NewFlagSet("carbon-emissions", flag.ExitOnError)

		dataFuelsFlags = flag.NewFlagSet("fuels", flag.ExitOnError)

		dataAggregateDataFlags = flag.NewFlagSet("aggregate-data", flag.ExitOnError)
	)
	dataFlags.Usage = dataUsage
	dataCarbonEmissionsFlags.Usage = dataCarbonEmissionsUsage
	dataFuelsFlags.Usage = dataFuelsUsage
	dataAggregateDataFlags.Usage = dataAggregateDataUsage

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
		case "data":
			svcf = dataFlags
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
		case "data":
			switch epn {
			case "carbon-emissions":
				epf = dataCarbonEmissionsFlags

			case "fuels":
				epf = dataFuelsFlags

			case "aggregate-data":
				epf = dataAggregateDataFlags

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
		case "data":
			c := datac.NewClient(cc, opts...)
			switch epn {
			case "carbon-emissions":
				endpoint = c.CarbonEmissions()
				data = nil
			case "fuels":
				endpoint = c.Fuels()
				data = nil
			case "aggregate-data":
				endpoint = c.AggregateData()
				data = nil
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// dataUsage displays the usage of the data command and its subcommands.
func dataUsage() {
	fmt.Fprintf(os.Stderr, `Service that provides forecasts to clickhouse from Carbonara API
Usage:
    %[1]s [globalflags] data COMMAND [flags]

COMMAND:
    carbon-emissions: query api getting search data for carbon_intensity event
    fuels: query api using a search call for a fuel event from Carbonara API
    aggregate-data: get the aggregate data for an event from clickhouse

Additional help:
    %[1]s data COMMAND --help
`, os.Args[0])
}
func dataCarbonEmissionsUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] data carbon-emissions

query api getting search data for carbon_intensity event

Example:
    %[1]s data carbon-emissions
`, os.Args[0])
}

func dataFuelsUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] data fuels

query api using a search call for a fuel event from Carbonara API

Example:
    %[1]s data fuels
`, os.Args[0])
}

func dataAggregateDataUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] data aggregate-data

get the aggregate data for an event from clickhouse

Example:
    %[1]s data aggregate-data
`, os.Args[0])
}
