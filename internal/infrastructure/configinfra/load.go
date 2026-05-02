package configinfra

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Load(config *Config) error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	v := viper.New()
	setDefaults(v)

	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	v.SetConfigName(fmt.Sprintf("env.%s", env))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to find config file: %w", err)
		}

		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	fmt.Printf("config loaded [env=%s, file=%s]\n", env, v.ConfigFileUsed())
	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.env", "production")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.service_name", "winsku")
	v.SetDefault("server.tls.enable", false)

	v.SetDefault("database.port", "5432")
	v.SetDefault("database.ssl_mode", "require")
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.max_idle_conns", 25)
	v.SetDefault("database.conn_max_lifetime", "30m")

	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt_auth.access_token_exp", time.Minute*15)
	v.SetDefault("jwt_auth.refresh_token_exp", time.Hour*168)
	v.SetDefault("jwt_auth.remember_me_token_exp", time.Hour*720)

	v.SetDefault("rate_limit.enable", true)
	v.SetDefault("rate_limit.limit", 60)
	v.SetDefault("rate_limit.period", time.Minute)

	v.SetDefault("cors.enable", true)
	v.SetDefault("cors.allow_credentials", false)
}
