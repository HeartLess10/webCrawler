package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"webScraping/server/handlers"
)

func StartServer(logger *log.Logger) {
	server := getServer(logger)
	go func() {
		logger.Println("Starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Println("Server error:", err)
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

func getServerServeMux(logger *log.Logger) *http.ServeMux {
	sm := http.NewServeMux()
	setUpHandlers(sm, logger)
	return sm
}

func setUpHandlers(sm *http.ServeMux, logger *log.Logger) {
	sm.Handle("/crawler", handlers.NewCrawlerHandler(logger))
	logger.Println("Done setting up crawler urls")
	sm.Handle("/authentication", handlers.NewAuthenticationHandler(logger))
	logger.Println("Done setting up authentication urls")

}

func getServer(logger *log.Logger) *http.Server {
	server := new(http.Server)
	server.Addr = ":9090"
	server.Handler = getServerServeMux(logger)

	return server
}
