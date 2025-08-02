package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Service2 struct {
	client *http.Client
}

type Response struct {
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data,omitempty"`
}

func main() {
	service := &Service2{}

	// Initialize mTLS client for communicating with service1
	if err := service.initMTLSClient(); err != nil {
		log.Fatalf("Failed to initialize mTLS client: %v", err)
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", service.healthHandler)
	mux.HandleFunc("/metrics", service.metricsHandler)
	mux.HandleFunc("/call-service1", service.callService1Handler)

	// Setup mTLS server
	server := &http.Server{
		Addr:    ":8444",
		Handler: mux,
	}

	tlsConfig, err := setupMTLSServer()
	if err != nil {
		log.Fatalf("Failed to setup mTLS server: %v", err)
	}
	server.TLSConfig = tlsConfig

	log.Println("Service2 starting on :8444 with mTLS...")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func setupMTLSServer() (*tls.Config, error) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair("/certs/service2-cert.pem", "/certs/service2-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %v", err)
	}

	// Load CA certificate
	caCert, err := os.ReadFile("/certs/ca.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		ServerName:   "service2",
	}, nil
}

func (s *Service2) initMTLSClient() error {
	// Load client certificate for service2
	cert, err := tls.LoadX509KeyPair("/certs/service2-client-cert.pem", "/certs/service2-client-key.pem")
	if err != nil {
		return fmt.Errorf("failed to load client certificate: %v", err)
	}

	// Load CA certificate
	caCert, err := os.ReadFile("/certs/ca.pem")
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   "service1",
	}

	s.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 30 * time.Second,
	}

	return nil
}

func (s *Service2) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Service:   "service2",
		Message:   "Service2 is healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service2) metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"requests_processed": 1024,
		"uptime_seconds":     3600,
		"memory_usage_mb":    128,
		"cpu_usage_percent":  25.5,
	}

	response := Response{
		Service:   "service2",
		Message:   "Metrics from service2",
		Timestamp: time.Now(),
		Data:      metrics,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service2) callService1Handler(w http.ResponseWriter, r *http.Request) {
	// Make mTLS call to service1
	resp, err := s.client.Get("https://service1:8443/data")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to call service1: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response from service1: %v", err), http.StatusInternalServerError)
		return
	}

	var service1Response Response
	if err := json.Unmarshal(body, &service1Response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse response from service1: %v", err), http.StatusInternalServerError)
		return
	}

	response := Response{
		Service:   "service2",
		Message:   "Successfully called service1 via mTLS",
		Timestamp: time.Now(),
		Data:      service1Response,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
