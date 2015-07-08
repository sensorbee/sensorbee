package topology

import (
	"github.com/Sirupsen/logrus"
	"pfi/sensorbee/sensorbee/server"
	"pfi/sensorbee/sensorbee/server/testutil"
)

func createTestServer() *testutil.Server {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	topologies := server.NewDefaultTopologyRegistry()

	root := server.SetUpRouter("/", server.ContextGlobalVariables{
		Logger:     logger,
		Topologies: topologies,
	})
	server.SetUpAPIRouter("/", root, nil)

	testutil.TestAPIWithRealHTTPServer = true
	return testutil.NewServer(root)
}
