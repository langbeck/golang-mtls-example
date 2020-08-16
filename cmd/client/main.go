package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var dialers = map[string]func(string, string) (net.Conn, error){
	"notls": dialNoTLS,
	"tls":   dialTLS,
	"mtls":  dialMTLS,
}

func main() {
	optMode := flag.String("mode", "notls", "Select mode of operation: notls, tls, mtls")
	flag.Parse()

	dial, found := dialers[*optMode]
	if !found {
		logrus.Errorf("Invalid mode %q", *optMode)
		return
	}

	c, err := dial("tcp", "localhost:9900")
	if err != nil {
		logrus.Errorf("Failed to connect to the server: %v", err)
		return
	}

	defer c.Close()

	go io.Copy(os.Stderr, c)

	for i := 0; ; i++ {
		_, err = fmt.Fprintln(c, i)
		if err != nil {
			logrus.Errorf("I/O error: %v", err)
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}
