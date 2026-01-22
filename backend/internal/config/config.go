package config

import "github.com/spf13/viper"

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
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // 自动读取环境变量

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.SetDefault("DB_DSN", "postgres://postgres:password@localhost:5432/postgres")
	viper.SetDefault("MATCHING_ENABLED", true)
	viper.SetDefault("MATCHING_LEASE_TTL", "20s")
	viper.SetDefault("MATCHING_LEASE_REFRESH", "5s")
	viper.SetDefault("MATCHING_SNAPSHOT_EVERY", "2s")
	viper.SetDefault("MATCHING_OWNER", "")

	err = viper.ReadInConfig()
	// If config file not found, we can still use env vars
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = nil
	}

	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&config)
	return
}
