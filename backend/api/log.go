package api

import (
	"net/url"
	"strings"
)

type ariadneAPIResponse struct {
	Type        string   `json:"type"`
	APIResponse Response `json:"apiResponse"`
}

type ariadneAPIResponseLogEntry struct {
	sharedLogEntryFields
	APIResponse ariadneAPIResponse `json:"apiResponse"`
}

type primoAPIFRBRMemberRequest struct {
	Type                        string `json:"type"`
	DumpedFRBRMemberHTTPRequest string `json:"dumpedFRBRMemberHTTPRequest"`
}

type primoAPIFRBRMemberResponse struct {
	Type                         string `json:"type"`
	DumpedFRBRMemberHTTPResponse string `json:"dumpedFRBRMemberHTTPResponse"`
}

type primoAPIFRBRMemberRequestLogEntry struct {
	sharedLogEntryFields
	APIRequest primoAPIFRBRMemberRequest `json:"apiRequest"`
}

type primoAPIFRBRMemberResponseLogEntry struct {
	sharedLogEntryFields
	APIResponse primoAPIFRBRMemberResponse `json:"apiResponse"`
}

type primoAPIISBNSearchRequest struct {
	Type                        string `json:"type"`
	DumpedISBNSearchHTTPRequest string `json:"dumpedISBNSearchHTTPRequest"`
}

type primoAPIISBNSearchResponse struct {
	Type                         string `json:"type"`
	DumpedISBNSearchHTTPResponse string `json:"dumpedISBNSearchHTTPResponse"`
}

type primoAPIISBNSearchRequestLogEntry struct {
	sharedLogEntryFields
	APIRequest primoAPIISBNSearchRequest `json:"apiRequest"`
}

type primoAPIISBNSearchResponseLogEntry struct {
	sharedLogEntryFields
	APIResponse primoAPIISBNSearchResponse `json:"apiResponse"`
}

type sharedLogEntryFields struct {
	QueryString string     `json:"queryString"`
	QueryParams url.Values `json:"queryParams"`
}

type sfxAPIRequest struct {
	Type              string `json:"type"`
	DumpedHTTPRequest string `json:"dumpedHTTPRequest"`
}

type sfxAPIResponse struct {
	Type               string `json:"type"`
	DumpedHTTPResponse string `json:"dumpedHTTPResponse"`
}

type sfxAPIRequestLogEntry struct {
	sharedLogEntryFields
	APIRequest sfxAPIRequest `json:"apiRequest"`
}

type sfxAPIResponseLogEntry struct {
	sharedLogEntryFields
	APIResponse sfxAPIResponse `json:"apiResponse"`
}

const AriadneKey = "ariadne"
const MessageKey = "message"

const prefixToTrim = "?"

func getSharedLogEntryFields(queryString string) sharedLogEntryFields {
	if strings.HasPrefix(queryString, prefixToTrim) {
		queryString = strings.TrimPrefix(queryString, prefixToTrim)
	}

	// We are assuming errors are not possible, since bad query strings are caught
	// early.
	params, _ := url.ParseQuery(queryString)

	return sharedLogEntryFields{
		QueryString: queryString,
		QueryParams: params,
	}
}

func makeAriadneAPIResponseLogEntry(queryString string, apiResponse Response) ariadneAPIResponseLogEntry {
	sharedLogEntryFields := getSharedLogEntryFields(queryString)

	return ariadneAPIResponseLogEntry{
		sharedLogEntryFields: sharedLogEntryFields,
		APIResponse: ariadneAPIResponse{
			Type:        "api.Response",
			APIResponse: apiResponse,
		},
	}
}

func makePrimoAPIFRBRMemberRequestLogEntry(queryString string, dumpedHTTPRequest string) primoAPIFRBRMemberRequestLogEntry {
	sharedLogEntryFields := getSharedLogEntryFields(queryString)

	return primoAPIFRBRMemberRequestLogEntry{
		sharedLogEntryFields: sharedLogEntryFields,
		APIRequest: primoAPIFRBRMemberRequest{
			Type:                        "primoRequest",
			DumpedFRBRMemberHTTPRequest: dumpedHTTPRequest,
		},
	}
}

func makePrimoAPIFRBRMemberResponseLogEntry(queryString string, dumpedHTTPResponse string) primoAPIFRBRMemberResponseLogEntry {
	sharedLogEntryFields := getSharedLogEntryFields(queryString)

	return primoAPIFRBRMemberResponseLogEntry{
		sharedLogEntryFields: sharedLogEntryFields,
		APIResponse: primoAPIFRBRMemberResponse{
			Type:                         "primoResponse",
			DumpedFRBRMemberHTTPResponse: dumpedHTTPResponse,
		},
	}
}

func makePrimoAPIISBNSearchRequestLogEntry(queryString string, dumpedHTTPRequest string) primoAPIISBNSearchRequestLogEntry {
	sharedLogEntryFields := getSharedLogEntryFields(queryString)

	return primoAPIISBNSearchRequestLogEntry{
		sharedLogEntryFields: sharedLogEntryFields,
		APIRequest: primoAPIISBNSearchRequest{
			Type:                        "primoRequest",
			DumpedISBNSearchHTTPRequest: dumpedHTTPRequest,
		},
	}
}

func makePrimoAPIISBNSearchResponseLogEntry(queryString string, dumpedHTTPResponse string) primoAPIISBNSearchResponseLogEntry {
	sharedLogEntryFields := getSharedLogEntryFields(queryString)

	return primoAPIISBNSearchResponseLogEntry{
		sharedLogEntryFields: sharedLogEntryFields,
		APIResponse: primoAPIISBNSearchResponse{
			Type:                         "primoResponse",
			DumpedISBNSearchHTTPResponse: dumpedHTTPResponse,
		},
	}
}

func makeNewSFXAPIRequestLogEntry(queryString string, dumpedHTTPRequest string) sfxAPIRequestLogEntry {
	sharedLogEntryFields := getSharedLogEntryFields(queryString)

	return sfxAPIRequestLogEntry{
		sharedLogEntryFields: sharedLogEntryFields,
		APIRequest: sfxAPIRequest{
			Type:              "sfxRequest",
			DumpedHTTPRequest: dumpedHTTPRequest,
		},
	}
}

func makeNewSFXAPIResponseLogEntry(queryString string, dumpedHTTPResponse string) sfxAPIResponseLogEntry {
	sharedLogEntryFields := getSharedLogEntryFields(queryString)

	return sfxAPIResponseLogEntry{
		sharedLogEntryFields: sharedLogEntryFields,
		APIResponse: sfxAPIResponse{
			Type:               "sfxResponse",
			DumpedHTTPResponse: dumpedHTTPResponse,
		},
	}
}
