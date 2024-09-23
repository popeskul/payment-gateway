package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Auth          AuthConfig
	AcquiringBank AcquiringBankConfig
	Logging       LoggingConfig
	Metrics       MetricsConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConnections  int
	MinConnections  int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type AuthConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

type AcquiringBankConfig struct {
	ProcessingDelay time.Duration
	FailureRate     float64
}

type LoggingConfig struct {
	Level  string
	Format string
}

type MetricsConfig struct {
	Enabled bool
	Port    string
}

func (d DatabaseConfig) GetHost() string {
	return d.Host
}

func (d DatabaseConfig) GetPort() string {
	return d.Port
}

func (d DatabaseConfig) GetUser() string {
	return d.User
}

func (d DatabaseConfig) GetPassword() string {
	return d.Password
}

func (d DatabaseConfig) GetDBName() string {
	return d.DBName
}

func (d DatabaseConfig) GetMaxConnections() int {
	return d.MaxConnections
}

func (d DatabaseConfig) GetMinConnections() int {
	return d.MinConnections
}

func (d DatabaseConfig) GetMaxConnLifetime() int {
	return int(d.MaxConnLifetime)
}

func (d DatabaseConfig) GetMaxConnIdleTime() int {
	return int(d.MaxConnIdleTime)
}

func (a AuthConfig) GetAccessTokenSecret() string {
	return a.AccessTokenSecret
}

func (a AuthConfig) GetRefreshTokenSecret() string {
	return a.RefreshTokenSecret
}

func (a AuthConfig) GetAccessTokenTTL() int {
	return int(a.AccessTokenTTL)
}

func (a AuthConfig) GetRefreshTokenTTL() int {
	return int(a.RefreshTokenTTL)
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	config.Auth.AccessTokenSecret = viper.GetString("ACCESS_TOKEN_SECRET")
	config.Auth.RefreshTokenSecret = viper.GetString("REFRESH_TOKEN_SECRET")
	config.Database.Password = viper.GetString("DB_PASSWORD")

	if metricsEnabled := viper.GetString("METRICS_ENABLED"); metricsEnabled != "" {
		enabled, err := strconv.ParseBool(metricsEnabled)
		if err != nil {
			return nil, fmt.Errorf("invalid METRICS_ENABLED value: %v", err)
		}
		config.Metrics.Enabled = enabled
	}

	if dbHost := viper.GetString("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := viper.GetString("DB_PORT"); dbPort != "" {
		config.Database.Port = dbPort
	}
	if dbUser := viper.GetString("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbName := viper.GetString("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}
	if dbSSLMode := viper.GetString("DB_SSLMODE"); dbSSLMode != "" {
		config.Database.SSLMode = dbSSLMode
	}
	if logLevel := viper.GetString("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}
	if logFormat := viper.GetString("LOG_FORMAT"); logFormat != "" {
		config.Logging.Format = logFormat
	}
	if metricsPort := viper.GetString("METRICS_PORT"); metricsPort != "" {
		config.Metrics.Port = metricsPort
	}

	setDefaults(&config)

	return &config, nil
}

func setDefaults(config *Config) {
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 10 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 10 * time.Second
	}
	if config.Server.IdleTimeout == 0 {
		config.Server.IdleTimeout = 120 * time.Second
	}
	if config.Database.MaxConnections == 0 {
		config.Database.MaxConnections = 10
	}
	if config.Database.MinConnections == 0 {
		config.Database.MinConnections = 2
	}
	if config.Database.MaxConnLifetime == 0 {
		config.Database.MaxConnLifetime = 1 * time.Hour
	}
	if config.Database.MaxConnIdleTime == 0 {
		config.Database.MaxConnIdleTime = 30 * time.Minute
	}
	if config.Auth.AccessTokenTTL == 0 {
		config.Auth.AccessTokenTTL = 15 * time.Minute
	}
	if config.Auth.RefreshTokenTTL == 0 {
		config.Auth.RefreshTokenTTL = 7 * 24 * time.Hour // 7 days
	}
	if config.AcquiringBank.ProcessingDelay == 0 {
		config.AcquiringBank.ProcessingDelay = 200 * time.Millisecond
	}
	if config.AcquiringBank.FailureRate == 0 {
		config.AcquiringBank.FailureRate = 0.05
	}
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
	if !config.Metrics.Enabled {
		config.Metrics.Enabled = true
	}
	if config.Metrics.Port == "" {
		config.Metrics.Port = "9090"
	}
}
