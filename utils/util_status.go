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
func Status201() *Status {
	return genericStatus(http.StatusCreated)
}
func Status202() *Status {
	return genericStatus(http.StatusAccepted)
}
func Status203() *Status {
	return genericStatus(http.StatusNonAuthoritativeInfo)
}
func Status204() *Status {
	return genericStatus(http.StatusNoContent)
}
func Status205() *Status {
	return genericStatus(http.StatusResetContent)
}
func Status206() *Status {
	return genericStatus(http.StatusPartialContent)
}

func Status300() *Status {
	return genericStatus(http.StatusMultipleChoices)
}
func Status301() *Status {
	return genericStatus(http.StatusMovedPermanently)
}
func Status302() *Status {
	return genericStatus(http.StatusFound)
}
func Status303() *Status {
	return genericStatus(http.StatusSeeOther)
}
func Status304() *Status {
	return genericStatus(http.StatusNotModified)
}
func Status305() *Status {
	return genericStatus(http.StatusUseProxy)
}
func Status307() *Status {
	return genericStatus(http.StatusTemporaryRedirect)
}

func Status400() *Status {
	return genericStatus(http.StatusBadRequest)
}
func Status400WithDescription(description string) *Status {
	status := genericStatus(http.StatusBadRequest)
	status.Description = description
	return status
}
func Status401() *Status {
	return genericStatus(http.StatusUnauthorized)
}
func Status401WithDescription(description string) *Status {
	status := genericStatus(http.StatusBadRequest)
	status.Description = description
	return status
}
func Status402() *Status {
	return genericStatus(http.StatusPaymentRequired)
}
func Status403() *Status {
	return genericStatus(http.StatusForbidden)
}
func Status404() *Status {
	return genericStatus(http.StatusNotFound)
}
func Status405() *Status {
	return genericStatus(http.StatusMethodNotAllowed)
}
func Status406() *Status {
	return genericStatus(http.StatusNotAcceptable)
}
func Status407() *Status {
	return genericStatus(http.StatusProxyAuthRequired)
}
func Status408() *Status {
	return genericStatus(http.StatusRequestTimeout)
}
func Status409() *Status {
	return genericStatus(http.StatusConflict)
}
func Status410() *Status {
	return genericStatus(http.StatusGone)
}
func Status411() *Status {
	return genericStatus(http.StatusLengthRequired)
}
func Status412() *Status {
	return genericStatus(http.StatusPreconditionFailed)
}
func Status413() *Status {
	return genericStatus(http.StatusRequestEntityTooLarge)
}
func Status414() *Status {
	return genericStatus(http.StatusRequestURITooLong)
}
func Status415() *Status {
	return genericStatus(http.StatusUnsupportedMediaType)
}
func Status416() *Status {
	return genericStatus(http.StatusRequestedRangeNotSatisfiable)
}
func Status417() *Status {
	return genericStatus(http.StatusExpectationFailed)
}
func Status418() *Status {
	return genericStatus(http.StatusTeapot)
}
func Status422() *Status {
	return specificStatus(422, "Unprocessable Entity")
}
func Status423() *Status {
	return specificStatus(423, "Locked")
}
func Status424() *Status {
	return specificStatus(424, "Failed Dependency")
}
func Status425() *Status {
	return specificStatus(425, "Unordered Collection")
}
func Status426() *Status {
	return specificStatus(426, "Upgrade Required")
}
func Status428() *Status {
	return specificStatus(428, "Precondition Required")
}
func Status429() *Status {
	return specificStatus(429, "Too Many Requests")
}
func Status431() *Status {
	return specificStatus(431, "Request Header Fields Too Large")
}

func Status500() *Status {
	return genericStatus(http.StatusInternalServerError)
}
func Status501() *Status {
	return genericStatus(http.StatusNotImplemented)
}
func Status502() *Status {
	return genericStatus(http.StatusBadGateway)
}
func Status503() *Status {
	return genericStatus(http.StatusServiceUnavailable)
}
func Status504() *Status {
	return genericStatus(http.StatusGatewayTimeout)
}
func Status505() *Status {
	return genericStatus(http.StatusHTTPVersionNotSupported)
}
func Status506() *Status {
	return specificStatus(506, "Variant Also Negotiates")
}
func Status507() *Status {
	return specificStatus(507, "Insufficient Storage")
}
func Status508() *Status {
	return specificStatus(508, "Loop Detected")
}
func Status511() *Status {
	return specificStatus(511, "Network Authentication Required")
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
