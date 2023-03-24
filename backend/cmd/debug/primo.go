package debug

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"

	"ariadne/primo"
)

func init() {
	DebugCmd.AddCommand(dumpPrimoAPIResponsesCmd)
	DebugCmd.AddCommand(dumpPrimoFRBRMemberRequestsCmd)
	DebugCmd.AddCommand(dumpPrimoHTTPRequestCmd)
	DebugCmd.AddCommand(dumpPrimoHTTPResponsesCmd)
	DebugCmd.AddCommand(primoLinksJSONCmd)
}

var dumpPrimoAPIResponsesCmd = &cobra.Command{
	Use:     "primo-api-responses [query string]",
	Short:   "Dump Primo API response bodies for query string",
	Example: "ariadne debug primo-api-responses '?sid=&aulast=Shakespeare&aufirst=William&genre=book&title=The%20Oxford%20Shakespeare:%20Hamlet&date=1987&isbn=9780198129103'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		dump, err := dumpPrimoAPIResponses(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(dump)
	},
}

var dumpPrimoFRBRMemberRequestsCmd = &cobra.Command{
	Use:     "primo-frbr-member-requests [query string]",
	Short:   "Dump Primo HTTP requests for query string: all FRBR member requests after the initial ISBN search request",
	Example: "ariadne debug primo-frbr-member-requests '?sid=&aulast=Shakespeare&aufirst=William&genre=book&title=The%20Oxford%20Shakespeare:%20Hamlet&date=1987&isbn=9780198129103'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		dump, err := dumpPrimoFRBRMemberRequests(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(dump)
	},
}

var dumpPrimoHTTPRequestCmd = &cobra.Command{
	Use:     "primo-request [query string]",
	Short:   "Dump Primo HTTP request for query string: initial ISBN search request only",
	Example: "ariadne debug primo-request '?sid=&aulast=Shakespeare&aufirst=William&genre=book&title=The%20Oxford%20Shakespeare:%20Hamlet&date=1987&isbn=9780198129103'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		dump, err := dumpPrimoHTTPRequest(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(dump)
	},
}

var dumpPrimoHTTPResponsesCmd = &cobra.Command{
	Use:     "primo-responses [query string]",
	Short:   "Dump Primo HTTP responses for query string",
	Example: "ariadne debug primo-responses '?sid=&aulast=Shakespeare&aufirst=William&genre=book&title=The%20Oxford%20Shakespeare:%20Hamlet&date=1987&isbn=9780198129103'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		dump, err := dumpPrimoHTTPResponses(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(dump)
	},
}

var primoLinksJSONCmd = &cobra.Command{
	Use:     "primo-links [query string]",
	Short:   "Return JSON array of link objects returned by Primo response for query string",
	Example: "ariadne debug primo-links '?sid=&aulast=Shakespeare&aufirst=William&genre=book&title=The%20Oxford%20Shakespeare:%20Hamlet&date=1987&isbn=9780198129103'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		linksJSON, err := linksJSON(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(linksJSON)
	},
}

func dumpPrimoAPIResponses(queryString string) (string, error) {
	primoRequest, err := primo.NewPrimoRequest(queryString)
	if err != nil {
		return queryString, err
	}

	primoResponse, err := primo.Do(primoRequest)
	if err != nil {
		return queryString, err
	}

	apiResponsesJSON, err := json.MarshalIndent(primoResponse.APIResponses, "", "    ")
	if err != nil {
		return "", err
	}

	return string(apiResponsesJSON), nil
}

func dumpPrimoFRBRMemberRequests(queryString string) (string, error) {
	primoRequest, err := primo.NewPrimoRequest(queryString)
	if err != nil {
		return queryString, err
	}

	primoResponse, err := primo.Do(primoRequest)
	if err != nil {
		return queryString, err
	}

	var output string
	for i, dumpedHTTPRequest := range primoResponse.DumpedHTTPRequests {
		output += formatDumpedHTTPRequestEntry(dumpedHTTPRequest, i)
	}

	return output, nil
}

func dumpPrimoHTTPRequest(queryString string) (string, error) {
	primoRequest, err := primo.NewPrimoRequest(queryString)
	if err != nil {
		return queryString, err
	}

	return primoRequest.DumpedHTTPRequest, nil
}

func dumpPrimoHTTPResponses(queryString string) (string, error) {
	primoRequest, err := primo.NewPrimoRequest(queryString)
	if err != nil {
		return queryString, err
	}

	primoResponse, err := primo.Do(primoRequest)
	if err != nil {
		return queryString, err
	}

	var output string
	for i, dumpedHTTPResponse := range primoResponse.DumpedHTTPResponses {
		output += formatDumpedHTTPResponseEntry(dumpedHTTPResponse, i)
	}

	return output, nil
}

func linksJSON(queryString string) (string, error) {
	primoRequest, err := primo.NewPrimoRequest(queryString)
	if err != nil {
		return queryString, err
	}

	primoResponse, err := primo.Do(primoRequest)
	if err != nil {
		return queryString, err
	}

	linksJSON, err := json.MarshalIndent(primoResponse.Links, "", "    ")
	if err != nil {
		return "", err
	}

	return string(linksJSON), nil
}

func formatDumpedHTTPRequestEntry(dumpedHTTPRequest string, i int) string {
	return formatDumpedEntry("DumpedHTTPRequest", dumpedHTTPRequest, i)
}

func formatDumpedHTTPResponseEntry(dumpedHTTPResponse string, i int) string {
	return formatDumpedEntry("DumpedHTTPResponse", dumpedHTTPResponse, i)
}

func formatDumpedEntry(entryType string, dumpedEntry string, i int) string {
	return fmt.Sprintf(
		`=============================
%s #%d: BEGIN
=============================
%s
=============================
%s #%d: END
=============================


`, entryType, i, dumpedEntry, entryType, i)
}
