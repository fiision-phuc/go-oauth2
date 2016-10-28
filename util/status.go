package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/// Status describes a HTTP status component.
type Status struct {
	Code        int         `json:"status,omitempty"`
	Error       string      `json:"error,omitempty"`
	Description string      `json:"error_description,omitempty"`
	StackTrace  interface{} `json:"stack_trace,omitempty"`
}

func Status200() *Status {
	return genericStatus(http.StatusOK)
}
func Status200WithDescription(description string) *Status {
	return specificStatus(http.StatusOK, description)
}
func Status201() *Status {
	return genericStatus(http.StatusCreated)
}
func Status201WithDescription(description string) *Status {
	return specificStatus(http.StatusCreated, description)
}
func Status202() *Status {
	return genericStatus(http.StatusAccepted)
}
func Status202WithDescription(description string) *Status {
	return specificStatus(http.StatusAccepted, description)
}
func Status203() *Status {
	return genericStatus(http.StatusNonAuthoritativeInfo)
}
func Status203WithDescription(description string) *Status {
	return specificStatus(http.StatusNonAuthoritativeInfo, description)
}
func Status204() *Status {
	return genericStatus(http.StatusNoContent)
}
func Status204WithDescription(description string) *Status {
	return specificStatus(http.StatusNoContent, description)
}
func Status205() *Status {
	return genericStatus(http.StatusResetContent)
}
func Status205WithDescription(description string) *Status {
	return specificStatus(http.StatusResetContent, description)
}
func Status206() *Status {
	return genericStatus(http.StatusPartialContent)
}
func Status206WithDescription(description string) *Status {
	return specificStatus(http.StatusPartialContent, description)
}
func Status300() *Status {
	return genericStatus(http.StatusMultipleChoices)
}
func Status300WithDescription(description string) *Status {
	return specificStatus(http.StatusMultipleChoices, description)
}
func Status301() *Status {
	return genericStatus(http.StatusMovedPermanently)
}
func Status301WithDescription(description string) *Status {
	return specificStatus(http.StatusMovedPermanently, description)
}
func Status302() *Status {
	return genericStatus(http.StatusFound)
}
func Status302WithDescription(description string) *Status {
	return specificStatus(http.StatusFound, description)
}
func Status303() *Status {
	return genericStatus(http.StatusSeeOther)
}
func Status303WithDescription(description string) *Status {
	return specificStatus(http.StatusSeeOther, description)
}
func Status304() *Status {
	return genericStatus(http.StatusNotModified)
}
func Status304WithDescription(description string) *Status {
	return specificStatus(http.StatusNotModified, description)
}
func Status305() *Status {
	return genericStatus(http.StatusUseProxy)
}
func Status305WithDescription(description string) *Status {
	return specificStatus(http.StatusUseProxy, description)
}
func Status307() *Status {
	return genericStatus(http.StatusTemporaryRedirect)
}
func Status307WithDescription(description string) *Status {
	return specificStatus(http.StatusTemporaryRedirect, description)
}
func Status400() *Status {
	return genericStatus(http.StatusBadRequest)
}

