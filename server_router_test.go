package oauth2

//import (
//	"os"
//	"testing"
//)

//func Test_addRoute(t *testing.T) {
//	defer os.Remove(ConfigFile)
//	server := DefaultServer()

//	server.addRoute(GET, "/login", func() {})
//	server.addRoute(POST, "login", func() {})
//	server.addRoute(PATCH, "/..//login", func() {})

//	if len(server.routes) != 1 {
//		t.Errorf("Expected only 1 route object but found %d", len(server.routes))
//	}

//	server.Group("/user", func(s *Server) {
//		s.Get("name", func() {})
//		s.Post("/../name", func() {})
//	})
//	if len(server.routes) != 2 {
//		t.Errorf("Expected only 2 routes objects but found %d", len(server.routes))
//	}

//	server.Group("/user/:userId", func(s *Server) {
//		s.Get("name", func() {})
//		s.Post("/../name", func() {})
//	})
//	if len(server.routes) != 3 {
//		t.Errorf("Expected only 3 routes objects but found %d", len(server.routes))
//	}
//}
