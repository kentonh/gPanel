// Package gpaccount handles the logic of the gPanel account server
package gpaccount

import (
	"context"
	"errors"
	"log"
)

func (con *Controller) Start() error {
	if con.Status == 1 {
		return errors.New("account server is already on")
	}

	con.Status = 1
	go con.Server.ListenAndServe()
	log.Printf("gPanel Account %v now serving out of %s on port %d\n", con.Name, con.DocumentRoot, con.Port)
	return nil
}

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
