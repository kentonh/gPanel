// Package gpaccount handles the logic of the gPanel account server
package gpaccount

import (
	"context"
	"errors"
	"fmt"
)

func (con *Controller) Start() error {
	if con.Status == 1 {
		return errors.New("Account server is already on.")
	}

	con.Status = 1
	go httpserver.ListenAndServe()
	return nil
}

func (con *Controller) Stop(graceful bool) error {
	if graceful {
		context, cancel := context.WithTimeout(context.Background(), con.GracefulShutdownTimeout)
		defer cancel()

		err := httpserver.Shutdown(context)
		if err == nil {
			return nil
		}

		fmt.Printf("Graceful shutdown failed attempting forced: %v\n", err)
	}

	if err := httpserver.Close(); err != nil {
		return err
	}

	con.Status = 0
	return nil
}
