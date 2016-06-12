package oauth2

import (
	"os"
	"testing"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_DefaultServer(t *testing.T) {
	defer os.Remove(debug)
	server := DefaultServer(true)

	if cfg == nil {
		t.Error(test.ExpectedNotNil)
	}
	if objectFactory == nil {
		t.Error(test.ExpectedNotNil)
	}
	if tokenStore == nil {
		t.Error(test.ExpectedNotNil)
	}
	if redirectPaths == nil {
		t.Error(test.ExpectedNotNil)
	}
	if grantsValidation == nil {
		t.Error(test.ExpectedNotNil)
	}
	if methodsValidation == nil {
		t.Error(test.ExpectedNotNil)
	}
	if server.router == nil {
		t.Error(test.ExpectedNotNil)
	}
}
