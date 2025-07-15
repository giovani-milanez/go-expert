package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type CepInput struct {
	Cep string `json:"cep"`
}
var SERVICE_B_URL string
func main() {
	collectorUrl := os.Getenv("OTEL_COLLECTOR_URL")
	if collectorUrl == "" {
		log.Fatal("OTEL_COLLECTOR_URL environment variable is not set")
	}
	SERVICE_B_URL = os.Getenv("SERVICEB_URL")
	if SERVICE_B_URL == "" {
		log.Fatal("SERVICEB_URL environment variable is not set")
	}

	if err := run(collectorUrl); err != nil {
		log.Fatalln(err)
	}
}

func run(collectorUrl string) (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK("service-a", collectorUrl, ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":8081",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	// Register handlers.
	handleFunc("/clima", climaHandler)

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func climaHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("service-a")
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	_, span := tracer.Start(ctx, "get-clima-a")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input CepInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	cep := input.Cep
	if cep == "" || len(cep) != 8 {
		http.Error(w, "Invalid CEP", http.StatusUnprocessableEntity)
		return
	}

	ret, code := callServiceB(cep, ctx)
	if code != http.StatusOK {
		http.Error(w, fmt.Sprintf("Error calling service B: %v", ret), code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(ret)); err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		return
	}
}

func callServiceB(cep string, ctx context.Context) (string, int) {
	url := fmt.Sprintf("http://%s/clima?cep=%s", SERVICE_B_URL, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", http.StatusInternalServerError
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Println("Error response from service B:", string(body))
		return "", resp.StatusCode
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusOK
	}

	return string(body), http.StatusOK
}