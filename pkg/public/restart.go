// Package public handles the logic of the public facing website
package public

import (
	"context"
	"fmt"
)

// Restart function combines both the start and stop function, using different
// status codes, as it is restarting.
func (con *Controller) Restart(graceful bool) error {
	con.Status = 3

	if graceful {
		context, cancel := context.WithTimeout(context.Background(), con.GracefulShutdownTimeout)
		defer cancel()

		err := server.Shutdown(context)
		if err != nil {
			fmt.Printf("Graceful shutdown failed attempting forced: %v\n", err)

			err = server.Close()
			if err != nil {
				return err
			}
		}
	}

	err := server.Close()
	if err != nil {
		return err
	}

	con.Status = 1
	go server.ListenAndServe()
	return nil
}

func (con *Controller) Maintenance() {
	con.Status = 2
}
