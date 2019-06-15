package jobs

import (
	"github.com/concepts-system/go-paperless/errors"

	faktory "github.com/contribsys/faktory/client"
)

var client *faktory.Client

// Client returns the singleton job client instance. If not connected, this method will
// create a new client and try to connect to the configured message broker.
func Client() *faktory.Client {
	if client == nil {
		var err error
		client, err = faktory.Open()

		if err != nil {
			panic(errors.Wrapf(err, "Failed to open connection to message-broker"))
		}
	}

	return client
}
