package main

import (
	"log"
	"microservice/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/context"
)

func main() {
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)

	// create the handlers
	ph := handlers.NewProducts(l)

	// create a new serve-mux(router) and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", ph)

	// create a new server
	server := &http.Server{
		Addr:         ":8080",           // configure the address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 1 * time.Second,   // max time to write response to the client
	}

	// start the server
	go func() {
		l.Println("Starting the server on port 8080...")
		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting the server: %s\n...", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// block until a signal is received
	sig := <-sigChan
	l.Println("Got signal:", sig)

	// gracefully shutdown the serve, waiting max 30 seconds for current operations to coomplete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}

// git push -u origin <branch>
