package test

// Template error message
const (
	ExpectedInvalidParameter     = "Expected \"Invalid %s parameter.\" but found \"%s\"."
	ExpectedNil                  = "Expected nil but found not nil."
	ExpectedNotNil               = "Expected not nil but found nil."
	ExpectedPanic                = "Expected panic but never received."
	ExpectedBoolButFoundBool     = "Expected \"%t\" but found \"%t\"."
	ExpectedNumberButFoundNumber = "Expected \"%d\" but found \"%d\"."
	ExpectedStringButFoundString = "Expected \"%s\" but found \"%s\"."
)
