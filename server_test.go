package oauth2

import "testing"

func Test_DefaultServer(t *testing.T) {
	server := DefaultServer()

	if server.Config == nil {
		t.Errorf("Expected config file should be loaded at creation time but found nil.")
	}
}
