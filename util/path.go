package util

import (
	"fmt"
	"regexp"
)

var (
	// Globs regex
	globsFinder = regexp.MustCompile(`\*\*`)

	// Path regex
	pathFinder = regexp.MustCompile(`{[^/#?()\.\\]+}`)
)

// ConvertPath converts raw path to regular expression rule to match request's path.
//
// - parameter path: url path
func ConvertPath(path string) string {
	regexPattern := pathFinder.ReplaceAllStringFunc(path, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)-1])
	})
	regexPattern = globsFinder.ReplaceAllStringFunc(regexPattern, func(m string) string {
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
	})

	if len(regexPattern) == 1 && regexPattern == "/" {
		regexPattern = fmt.Sprintf("^%s?$", regexPattern)
	} else {
		regexPattern = fmt.Sprintf("^%s/?$", regexPattern)
	}

	return regexPattern
}
