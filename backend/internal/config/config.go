package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var (
	ErrMissingDBDSN       = errors.New("database DSN is required")
	ErrInvalidPort        = errors.New("invalid port number")
	ErrInvalidDuration    = errors.New("invalid duration format")
	ErrMissingExchangeURL = errors.New("exchange URL is required for trading features")
)

type Config struct {
	ExchangeURL           string `mapstructure:"EXCHANGE_URL"`
	DB_DSN                string `mapstructure:"DB_DSN"`
	NatsURL               string `mapstructure:"NATS_URL"`
	Port                  string `mapstructure:"PORT"`
	MatchingEnabled       bool   `mapstructure:"MATCHING_ENABLED"`
	MatchingLeaseTTL      string `mapstructure:"MATCHING_LEASE_TTL"`
	MatchingLeaseRefresh  string `mapstructure:"MATCHING_LEASE_REFRESH"`
	MatchingSnapshotEvery string `mapstructure:"MATCHING_SNAPSHOT_EVERY"`
	MatchingOwner         string `mapstructure:"MATCHING_OWNER"`

	// 新增配置项
	LogLevel            string `mapstructure:"LOG_LEVEL"`
	EnableMetrics       bool   `mapstructure:"ENABLE_METRICS"`
	MetricsPort         string `mapstructure:"METRICS_PORT"`
	StripeAPIKey        string `mapstructure:"STRIPE_API_KEY"`
	StripeWebhookSecret string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	// 验证必需的数据库配置
	if c.DB_DSN == "" {
		return ErrMissingDBDSN
	}

	// 验证端口格式
	if c.Port == "" {
		c.Port = "8080"
	}

	// 验证时间间隔格式
	if c.MatchingEnabled {
		if _, err := time.ParseDuration(c.MatchingLeaseTTL); err != nil {
			return fmt.Errorf("%w: MATCHING_LEASE_TTL=%s", ErrInvalidDuration, c.MatchingLeaseTTL)
		}
		if _, err := time.ParseDuration(c.MatchingLeaseRefresh); err != nil {
			return fmt.Errorf("%w: MATCHING_LEASE_REFRESH=%s", ErrInvalidDuration, c.MatchingLeaseRefresh)
		}
		if _, err := time.ParseDuration(c.MatchingSnapshotEvery); err != nil {
			return fmt.Errorf("%w: MATCHING_SNAPSHOT_EVERY=%s", ErrInvalidDuration, c.MatchingSnapshotEvery)
		}
	}

	// 验证日志级别
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be debug, info, warn, or error)", c.LogLevel)
	}

	return nil
}

// GetMatchingLeaseTTL 返回解析后的 LeaseTTL 时间
func (c *Config) GetMatchingLeaseTTL() time.Duration {
	d, _ := time.ParseDuration(c.MatchingLeaseTTL)
	if d == 0 {
		return 20 * time.Second
	}
	return d
}

// GetMatchingLeaseRefresh 返回解析后的 LeaseRefresh 时间
func (c *Config) GetMatchingLeaseRefresh() time.Duration {
	d, _ := time.ParseDuration(c.MatchingLeaseRefresh)
	if d == 0 {
		return 5 * time.Second
	}
	return d
}

// GetMatchingSnapshotEvery 返回解析后的 SnapshotEvery 时间
func (c *Config) GetMatchingSnapshotEvery() time.Duration {
	d, _ := time.ParseDuration(c.MatchingSnapshotEvery)
	if d == 0 {
		return 2 * time.Second
	}
	return d
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // 自动读取环境变量

	// 设置默认值
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.SetDefault("DB_DSN", "postgres://postgres:password@localhost:5432/postgres")
	viper.SetDefault("MATCHING_ENABLED", true)
	viper.SetDefault("MATCHING_LEASE_TTL", "20s")
	viper.SetDefault("MATCHING_LEASE_REFRESH", "5s")
	viper.SetDefault("MATCHING_SNAPSHOT_EVERY", "2s")
	viper.SetDefault("MATCHING_OWNER", "")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENABLE_METRICS", true)
	viper.SetDefault("METRICS_PORT", "9090")

	err = viper.ReadInConfig()
	// If config file not found, we can still use env vars
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = nil
	}

	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}
