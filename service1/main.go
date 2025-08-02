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

type Service1 struct {
	client *http.Client
}

type Response struct {
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data,omitempty"`
}

func main() {
	service := &Service1{}

	// Initialize mTLS client for communicating with service2
	if err := service.initMTLSClient(); err != nil {
		log.Fatalf("Failed to initialize mTLS client: %v", err)
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", service.healthHandler)
	mux.HandleFunc("/data", service.dataHandler)
	mux.HandleFunc("/call-service2", service.callService2Handler)

	// Setup mTLS server
	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}

	tlsConfig, err := setupMTLSServer()
	if err != nil {
		log.Fatalf("Failed to setup mTLS server: %v", err)
	}
	server.TLSConfig = tlsConfig

	log.Println("Service1 starting on :8443 with mTLS...")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func setupMTLSServer() (*tls.Config, error) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair("/certs/service1-cert.pem", "/certs/service1-key.pem")
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
		ServerName:   "service1",
	}, nil
}

func (s *Service1) initMTLSClient() error {
	// Load client certificate for service1
	cert, err := tls.LoadX509KeyPair("/certs/service1-client-cert.pem", "/certs/service1-client-key.pem")
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
		ServerName:   "service2",
	}

	s.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 30 * time.Second,
	}

	return nil
}

func (s *Service1) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Service:   "service1",
		Message:   "Service1 is healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service1) dataHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"users":  []string{"alice", "bob", "charlie"},
		"status": "active",
		"count":  3,
	}

	response := Response{
		Service:   "service1",
		Message:   "Data from service1",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service1) callService2Handler(w http.ResponseWriter, r *http.Request) {
	// Make mTLS call to service2
	resp, err := s.client.Get("https://service2:8444/health")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to call service2: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response from service2: %v", err), http.StatusInternalServerError)
		return
	}

	var service2Response Response
	if err := json.Unmarshal(body, &service2Response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse response from service2: %v", err), http.StatusInternalServerError)
		return
	}

	response := Response{
		Service:   "service1",
		Message:   "Successfully called service2 via mTLS",
		Timestamp: time.Now(),
		Data:      service2Response,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
