package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"pfi/sensorbee/sensorbee/server/testutil"
)

func createTestServer() *testutil.Server {
	return createTestServerWithCustomRoute(nil)
}

func createTestServerWithCustomRoute(route func(prefix string, r *web.Router)) *testutil.Server {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	topologies := NewDefaultTopologyRegistry()

	root := SetUpRouter("/", ContextGlobalVariables{
		Logger:     logger,
		Topologies: topologies,
	})
	SetUpAPIRouter("/", root, route)
	return testutil.NewServer(root)
}

var jscan = testutil.JScan
