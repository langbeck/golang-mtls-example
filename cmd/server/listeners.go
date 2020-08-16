package main

import (
	"crypto/tls"
	"event-listener/pkg/common"
	"fmt"
	"net"
)

const (
	clientRootCAPath = "creds/ca-client-cert.pem"

	tlsCertPath = "creds/tls-server-cert.pem"
	tlsKeyPath  = "creds/tls-server-key.pem"
)

func listenNoTLS(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func getConfigForTLS() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not load server TLS credentials: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

func listenTLS(network, address string) (net.Listener, error) {
	tlsConf, err := getConfigForTLS()
	if err != nil {
		return nil, err
	}

	return tls.Listen(network, address, tlsConf)
}

func listenMTLS(network, address string) (net.Listener, error) {
	tlsConf, err := getConfigForTLS()
	if err != nil {
		return nil, err
	}

	// Enforce mTLS (no verification)
	tlsConf.ClientAuth = tls.RequireAnyClientCert

	return tls.Listen(network, address, tlsConf)
}

func listenVerifiedMTLS(network, address string) (net.Listener, error) {
	tlsConf, err := getConfigForTLS()
	if err != nil {
		return nil, err
	}

	trusted, err := common.LoadCertPool(clientRootCAPath)
	if err != nil {
		return nil, fmt.Errorf("could not load trusted bundle: %w", err)
	}

	// Enforce mTLS
	tlsConf.ClientAuth = tls.RequireAndVerifyClientCert
	tlsConf.ClientCAs = trusted

	return tls.Listen(network, address, tlsConf)
}
