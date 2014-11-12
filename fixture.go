package procspy

// SetFixtures declares constant Connection and ConnectionProcs which will
// always be returned by the package-level Connections and Processes
// functions. It's designed to be used in tests.

type fixtConnIter []Connection

func (f *fixtConnIter) Next() *Connection {
	if len(*f) == 0 {
		return nil
	}
	car := (*f)[0]
	*f = (*f)[1:]
	return &car
}

// SetFixtures is used in test scenarios to have known output.
func SetFixtures(c []Connection, p Procs) {
	cbConnections = func() (ConnIter, error) {
		f := fixtConnIter(c)
		return &f, nil
	}

	cbProcesses = func() (Procs, error) {
		return p, nil
	}
}
