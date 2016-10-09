package oauth2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type log_message struct {
	Uri         string `json:"uri,omitempty"`
	Method      string `json:"method,omitempty"`
	RequestTime string `json:"request_time,omitempty"`

	Trace   []string     `json:"trace,omitempty"`
	Body    *log_body    `json:"body,omitempty"`
	Request *log_request `json:"request,omitempty"`
}

type log_body struct {
	ContentType   string            `json:"content_type,omitempty"`
	RequestBody   map[string]string `json:"request_body,omitempty"`
	RequestParams map[string]string `json:"request_params,omitempty"`
}

type log_request struct {
	UserAgent     string `json:"user_agent,omitempty"`
	HttpReferer   string `json:"http_referer,omitempty"`
	RemoteAddress string `json:"remote_address,omitempty"`
}

/** Recovery function when there is a panic. */
func RecoveryInternal(logger *log.Logger) {
	if err := recover(); err != nil {
		log := log_message{Trace: callStack(3)}

		cause, _ := json.Marshal(log)
		logger.Println(string(cause))
	}
}
func RecoveryRequest(request *http.Request, response http.ResponseWriter, isDevelopment bool) {
	//	if err := recover(); err != nil {
	//		log := log_message{
	//			Uri:         c.Path,
	//			Method:      fmt.Sprintf("%s|%s", c.request.Proto, c.request.Method),
	//			RequestTime: time.Now().UTC().Format(time.RFC822),

	//			Trace: callStack(3),

	//			Body: &log_body{
	//				ContentType:   c.Header["content-type"],
	//				RequestBody:   c.QueryParams,
	//				RequestParams: c.PathParams,
	//			},

	//			Request: &log_request{
	//				UserAgent:     c.request.UserAgent(),
	//				HttpReferer:   c.request.Referer(),
	//				RemoteAddress: c.request.RemoteAddr,
	//			},
	//		}

	//		// Define status error
	//		var httpError *util.Status
	//		if status, ok := reflect.ValueOf(err).Interface().(util.Status); ok {
	//			httpError = &status
	//		} else {
	//			httpError = util.Status500()
	//			httpError.Description = fmt.Sprintf("%s", err)
	//		}

	//		// Should include stack trace or not
	//		if isDevelopment {
	//			httpError.StackTrace = log
	//		}

	//		c.OutputError(httpError)
	//	}
}

// MARK: Private functions
func callStack(skip int) []string {
	// FIX FIX FIX: What if we have more than 1 go path???
	srcPath := fmt.Sprintf("%s/src", os.Getenv("GOPATH"))
	traces := make([]string, 5)

	for i, j := skip, 0; ; i++ {
		// Condition validation: Stop if there is nothing else
		pc, file, line, ok := runtime.Caller(i)
		if !ok || j >= 5 {
			break
		}
		fmt.Println(file, line)

		// Condition validation: Skip go root
		if !strings.HasPrefix(file, srcPath) {
			continue
		}

		// Trim prefix
		file = file[len(srcPath):]

		// Print this much at least. If we can't find the source, it won't show.
		traces[j] = fmt.Sprintf("%s: %s (%d)", file, callFunction(pc), line)
		j++
	}
	return traces
}
func callFunction(pc uintptr) string {
	fn := runtime.FuncForPC(pc)

	// Condition validation: return don't know if function is not available
	if fn == nil {
		return "???"
	}

	// Convert function name to byte array for modification
	name := []byte(fn.Name())

	// Eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}

	// Eliminate period prefix
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}

	// Convert center dot to dot
	name = bytes.Replace(name, centerDot, dot, -1)
	return string(name)
}
