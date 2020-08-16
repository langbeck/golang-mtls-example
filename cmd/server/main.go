package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
)

var listeners = map[string]func(string, string) (net.Listener, error){
	"notls":       listenNoTLS,
	"tls":         listenTLS,
	"mtls":        listenMTLS,
	"mtls-verify": listenVerifiedMTLS,
}

func main() {
	optMode := flag.String("mode", "notls", "Select mode of operation: notls, tls, mtls, mtls-verify")
	flag.Parse()

	listen, found := listeners[*optMode]
	if !found {
		logrus.Errorf("Invalid mode %q", *optMode)
		return
	}

	l, err := listen("tcp", ":9900")
	if err != nil {
		logrus.WithError(err).Error("could not start the server")
		return
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			logrus.WithError(err).Error("could not accept a connection")
			return
		}

		go func(c net.Conn) {
			defer c.Close()

			log := logrus.WithField("remote", c.RemoteAddr())

			// Check if it's a TLS connection and it's properties.
			// This isn't necessary. Just showing how an TCP application
			// (including HTTP servers and clients) can try to manually
			// check if they are in a secure connection and get any
			// certificate that the client may have presented.
			tlsConn, ok := c.(*tls.Conn)
			if ok {
				log.Infof("TLS connection detected")

				// Ensure TLS handshake is complete so we can access
				// client certificate (if any).
				err = tlsConn.Handshake()
				if err != nil {
					log.WithError(err).Error("TLS handshake failed")
					return
				}

				state := tlsConn.ConnectionState()
				if len(state.PeerCertificates) == 0 {
					log.Info("Cliend did not provided any certificate")
				} else {
					for i, cert := range state.PeerCertificates {
						log.Infof("  [#%d] Subject: %s", i+1, cert.Subject)
					}
				}
			}

			// Just a dummy server that sends even numbers back to the client.
			// Nothing interesting here.
			r := bufio.NewReader(c)
			for {
				line, _, err := r.ReadLine()
				if err != nil {
					if err == io.EOF {
						break
					}

					log.WithError(err).Error("I/O error")
					return
				}

				n, err := strconv.Atoi(string(line))
				if err != nil {
					log.WithError(err).Error("parse error")
					return
				}

				if n%2 == 0 {
					fmt.Fprintf(c, "Don't want %d\n", n)
				} else {
					log.Infof("Keeping %d\n", n)
				}
			}
		}(c)
	}
}
