package commands

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Peltoche/zapette/assets"
	"github.com/Peltoche/zapette/internal/server"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/logger"
	"github.com/Peltoche/zapette/internal/tools/response"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrConflictTLSConfig = errors.New("can't use --self-signed-cert and --tls-key at the same time")
	ErrDevFlagRequire    = errors.New("this flag require the --dev flag setup")
)

type Config struct {
	LogLevel       string   `mapstructure:"log-level"`
	Folder         string   `mapstructure:"folder"`
	TLSCert        string   `mapstructure:"tls-cert"`
	TLSKey         string   `mapstructure:"tls-key"`
	HTTPHost       string   `mapstructure:"http-host"`
	HTTPHostnames  []string `mapstructure:"http-hosts"`
	HTTPPort       int      `mapstructure:"http-port"`
	MemoryFS       bool     `mapstructure:"memory-fs"`
	SelfSignedCert bool     `mapstructure:"self-signed-cert"`
	Debug          bool     `mapstructure:"debug"`
	Dev            bool     `mapstructure:"dev"`
	HotReload      bool     `mapstructure:"hot-reload"`
}

func NewConfigFromCmd(cmd *cobra.Command) (server.Config, error) {
	var cfg Config

	viper.AutomaticEnv()
	viper.SetEnvPrefix("zapette")

	viper.BindPFlags(cmd.Flags())

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return server.Config{}, fmt.Errorf("config error: %w", err)
	}

	if cfg.HotReload && !cfg.Dev {
		return server.Config{}, fmt.Errorf("--hot-reload: %w", ErrDevFlagRequire)
	}

	if cfg.MemoryFS && !cfg.Dev {
		return server.Config{}, fmt.Errorf("--memory-fs: %w", ErrDevFlagRequire)
	}

	var logLevel slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "err", "error":
		logLevel = slog.LevelError
	default:
		return server.Config{}, fmt.Errorf("invalid log level")
	}

	if cfg.Debug {
		logLevel = slog.LevelDebug
	}

	var fs afero.Fs
	var storagePath string
	if cfg.MemoryFS {
		fs = afero.NewMemMapFs()
		storagePath = ":memory:"
	} else {
		fs = afero.NewOsFs()
		storagePath = path.Join(cfg.Folder, "db.sqlite")
	}

	err = fs.MkdirAll(cfg.Folder, 0o755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return server.Config{}, fmt.Errorf("failed to create %q: %w", cfg.Folder, err)
	}

	if cfg.SelfSignedCert {
		if cfg.TLSCert != "" || cfg.TLSKey != "" {
			return server.Config{}, ErrConflictTLSConfig
		}

		cfg.TLSCert, cfg.TLSKey, err = generateSelfSignedCertificate(cfg.HTTPHostnames, cfg.Folder, fs)
		if err != nil {
			return server.Config{}, fmt.Errorf("failed to generate the self-signed certificate: %w", err)
		}
	}

	isTLSEnabled := cfg.TLSCert != "" || cfg.TLSKey != ""

	return server.Config{
		FS: fs,
		Listener: router.Config{
			Addr:      net.JoinHostPort(cfg.HTTPHost, strconv.Itoa(cfg.HTTPPort)),
			TLS:       isTLSEnabled,
			Secure:    !cfg.Dev,
			CertFile:  cfg.TLSCert,
			KeyFile:   cfg.TLSKey,
			HostNames: cfg.HTTPHostnames,
		},
		Storage: sqlstorage.Config{
			Path: storagePath,
		},
		Assets: assets.Config{
			HotReload: cfg.HotReload,
		},
		Tools: tools.Config{
			Response: response.Config{
				PrettyRender: cfg.Dev,
			},
			Log: logger.Config{
				Level:  logLevel,
				Output: os.Stderr,
			},
		},
		Folder: server.Folder(cfg.Folder),
		HTML: html.Config{
			PrettyRender: cfg.Dev,
			HotReload:    cfg.HotReload,
		},
	}, nil
}

func generateSelfSignedCertificate(hostnames []string, folderPath string, fs afero.Fs) (string, string, error) {
	sslfolder := path.Join(folderPath, "ssl")
	certificatePath := path.Join(sslfolder, "cert.pem")
	privateKeyPath := path.Join(sslfolder, "key.pem")

	err := fs.MkdirAll(sslfolder, 0o700)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the SSL folder: %w", err)
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to GenerateKey: %w", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Duck Corp"},
		},
		DNSNames:  hostnames,
		NotBefore: time.Now(),
		// NotAfter:  time.Now().Add(3 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %w", err)
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		return "", "", errors.New("failed to encode certificate to PEM")
	}

	if err := afero.WriteFile(fs, certificatePath, pemCert, 0o644); err != nil {
		return "", "", fmt.Errorf("failed to write the certificate into the data folder: %w", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("unable to marshal private key: %w", err)
	}

	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		return "", "", errors.New("failed to encode key to PEM")
	}

	if err := afero.WriteFile(fs, privateKeyPath, pemKey, 0o600); err != nil {
		return "", "", fmt.Errorf("failed to write the certificate into the data folder: %w", err)
	}

	return certificatePath, privateKeyPath, nil
}
