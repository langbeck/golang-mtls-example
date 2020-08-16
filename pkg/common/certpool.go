package common

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func LoadCertPool(certPaths ...string) (*x509.CertPool, error) {
	trusted := x509.NewCertPool()

	for _, certPath := range certPaths {
		caCertBytes, err := ioutil.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("could not read certificate at %q: %w", certPath, err)
		}

		wasAdded := trusted.AppendCertsFromPEM(caCertBytes)
		if !wasAdded {
			return nil, fmt.Errorf("could not load certificate at %q: %w", certPath, err)
		}
	}

	return trusted, nil
}
