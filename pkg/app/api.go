package app

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Shop struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock bool   `json:"stock"`
	Url   string

	lastUpdate time.Time
}

type Beer struct {
	key    string
	Name   string           `json:"name"`
	Brand  string           `json:"brand"`
	Type   string           `json:"type"`
	Degree int              `json:"degree"`
	Size   int              `json:"size"`
	Shops  map[string]*Shop `json:"shops"`
}

func (a *App) StartServer() {
	// update all data
	a.updateAll()

	// start auto updater
	ctx, updaterCancel := context.WithCancel(context.Background())
	go a.RunUpdater(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.handleDefault)
	mux.HandleFunc("/beers", a.handleBeers)
	mux.Handle("/metrics", promhttp.HandlerFor(a.monitor.Registry, promhttp.HandlerOpts{
		ErrorLog:          a.logger,
		ErrorHandling:     0,
		Registry:          a.monitor.Registry,
		EnableOpenMetrics: true,
	}))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Create a channel to receive signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine
	go func() {
		a.logger.Infof("Server listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for a signal to shutdown the server
	sig := <-signalCh
	a.logger.Debugf("Received signal: %v\n", sig)
	updaterCancel()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		a.logger.Fatalf("Server shutdown failed: %v\n", err)
	}

	a.logger.Infof("Server shutdown gracefully")
}

func (a *App) handleBeers(w http.ResponseWriter, r *http.Request) {
	a.mtx.RLock()
	defer a.mtx.RUnlock()

	res, err := json.Marshal(a.beers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}

	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(res)
}

func (a *App) handleDefault(w http.ResponseWriter, r *http.Request) {
	a.mtx.RLock()
	defer a.mtx.RUnlock()

	tmlData := struct {
		Beers  []*Beer
		OShops []string
	}{
		Beers:  a.beers,
		OShops: a.config.Shops,
	}

	if err := a.htmlTemplate.Execute(w, tmlData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
