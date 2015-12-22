package utils

import "net/http"

type Status struct {
	Code        int    `json:"status,omitempty"`
	Error       string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`

	StackTrace interface{} `json:"stack_trace,omitempty"`
}

// MARK: Struct's constructors
func Status200() *Status {
	return genericStatus(http.StatusOK)
}
func Status200WithDescription(description string) *Status {
	status := Status200()
	status.Description = description
	return status
}
func Status201() *Status {
	return genericStatus(http.StatusCreated)
}
func Status201WithDescription(description string) *Status {
	status := Status201()
	status.Description = description
	return status
}
func Status202() *Status {
	return genericStatus(http.StatusAccepted)
}
func Status202WithDescription(description string) *Status {
	status := Status202()
	status.Description = description
	return status
}
func Status203() *Status {
	return genericStatus(http.StatusNonAuthoritativeInfo)
}
func Status203WithDescription(description string) *Status {
	status := Status203()
	status.Description = description
	return status
}
func Status204() *Status {
	return genericStatus(http.StatusNoContent)
}
func Status204WithDescription(description string) *Status {
	status := Status204()
	status.Description = description
	return status
}
func Status205() *Status {
	return genericStatus(http.StatusResetContent)
}
func Status205WithDescription(description string) *Status {
	status := Status205()
	status.Description = description
	return status
}
func Status206() *Status {
	return genericStatus(http.StatusPartialContent)
}
func Status206WithDescription(description string) *Status {
	status := Status206()
	status.Description = description
	return status
}
func Status300() *Status {
	return genericStatus(http.StatusMultipleChoices)
}
func Status300WithDescription(description string) *Status {
	status := Status300()
	status.Description = description
	return status
}
func Status301() *Status {
	return genericStatus(http.StatusMovedPermanently)
}
func Status301WithDescription(description string) *Status {
	status := Status301()
	status.Description = description
	return status
}
func Status302() *Status {
	return genericStatus(http.StatusFound)
}
func Status302WithDescription(description string) *Status {
	status := Status302()
	status.Description = description
	return status
}
func Status303() *Status {
	return genericStatus(http.StatusSeeOther)
}
func Status303WithDescription(description string) *Status {
	status := Status303()
	status.Description = description
	return status
}
func Status304() *Status {
	return genericStatus(http.StatusNotModified)
}
func Status304WithDescription(description string) *Status {
	status := Status304()
	status.Description = description
	return status
}
func Status305() *Status {
	return genericStatus(http.StatusUseProxy)
}
func Status305WithDescription(description string) *Status {
	status := Status305()
	status.Description = description
	return status
}
func Status307() *Status {
	return genericStatus(http.StatusTemporaryRedirect)
}
func Status307WithDescription(description string) *Status {
	status := Status307()
	status.Description = description
	return status
}
func Status400() *Status {
	return genericStatus(http.StatusBadRequest)
}

func Status400WithDescription(description string) *Status {
	status := Status400()
	status.Description = description
	return status
}
func Status401() *Status {
	return genericStatus(http.StatusUnauthorized)
}
func Status401WithDescription(description string) *Status {
	status := Status401()
	status.Description = description
	return status
}
func Status402() *Status {
	return genericStatus(http.StatusPaymentRequired)
}
func Status402WithDescription(description string) *Status {
	status := Status402()
	status.Description = description
	return status
}
func Status403() *Status {
	return genericStatus(http.StatusForbidden)
}
func Status403WithDescription(description string) *Status {
	status := Status403()
	status.Description = description
	return status
}
func Status404() *Status {
	return genericStatus(http.StatusNotFound)
}
func Status404WithDescription(description string) *Status {
	status := Status404()
	status.Description = description
	return status
}
func Status405() *Status {
	return genericStatus(http.StatusMethodNotAllowed)
}
func Status405WithDescription(description string) *Status {
	status := Status405()
	status.Description = description
	return status
}
func Status406() *Status {
	return genericStatus(http.StatusNotAcceptable)
}
func Status406WithDescription(description string) *Status {
	status := Status406()
	status.Description = description
	return status
}
func Status407() *Status {
	return genericStatus(http.StatusProxyAuthRequired)
}
func Status407WithDescription(description string) *Status {
	status := Status407()
	status.Description = description
	return status
}
func Status408() *Status {
	return genericStatus(http.StatusRequestTimeout)
}
func Status408WithDescription(description string) *Status {
	status := Status408()
	status.Description = description
	return status
}
func Status409() *Status {
	return genericStatus(http.StatusConflict)
}
func Status409WithDescription(description string) *Status {
	status := Status409()
	status.Description = description
	return status
}
func Status410() *Status {
	return genericStatus(http.StatusGone)
}
func Status410WithDescription(description string) *Status {
	status := Status410()
	status.Description = description
	return status
}
func Status411() *Status {
	return genericStatus(http.StatusLengthRequired)
}
func Status411WithDescription(description string) *Status {
	status := Status411()
	status.Description = description
	return status
}
func Status412() *Status {
	return genericStatus(http.StatusPreconditionFailed)
}
func Status412WithDescription(description string) *Status {
	status := Status412()
	status.Description = description
	return status
}
func Status413() *Status {
	return genericStatus(http.StatusRequestEntityTooLarge)
}
func Status413WithDescription(description string) *Status {
	status := Status413()
	status.Description = description
	return status
}
func Status414() *Status {
	return genericStatus(http.StatusRequestURITooLong)
}
func Status414WithDescription(description string) *Status {
	status := Status414()
	status.Description = description
	return status
}
func Status415() *Status {
	return genericStatus(http.StatusUnsupportedMediaType)
}
func Status415WithDescription(description string) *Status {
	status := Status415()
	status.Description = description
	return status
}
func Status416() *Status {
	return genericStatus(http.StatusRequestedRangeNotSatisfiable)
}
func Status416WithDescription(description string) *Status {
	status := Status416()
	status.Description = description
	return status
}
func Status417() *Status {
	return genericStatus(http.StatusExpectationFailed)
}
func Status417WithDescription(description string) *Status {
	status := Status417()
	status.Description = description
	return status
}
func Status418() *Status {
	return genericStatus(http.StatusTeapot)
}
func Status418WithDescription(description string) *Status {
	status := Status418()
	status.Description = description
	return status
}
func Status422() *Status {
	return specificStatus(422, "Unprocessable Entity")
}
func Status422WithDescription(description string) *Status {
	status := Status422()
	status.Description = description
	return status
}
func Status423() *Status {
	return specificStatus(423, "Locked")
}
func Status423WithDescription(description string) *Status {
	status := Status423()
	status.Description = description
	return status
}
func Status424() *Status {
	return specificStatus(424, "Failed Dependency")
}
func Status424WithDescription(description string) *Status {
	status := Status424()
	status.Description = description
	return status
}
func Status425() *Status {
	return specificStatus(425, "Unordered Collection")
}
func Status425WithDescription(description string) *Status {
	status := Status425()
	status.Description = description
	return status
}
func Status426() *Status {
	return specificStatus(426, "Upgrade Required")
}
func Status426WithDescription(description string) *Status {
	status := Status426()
	status.Description = description
	return status
}
func Status428() *Status {
	return specificStatus(428, "Precondition Required")
}
func Status428WithDescription(description string) *Status {
	status := Status428()
	status.Description = description
	return status
}
func Status429() *Status {
	return specificStatus(429, "Too Many Requests")
}
func Status429WithDescription(description string) *Status {
	status := Status429()
	status.Description = description
	return status
}
func Status431() *Status {
	return specificStatus(431, "Request Header Fields Too Large")
}
func Status431WithDescription(description string) *Status {
	status := Status431()
	status.Description = description
	return status
}
func Status500() *Status {
	return genericStatus(http.StatusInternalServerError)
}
func Status500WithDescription(description string) *Status {
	status := Status500()
	status.Description = description
	return status
}
func Status501() *Status {
	return genericStatus(http.StatusNotImplemented)
}
func Status501WithDescription(description string) *Status {
	status := Status501()
	status.Description = description
	return status
}
func Status502() *Status {
	return genericStatus(http.StatusBadGateway)
}
func Status502WithDescription(description string) *Status {
	status := Status502()
	status.Description = description
	return status
}
func Status503() *Status {
	return genericStatus(http.StatusServiceUnavailable)
}
func Status503WithDescription(description string) *Status {
	status := Status503()
	status.Description = description
	return status
}
func Status504() *Status {
	return genericStatus(http.StatusGatewayTimeout)
}
func Status504WithDescription(description string) *Status {
	status := Status504()
	status.Description = description
	return status
}
func Status505() *Status {
	return genericStatus(http.StatusHTTPVersionNotSupported)
}
func Status505WithDescription(description string) *Status {
	status := Status505()
	status.Description = description
	return status
}
func Status506() *Status {
	return specificStatus(506, "Variant Also Negotiates")
}
func Status506WithDescription(description string) *Status {
	status := Status506()
	status.Description = description
	return status
}
func Status507() *Status {
	return specificStatus(507, "Insufficient Storage")
}
func Status507WithDescription(description string) *Status {
	status := Status507()
	status.Description = description
	return status
}
func Status508() *Status {
	return specificStatus(508, "Loop Detected")
}
func Status508WithDescription(description string) *Status {
	status := Status508()
	status.Description = description
	return status
}
func Status511() *Status {
	return specificStatus(511, "Network Authentication Required")
}
func Status511WithDescription(description string) *Status {
	status := Status511()
	status.Description = description
	return status
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
