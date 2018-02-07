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
		return errors.New("public server is already on")
	}

	con.Status = 1
	go con.Server.ListenAndServe()
	log.Printf("Public server now serving out of %s on port %d\n", con.Directory+"document_root/", con.Port)
	return nil
}

// Stop function stops the server gracefully or forceful, depending on the boolean input
func (con *Controller) Stop(graceful bool) error {
	if graceful {
		context, cancel := context.WithTimeout(context.Background(), con.GracefulShutdownTimeout)
		defer cancel()

		err := con.Server.Shutdown(context)
		if err == nil {
			return nil
		}

		log.Printf("Graceful shutdown failed attempting forced: %v\n", err)
	}

	if err := con.Server.Close(); err != nil {
		return err
	}

	con.Status = 0
	return nil
}
