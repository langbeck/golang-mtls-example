package main

import (
	"crypto/tls"
	"event-listener/pkg/common"
	"fmt"
	"net"
)

const (
	serverRootCAPath = "creds/ca-server-cert.pem"

	tlsCertPath = "creds/tls-client-cert.pem"
	tlsKeyPath  = "creds/tls-client-key.pem"
)

func dialNoTLS(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

func getConfigForTLS() (*tls.Config, error) {
	trusted, err := common.LoadCertPool(serverRootCAPath)
	if err != nil {
		return nil, fmt.Errorf("could not load trusted bundle: %w", err)
	}

	return &tls.Config{
		RootCAs: trusted,
	}, nil
}

func dialTLS(network, address string) (net.Conn, error) {
	tlsConf, err := getConfigForTLS()
	if err != nil {
		return nil, err
	}

	return tls.Dial(network, address, tlsConf)
}

func dialMTLS(network, address string) (net.Conn, error) {
	tlsConf, err := getConfigForTLS()
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not load client TLS credentials: %w", err)
	}

	// Enable mTLS
	tlsConf.Certificates = []tls.Certificate{cert}

	return tls.Dial(network, address, tlsConf)
}
