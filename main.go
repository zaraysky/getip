package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// func main() {

// 	http.HandleFunc("/ip", func(w http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintf(w, "%s", strings.Split(req.RemoteAddr, ":")[0])
// 	})
// 	http.ListenAndServe(":8090", nil)
// }

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// )

func main() {
	srv := &http.Server{
		Addr: ":8090",
	}

	idleConnsClosed := make(chan struct{}, 1)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	http.HandleFunc("/ip", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", strings.Split(req.RemoteAddr, ":")[0])
	})

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Print("Server Stopped")
}
