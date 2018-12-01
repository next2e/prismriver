package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/server/routes/media"
	"gitlab.com/ttpcodes/prismriver/internal/app/server/routes/player"
	"gitlab.com/ttpcodes/prismriver/internal/app/server/routes/queue"
	"gitlab.com/ttpcodes/prismriver/internal/app/server/ws/routes"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "gitlab.com/ttpcodes/prismriver/statik"
)

func CreateRouter() {
	wait := time.Duration(15)

	r := mux.NewRouter()
	r.HandleFunc("/media/random", media.RandomHandler).Methods("GET")
	r.HandleFunc("/media/search", media.SearchHandler).Methods("GET")
	r.HandleFunc("/player", player.UpdateHandler).Methods("PUT")
	r.HandleFunc("/queue", queue.IndexHandler).Methods("GET")
	r.HandleFunc("/queue", queue.StoreHandler).Methods("POST")
	r.HandleFunc("/queue/{id}", queue.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/ws/queue", routes.WebsocketQueueHandler)

	statikFS, err := fs.New()
	if err != nil {
		logrus.Error("Error on loading static assets:")
		logrus.Error(err)
	}

	r.PathPrefix("/").Handler(http.FileServer(statikFS))

	srv := &http.Server{
		Addr:         "0.0.0.0:80",
		Handler:      r,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error("Error on starting HTTP server:")
			logrus.Error(err)
		}
	}()
	logrus.Info("HTTP server now listening on port 80.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	logrus.Info("HTTP server gracefully shut down.")
}
