package procspy

// SetFixtures declares constant Connection and ConnectionProcs which will
// always be returned by the package-level Connections and Processes
// functions. It's designed to be used in tests.
func SetFixtures(c []Connection, p []ConnectionProc) {
	cbConnections = func() ([]Connection, error) {
		return c, nil
	}

	cbProcesses = func([]Connection) ([]ConnectionProc, error) {
		return p, nil
	}
}
