package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func handlers() {
	routerPath("/api", httpGet, func(w http.ResponseWriter, r *http.Request) {
		log.Print("in")
		// composite response body
		var res = map[string]string{"result": "succ", "name": "test"}
		response, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

}

var (
	router   *mux.Router
	hostname string
	port     int
	path     string
	stopChan os.Signal
)

func init() {
	flag.StringVar(&hostname, "hostname", "0.0.0.0", "The hostname or IP on which the REST server will listen")
	flag.IntVar(&port, "port", 80, "The port on which the REST server will listen")
	flag.StringVar(&path, "path", "dist/", "Http Document Path")
	router = mux.NewRouter().StrictSlash(true)
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	flag.Parse()

	handlers()

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(path))))

	var address = fmt.Sprintf("%s:%d", hostname, port)

	srv := &http.Server{Addr: address, Handler: router}
	go func() {
		// service connections
		log.Printf("Server Start on http://%s:%d", hostname, port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan

	log.Println("Shutting down server...")
	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	if cf != nil {
	}
	srv.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}

type httpMethod string

const (
	httpGet  httpMethod = "GET"
	httpPost httpMethod = "POST"
)

func routerPath(path string, method httpMethod, handleFunc func(http.ResponseWriter, *http.Request)) {
	router.HandleFunc(path, handleFunc).Methods(string(method))
}
