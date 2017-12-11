// Package public handles the logic of the public facing website
package public

import (
	"context"
	"errors"
	"log"
)

// Start function starts listening on the public server
func (con *Controller) Start() error {
	if con.Status == 1 {
		return errors.New("Public server is already on.")
	}

	con.Status = 1
	go server.ListenAndServe()
	log.Printf("Public server now serving out of %s on port %d\n", con.Directory+"public/", con.Port)
	return nil
}

// Stop function stops the server gracefully or forceful, depending on the boolean input
func (con *Controller) Stop(graceful bool) error {
	if graceful {
		context, cancel := context.WithTimeout(context.Background(), con.GracefulShutdownTimeout)
		defer cancel()

		err := server.Shutdown(context)
		if err == nil {
			return nil
		}

		log.Printf("Graceful shutdown failed attempting forced: %v\n", err)
	}

	if err := server.Close(); err != nil {
		return err
	}

	con.Status = 0
	return nil
}
