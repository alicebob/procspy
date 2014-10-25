package procspy

import (
	"encoding/json"
	"os"
)

type jsonFormat struct {
	Connections []Connection
	Processes   []ConnectionProc
}

// Fixture uses the named json file as the fixed output of Connections() and
// Processes(). Used in integration test setups. See example_fixture.json for
// the file format.
func Fixture(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	var res jsonFormat
	if err := json.NewDecoder(f).Decode(&res); err != nil {
		return err
	}

	cbConnections = func() ([]Connection, error) {
		return res.Connections, nil
	}
	cbProcesses = func([]Connection) ([]ConnectionProc, error) {
		return res.Processes, nil
	}

	return nil
}
