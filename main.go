package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

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

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
		// indexfile := path + "index.html"
		// if _, err := os.Stat(indexfile); os.IsNotExist(err) {
		// 	var file, _ = os.Create(indexfile)
		// 	file.WriteString("Hello World!")
		// 	file.Sync()
		// }
	}

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
