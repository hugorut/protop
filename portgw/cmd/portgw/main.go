package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	gwhttp "github.com/hugorut/protop/portgw/internal/http"
	"github.com/hugorut/protop/portgw/internal/processor"
	gwos "github.com/hugorut/protop/portgw/internal/processor/os"
	"github.com/hugorut/protop/portgw/internal/store/grpc"
)

func envOrDefault(env string, def string) string {
	val := os.Getenv(env)
	if strings.TrimSpace(val) == "" {
		return def
	}

	return val
}

func main() {

	port := envOrDefault("HTTP_PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: newRouter(),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-done
	log.Print("server stopping")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %s", err)
	}

	log.Println("server shutdown success")
}

func newRouter() *mux.Router {
	router := mux.NewRouter()

	logger := logrus.New()

	handler := gwhttp.Handler{
		Logger: logger,
		ProcessorProvider: processor.Provider{
			"os": gwos.Processor{
				Processors: 5,
				Open:       os.Open,
				Decode:     gwos.JSONDecode,
				Store: grpc.PortStore{
					Logger: logger,
				},
				Logger: logger,
			},
		},
	}

	// setup all our routes here
	router.HandleFunc("/ports/file/{provider}/upload", handler.ProcessFile).Methods("POST")

	return router
}
