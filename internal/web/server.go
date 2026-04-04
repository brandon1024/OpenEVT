package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ListenAndServe(ctx context.Context, addr, path string, disableExporterMetrics bool) error {
	if !disableExporterMetrics {
		reg.MustRegister(collectors.NewGoCollector())
	}

	mux := http.NewServeMux()

	mux.Handle(path, promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		Registry: reg,
	}))
	mux.Handle("GET /inverter", http.HandlerFunc(GetInverter))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	return server.ListenAndServe()
}

func GetInverter(w http.ResponseWriter, req *http.Request) {
	status := get()

	if status == nil {
		http.Error(w, "inverter status currently unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(status)
}
