package server

import (
	"net/http"
	"time"

	"github.com/jbub/pgbouncer_exporter/collector"
	"github.com/jbub/pgbouncer_exporter/config"
	"github.com/jbub/pgbouncer_exporter/domain"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getLandingPage(telemetryPath string) []byte {
	return []byte(`
	<html>
	<head>
	<title>` + collector.Name + `</title>
	</head>
	<body>
	<h1>` + collector.Name + `</h1>
	<p><a href="` + telemetryPath + `">Metrics</a></p>
	</body>
	</html>`)
}

// New returns new prometheus exporter http server.
func New(cfg config.Config, exp *collector.Exporter, st domain.Store) *HTTPServer {
	reg := collector.NewRegistry(exp)
	mux := newHTTPMux(reg, cfg.TelemetryPath)
	srv := newHTTPServer(cfg.ListenAddress, mux)
	return &HTTPServer{
		cfg: cfg,
		st:  st,
		srv: srv,
	}
}

func newHTTPServer(listenAddr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              listenAddr,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Second,
	}
}

func newHTTPMux(reg prometheus.Gatherer, telemetryPath string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(telemetryPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write(getLandingPage(telemetryPath))
	})
	return mux
}

// HTTPServer represents prometheus exporter http server.
type HTTPServer struct {
	cfg config.Config
	st  domain.Store
	srv *http.Server
}

// Run runs http server.
func (s *HTTPServer) Run() error {
	return s.srv.ListenAndServe()
}
