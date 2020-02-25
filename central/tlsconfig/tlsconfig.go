package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/fileutils"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/mtls/verifier"
)

const (
	// TLSCertFileName is the tls certificate filename.
	TLSCertFileName = `tls.crt`
	// TLSKeyFileName is the private key filename.
	TLSKeyFileName = `tls.key`

	// DefaultCertPath is the path where the default TLS cert is located.
	DefaultCertPath = `/run/secrets/stackrox.io/default-tls-cert`
)

// NewCentralTLSConfigurer returns a new tls configurer to be used for Central.
func NewCentralTLSConfigurer() verifier.TLSConfigurer {
	return verifier.TLSConfigurerFunc(createTLSConfig)
}

func loadDefaultCertificate(dir string) (*tls.Certificate, error) {
	certFile := filepath.Join(dir, TLSCertFileName)
	keyFile := filepath.Join(dir, TLSKeyFileName)

	if filesExist, err := fileutils.AllExist(certFile, keyFile); err != nil || !filesExist {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, errors.Wrap(err, "parsing leaf certificate")
	}

	return &cert, nil
}

func loadInternalCertificateFromFiles() (*tls.Certificate, error) {
	if filesExist, err := fileutils.AllExist(mtls.CertFilePath(), mtls.KeyFilePath()); err != nil || !filesExist {
		return nil, err
	}

	cert, err := mtls.LeafCertificateFromFile()
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func issueInternalCertificate() (*tls.Certificate, error) {
	issuedCert, err := mtls.IssueNewCert(mtls.CentralSubject)
	if err != nil {
		return nil, errors.Wrap(err, "server keypair")
	}
	caPEM, err := mtls.CACertPEM()
	if err != nil {
		return nil, errors.Wrap(err, "CA cert retrieval")
	}
	serverCertBundle := append(issuedCert.CertPEM, caPEM...)

	serverTLSCert, err := tls.X509KeyPair(serverCertBundle, issuedCert.KeyPEM)
	if err != nil {
		return nil, errors.Wrap(err, "tls conversion")
	}
	return &serverTLSCert, nil
}

func getInternalCertificate() (*tls.Certificate, error) {
	// First try to load the internal certificate from files. If the files don't exist, issue
	// ourselves a cert.
	if certFromFiles, err := loadInternalCertificateFromFiles(); err != nil {
		return nil, err
	} else if certFromFiles != nil {
		return certFromFiles, nil
	}

	return issueInternalCertificate()
}

func serverCerts() ([]tls.Certificate, error) {
	var certs []tls.Certificate

	defaultCert, err := loadDefaultCertificate(DefaultCertPath)
	if err != nil {
		return nil, errors.Wrap(err, "loading default certificate")
	}
	if defaultCert != nil {
		certs = append(certs, *defaultCert)
	}

	internalCert, err := getInternalCertificate()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving internal certificate")
	} else if internalCert == nil {
		return nil, errors.New("no internal cert available")
	}
	certs = append(certs, *internalCert)
	return certs, nil

}

func createTLSConfig() (*tls.Config, error) {
	certPool, err := verifier.TrustedCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "loading trusted cert pool")
	}

	certs, err := serverCerts()
	if err != nil {
		return nil, err
	}

	cfg := verifier.DefaultTLSServerConfig(certPool, certs)

	return cfg, nil
}
