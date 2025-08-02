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

// Configuration from environment variables
type Config struct {
	Port           string
	CertPath       string
	KeyPath        string
	CAPath         string
	ClientCertPath string
	ClientKeyPath  string
	Service2URL    string
}

func getConfig() *Config {
	return &Config{
		Port:           getEnv("PORT", "8443"),
		CertPath:       getEnv("TLS_CERT_PATH", "/certs/service1-cert.pem"),
		KeyPath:        getEnv("TLS_KEY_PATH", "/certs/service1-key.pem"),
		CAPath:         getEnv("CA_CERT_PATH", "/certs/ca.pem"),
		ClientCertPath: getEnv("CLIENT_CERT_PATH", "/certs/service1-client-cert.pem"),
		ClientKeyPath:  getEnv("CLIENT_KEY_PATH", "/certs/service1-client-key.pem"),
		Service2URL:    getEnv("SERVICE2_URL", "https://service2:8444"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	config := getConfig()
	service := &Service1{}
	
	// Initialize mTLS client for communicating with service2
	if err := service.initMTLSClient(config); err != nil {
		log.Fatalf("Failed to initialize mTLS client: %v", err)
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", service.healthHandler)
	mux.HandleFunc("/data", service.dataHandler)
	mux.HandleFunc("/call-service2", func(w http.ResponseWriter, r *http.Request) {
		service.callService2Handler(w, r, config.Service2URL)
	})

	// Setup mTLS server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	tlsConfig, err := setupMTLSServer(config)
	if err != nil {
		log.Fatalf("Failed to setup mTLS server: %v", err)
	}
	server.TLSConfig = tlsConfig

	log.Printf("Service1 starting on :%s with mTLS...", config.Port)
	log.Printf("Service2 URL: %s", config.Service2URL)
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func setupMTLSServer(config *Config) (*tls.Config, error) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %v", err)
	}

	// Load CA certificate
	caCert, err := os.ReadFile(config.CAPath)
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

func (s *Service1) initMTLSClient(config *Config) error {
	// Load client certificate for service1
	cert, err := tls.LoadX509KeyPair(config.ClientCertPath, config.ClientKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load client certificate: %v", err)
	}

	// Load CA certificate
	caCert, err := os.ReadFile(config.CAPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:           caCertPool,
		InsecureSkipVerify: false, // Set to true only for testing
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

func (s *Service1) callService2Handler(w http.ResponseWriter, r *http.Request, service2URL string) {
	// Make mTLS call to service2
	resp, err := s.client.Get(service2URL + "/health")
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