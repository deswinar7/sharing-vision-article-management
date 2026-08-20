package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/example/article-management/backend/internal/config"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Open(cfg config.Config) (*sqlx.DB, error) {
	tlsMode, err := configureTLS(cfg)
	if err != nil {
		return nil, err
	}

	driverConfig := mysqlDriver.Config{
		User:                 cfg.DBUser,
		Passwd:               cfg.DBPassword,
		Net:                  "tcp",
		Addr:                 cfg.DBHost + ":" + cfg.DBPort,
		DBName:               cfg.DBName,
		ParseTime:            true,
		AllowNativePasswords: true,
		TLSConfig:            tlsMode,
	}

	db, err := sqlx.Open("mysql", driverConfig.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func configureTLS(cfg config.Config) (string, error) {
	if cfg.DBCACertPath == "" {
		return cfg.DBTLSMode, nil
	}

	certificate, err := os.ReadFile(cfg.DBCACertPath)
	if err != nil {
		return "", fmt.Errorf("read database CA certificate: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(certificate) {
		return "", fmt.Errorf("database CA certificate is invalid")
	}
	const name = "custom-ca"
	if err := mysqlDriver.RegisterTLSConfig(name, &tls.Config{MinVersion: tls.VersionTLS12, RootCAs: pool}); err != nil {
		return "", fmt.Errorf("register database TLS config: %w", err)
	}
	return name, nil
}
