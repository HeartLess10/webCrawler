package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"webScraping/handlers"
)

func main() {
	logger := log.New(os.Stdout, "Crawler-prefix", log.LstdFlags)
	server := getServer(logger)
	go func() {
		logger.Println("Starting server...")
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	setGracefulShutDown(server, logger)
}

func setGracefulShutDown(server *http.Server, logger *log.Logger) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	logger.Println("Recived terminate, graceful shutdown =>", sig)
	tc, _ := context.WithTimeout(context.Background(), 90*time.Second)
	server.Shutdown(tc)
}

func GetServerServeMux(logger *log.Logger) *http.ServeMux {
	sm := http.NewServeMux()
	crawlerHandler := handlers.NewGetCrawlerHandler(logger)
	sm.Handle("/crawler", crawlerHandler)
	return sm
}

func getServer(logger *log.Logger) *http.Server {
	server := new(http.Server)
	server.Addr = ":9090"
	server.Handler = GetServerServeMux(logger)

	return server

}
