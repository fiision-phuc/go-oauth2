package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/phuc0302/go-oauth2/util"
)

// Recovery recovers server from panic state.
func Recovery(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		var status *util.Status
		if httpError, ok := err.(*util.Status); ok {
			status = httpError
		} else {
			status = util.Status500()
		}

		// Return error
		if redirectURL := redirectPaths[status.Code]; len(redirectURL) > 0 {
			http.Redirect(w, r, redirectURL, status.Code)
		} else {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(status.Code)

			cause, _ := json.Marshal(status)
			w.Write(cause)
		}

		// Slack log
		go func() {
			// Generate error report
			var buffer bytes.Buffer
			buffer.WriteString(fmt.Sprintf("[%d] %s\n", status.Code, status.Description))

			buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "URI", r.URL.Path))
			buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "Address", r.RemoteAddr))
			buffer.WriteString(fmt.Sprintf("%-12s: %s | %s\n", "Method", r.Proto, r.Method))
			buffer.WriteString(fmt.Sprintf("%-12s: %s\n\n", "Request Time", time.Now().UTC().Format(time.RFC822)))

			// Write request
			buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "User Agent", r.UserAgent()))
			buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "Referer", r.Referer()))

			// Write header
			idx := 0
			for header, value := range r.Header {
				if header == "user-agent" || header == "referer" {
					continue
				}

				if idx == 0 {
					buffer.WriteString(fmt.Sprintf("%-12s: [%s] %s\n", "Header", header, value))
				} else {
					buffer.WriteString(fmt.Sprintf("%-12s: [%s] %s\n", "", header, value))
				}
				idx++
			}

			//			// Write Path Params
			//			if c.PathParams != nil && len(c.PathParams) > 0 {
			//				buffer.WriteString("\n")
			//				idx = 0
			//				for key, value := range c.PathParams {
			//					if idx == 0 {
			//						buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "Path Params", key, value))
			//					} else {
			//						buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "", key, value))
			//					}
			//					idx++
			//				}
			//			}

			//			// Write Query Params
			//			if c.QueryParams != nil && len(c.QueryParams) > 0 {
			//				buffer.WriteString("\n")
			//				idx = 0
			//				for key, value := range c.QueryParams {
			//					if idx == 0 {
			//						buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "Query Params", key, value))
			//					} else {
			//						buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "", key, value))
			//					}
			//					idx++
			//				}
			//			}

			//			// Write stack trace
			//			buffer.WriteString("\nStack Trace:\n")
			//			callStack(3, &buffer)

			// Log error
			logrus.Warningln(buffer.String())
		}()
	}
}

// callStack writes stack trace.
func callStack(skip int, w io.Writer) {
	srcPath := fmt.Sprintf("%s/src", os.Getenv("GOPATH"))
	paths := strings.Split(srcPath, ":")

	for _, path := range paths {
		for i := skip; ; i++ {
			// Condition validation: Stop if there is nothing else
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}

			// Condition validation: Skip go root
			if !strings.HasPrefix(file, path) {
				continue
			}

			// Trim prefix
			file = file[(len(path) + 4):]

			// Print this much at least. If we can't find the source, it won't show.
			file = fmt.Sprintf("%s: %s (%d)", file, callFunc(pc), line)
			io.WriteString(w, file)
			io.WriteString(w, "\n")
		}
	}
}

// callFunc returns func's name.
func callFunc(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "???"
	}
	name := fn.Name()

	// Eliminate the path prefix
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		name = name[lastslash+1:]
	}

	//	// Eliminate period prefix
	//	if period := strings.Index(name, "."); period >= 0 {
	//		name = name[period+1:]
	//	}

	//	// Convert center dot to dot
	//	name = strings.Replace(name, "·", ".", -1)
	//	return string(name)

	if tokens := strings.Split(name, "."); len(tokens) >= 2 {
		return tokens[1]
	}
	return "???"
}
