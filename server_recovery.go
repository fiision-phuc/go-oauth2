package oauth2

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/phuc0302/go-oauth2/util"
)

// recovery recovers server from panic state.
func recovery(c *RequestContext, isDevelopment bool) {
	if err := recover(); err != nil {
		var status *util.Status
		if httpError, ok := err.(*util.Status); ok {
			status = httpError
		} else {
			status = util.Status500()
		}

		// Generate error report
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("[%d] %s\n", status.Code, status.Description))

		buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "URI", c.Path))
		buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "Address", c.request.RemoteAddr))
		buffer.WriteString(fmt.Sprintf("%-12s: %s | %s\n", "Method", c.request.Proto, c.request.Method))
		buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "Request Time", time.Now().UTC().Format(time.RFC822)))

		// If not development environment, we output here
		if !isDevelopment {
			c.OutputText(status, buffer.String())
		}
		buffer.WriteString("\n")

		// Write request
		buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "User Agent", c.request.UserAgent()))
		buffer.WriteString(fmt.Sprintf("%-12s: %s\n", "Referer", c.request.Referer()))

		// Write header
		idx := 0
		for header, value := range c.Header {
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

		// Write Path Params
		if c.PathParams != nil && len(c.PathParams) > 0 {
			buffer.WriteString("\n")
			idx = 0
			for key, value := range c.PathParams {
				if idx == 0 {
					buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "Path Params", key, value))
				} else {
					buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "", key, value))
				}
				idx++
			}
		}

		// Write Query Params
		if c.QueryParams != nil && len(c.QueryParams) > 0 {
			buffer.WriteString("\n")
			idx = 0
			for key, value := range c.QueryParams {
				if idx == 0 {
					buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "Query Params", key, value))
				} else {
					buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "", key, value))
				}
				idx++
			}
		}

		// Write stack trace
		buffer.WriteString("\nStack Trace:\n")
		callStack(3, &buffer)

		// If development environment, we output here
		if isDevelopment {
			c.OutputText(status, buffer.String())
		}

		// Log error
		logrus.Warningln(buffer.String())

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
	} else {
		return "???"
	}
}