func Status400WithDescription(description string) *Status {
	return specificStatus(http.StatusBadRequest, description)
}
func Status401() *Status {
	return genericStatus(http.StatusUnauthorized)
}
func Status401WithDescription(description string) *Status {
	return specificStatus(http.StatusUnauthorized, description)
}
func Status402() *Status {
	return genericStatus(http.StatusPaymentRequired)
}
func Status402WithDescription(description string) *Status {
	return specificStatus(http.StatusPaymentRequired, description)
}
func Status403() *Status {
	return genericStatus(http.StatusForbidden)
}
func Status403WithDescription(description string) *Status {
	return specificStatus(http.StatusForbidden, description)
}
func Status404() *Status {
	return genericStatus(http.StatusNotFound)
}
func Status404WithDescription(description string) *Status {
	return specificStatus(http.StatusNotFound, description)
}
func Status405() *Status {
	return genericStatus(http.StatusMethodNotAllowed)
}
func Status405WithDescription(description string) *Status {
	return specificStatus(http.StatusMethodNotAllowed, description)
}
func Status406() *Status {
	return genericStatus(http.StatusNotAcceptable)
}
func Status406WithDescription(description string) *Status {
	return specificStatus(http.StatusNotAcceptable, description)
}
func Status407() *Status {
	return genericStatus(http.StatusProxyAuthRequired)
}
func Status407WithDescription(description string) *Status {
	return specificStatus(http.StatusProxyAuthRequired, description)
}
func Status408() *Status {
	return genericStatus(http.StatusRequestTimeout)
}
func Status408WithDescription(description string) *Status {
	return specificStatus(http.StatusRequestTimeout, description)
}
func Status409() *Status {
	return genericStatus(http.StatusConflict)
}
func Status409WithDescription(description string) *Status {
	return specificStatus(http.StatusConflict, description)
}
func Status410() *Status {
	return genericStatus(http.StatusGone)
}
func Status410WithDescription(description string) *Status {
	return specificStatus(http.StatusGone, description)
}
func Status411() *Status {
	return genericStatus(http.StatusLengthRequired)
}
func Status411WithDescription(description string) *Status {
	return specificStatus(http.StatusLengthRequired, description)
}
func Status412() *Status {
	return genericStatus(http.StatusPreconditionFailed)
}
func Status412WithDescription(description string) *Status {
	return specificStatus(http.StatusPreconditionFailed, description)
}
func Status413() *Status {
	return genericStatus(http.StatusRequestEntityTooLarge)
}
func Status413WithDescription(description string) *Status {
	return specificStatus(http.StatusRequestEntityTooLarge, description)
}
func Status414() *Status {
	return genericStatus(http.StatusRequestURITooLong)
}
func Status414WithDescription(description string) *Status {
	return specificStatus(http.StatusRequestURITooLong, description)
}
func Status415() *Status {
	return genericStatus(http.StatusUnsupportedMediaType)
}
func Status415WithDescription(description string) *Status {
	return specificStatus(http.StatusUnsupportedMediaType, description)
}
func Status416() *Status {
	return genericStatus(http.StatusRequestedRangeNotSatisfiable)
}
func Status416WithDescription(description string) *Status {
	return specificStatus(http.StatusRequestedRangeNotSatisfiable, description)
}
func Status417() *Status {
	return genericStatus(http.StatusExpectationFailed)
}
func Status417WithDescription(description string) *Status {
	return specificStatus(http.StatusExpectationFailed, description)
}
func Status418() *Status {
	return genericStatus(http.StatusTeapot)
}
func Status418WithDescription(description string) *Status {
	return specificStatus(http.StatusTeapot, description)
}
func Status422() *Status {
	return specificStatus(422, "Unprocessable Entity")
}
func Status422WithDescription(description string) *Status {
	return specificStatus(422, description)
}
func Status423() *Status {
	return specificStatus(423, "Locked")
}
func Status423WithDescription(description string) *Status {
	return specificStatus(423, description)
}
func Status424() *Status {
	return specificStatus(424, "Failed Dependency")
}
func Status424WithDescription(description string) *Status {
	return specificStatus(424, description)
}
func Status425() *Status {
	return specificStatus(425, "Unordered Collection")
}
func Status425WithDescription(description string) *Status {
	return specificStatus(425, description)
}
func Status426() *Status {
	return specificStatus(426, "Upgrade Required")
}
func Status426WithDescription(description string) *Status {
	return specificStatus(426, description)
}
func Status428() *Status {
	return specificStatus(428, "Precondition Required")
}
func Status428WithDescription(description string) *Status {
	return specificStatus(428, description)
}
func Status429() *Status {
	return specificStatus(429, "Too Many Requests")
}
func Status429WithDescription(description string) *Status {
	return specificStatus(429, description)
}
func Status431() *Status {
	return specificStatus(431, "Request Header Fields Too Large")
}
func Status431WithDescription(description string) *Status {
	return specificStatus(431, description)
}
func Status500() *Status {
	return genericStatus(http.StatusInternalServerError)
}
func Status500WithDescription(description string) *Status {
	return specificStatus(http.StatusInternalServerError, description)
}
func Status501() *Status {
	return genericStatus(http.StatusNotImplemented)
}
func Status501WithDescription(description string) *Status {
	return specificStatus(http.StatusNotImplemented, description)
}
func Status502() *Status {
	return genericStatus(http.StatusBadGateway)
}
func Status502WithDescription(description string) *Status {
	return specificStatus(http.StatusBadGateway, description)
}
func Status503() *Status {
	return genericStatus(http.StatusServiceUnavailable)
}
func Status503WithDescription(description string) *Status {
	return specificStatus(http.StatusServiceUnavailable, description)
}
func Status504() *Status {
	return genericStatus(http.StatusGatewayTimeout)
}
func Status504WithDescription(description string) *Status {
	return specificStatus(http.StatusGatewayTimeout, description)
}
func Status505() *Status {
	return genericStatus(http.StatusHTTPVersionNotSupported)
}
func Status505WithDescription(description string) *Status {
	return specificStatus(http.StatusHTTPVersionNotSupported, description)
}
func Status506() *Status {
	return specificStatus(506, "Variant Also Negotiates")
}
func Status506WithDescription(description string) *Status {
	return specificStatus(506, description)
}
func Status507() *Status {
	return specificStatus(507, "Insufficient Storage")
}
func Status507WithDescription(description string) *Status {
	return specificStatus(507, description)
}
func Status508() *Status {
	return specificStatus(508, "Loop Detected")
}
func Status508WithDescription(description string) *Status {
	return specificStatus(508, description)
}
func Status511() *Status {
	return specificStatus(511, "Network Authentication Required")
}
func Status511WithDescription(description string) *Status {
	return specificStatus(511, description)
}

// ParseStatus parses data into status object.
func ParseStatus(response *http.Response) *Status {
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	status := Status{}
	json.Unmarshal(data, &status)

	return &status
}

// MARK: Struct's private constructors
func genericStatus(statusCode int) *Status {
	return &Status{
		Code:        statusCode,
		Error:       http.StatusText(statusCode),
		Description: http.StatusText(statusCode),
	}
}
func specificStatus(statusCode int, title string) *Status {
	return &Status{
		Code:        statusCode,
		Error:       title,
		Description: title,
	}
}
