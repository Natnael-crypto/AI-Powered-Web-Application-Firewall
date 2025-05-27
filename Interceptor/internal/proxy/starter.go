package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"interceptor/internal/logger"
	"interceptor/internal/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Starter() {
	err := fetchConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	err = fetchApplicationConfig()
	if err != nil {
		log.Fatalf("Failed to fetch applications: %v", err)
		return
	}

	err = utils.InitHttpHandler()
	if err != nil {
		log.Fatalf("Failed to initialize Http Handler: %v", err)
	}

	err = utils.InitMlService()
	if err != nil {
		log.Fatalf("Failed to initialize Ml Handler: %v", err)
	}

	if remoteLogServer != "" {
		err = logger.InitializeLogger(remoteLogServer)
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
		defer logger.CloseLogger()
	}

	httpServer := &http.Server{
		Addr:    "0.0.0.0" + proxyPort,
		Handler: http.HandlerFunc(proxyRequest),
	}

	certMap := make(map[string]tls.Certificate)
	CertApp.mu.Lock()
	for _, cert := range CertApp.Certs {
		if cert.CertPath != "" && cert.KeyPath != "" {
			tlsCert, err := tls.LoadX509KeyPair(cert.CertPath, cert.KeyPath)
			if err != nil {
				log.Printf("Failed to load cert for %s: %v", cert.HostName, err)
				continue
			}
			certMap[cert.HostName] = tlsCert
		}
	}
	CertApp.mu.Unlock()

	getCertificate := func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, exists := certMap[chi.ServerName]; exists {
			return &cert, nil
		}
		log.Printf("No certificate found for %s, serving default certificate", chi.ServerName)
		for _, cert := range certMap {
			return &cert, nil
		}
		return nil, fmt.Errorf("no valid certificate found")
	}

	tlsConfig := &tls.Config{
		GetCertificate: getCertificate,
	}

	httpsListener, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		log.Fatalf("Failed to create HTTPS listener: %v", err)
	}

	httpsServer := &http.Server{
		Handler: http.HandlerFunc(proxyRequest),
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", proxyPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		log.Println("Starting HTTPS server on port 443 with SNI support")
		if err := httpsServer.Serve(httpsListener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	if err := httpsServer.Shutdown(ctx); err != nil {
		log.Printf("HTTPS server shutdown error: %v", err)
	}

	log.Println("Servers gracefully stopped")
}
