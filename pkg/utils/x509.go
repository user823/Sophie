package utils

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"os"
)

func LoadX509KeyPair(certFile, KeyFile string) tls.Certificate {
	certificate, err := tls.LoadX509KeyPair(certFile, KeyFile)
	if err != nil {
		panic(err.Error())
	}
	return certificate
}

func LoadCAs(CAPaths ...string) *x509.CertPool {
	caPool := x509.NewCertPool()
	for _, path := range CAPaths {
		caFile, err := os.Open(path)
		if err != nil {
			continue
		}
		certPEM, err := io.ReadAll(caFile)
		caFile.Close()
		if err != nil {
			continue
		}
		caPool.AppendCertsFromPEM(certPEM)
	}
	return caPool
}
