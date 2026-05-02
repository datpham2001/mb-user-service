package configinfra

import (
	"time"
)

type Config struct {
	Server    ServerConfig    `yaml:"server" mapstructure:"server"`
	Database  DatabaseConfig  `yaml:"database" mapstructure:"database"`
	Redis     RedisConfig     `yaml:"redis" mapstructure:"redis"`
	JwtAuth   JwtAuthConfig   `yaml:"jwt_auth" mapstructure:"jwt_auth"`
	OAuth2    OAuth2Config    `yaml:"oauth2" mapstructure:"oauth2"`
	RateLimit RateLimitConfig `yaml:"rate_limit" mapstructure:"rate_limit"`
	CORS      CORSConfig      `yaml:"cors" mapstructure:"cors"`
	S3        S3Config        `yaml:"s3" mapstructure:"s3"`
}

type ServerConfig struct {
	Env         string `yaml:"env" mapstructure:"env"`
	Host        string `yaml:"host" mapstructure:"host"`
	Port        string `yaml:"port" mapstructure:"port"`
	ServiceName string `yaml:"service_name" mapstructure:"service_name"`
	TLS         struct {
		Enable   bool   `yaml:"enable" mapstructure:"enable"`
		CertFile string `yaml:"cert_file" mapstructure:"cert_file"`
		KeyFile  string `yaml:"key_file" mapstructure:"key_file"`
	} `yaml:"tls" mapstructure:"tls"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host" mapstructure:"host"`
	Port            string        `yaml:"port" mapstructure:"port"`
	User            string        `yaml:"user" mapstructure:"user"`
	Password        string        `yaml:"password" mapstructure:"password"`
	DBName          string        `yaml:"db_name" mapstructure:"db_name"`
	SSLMode         string        `yaml:"ssl_mode" mapstructure:"ssl_mode"`
	MaxOpenConns    int           `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

type JwtAuthConfig struct {
	SecretKey          string        `yaml:"secret_key" mapstructure:"secret_key"`
	AccessTokenExp     time.Duration `yaml:"access_token_exp" mapstructure:"access_token_exp"`
	RefreshTokenExp    time.Duration `yaml:"refresh_token_exp" mapstructure:"refresh_token_exp"`
	RememberMeTokenExp time.Duration `yaml:"remember_me_token_exp" mapstructure:"remember_me_token_exp"`
}

type RateLimitConfig struct {
	Limit  int           `yaml:"limit" mapstructure:"limit"`
	Period time.Duration `yaml:"period" mapstructure:"period"`
	Enable bool          `yaml:"enable" mapstructure:"enable"`
}

type CORSConfig struct {
	Enable           bool     `yaml:"enable" mapstructure:"enable"`
	AllowedOrigins   []string `yaml:"allowed_origins" mapstructure:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods" mapstructure:"allowed_methods"`
	AllowCredentials bool     `yaml:"allow_credentials" mapstructure:"allow_credentials"`
}

type S3Config struct{}

type OAuth2Config struct {
	Google OAuth2ProviderConfig `yaml:"google" mapstructure:"google"`
}

type OAuth2ProviderConfig struct {
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
}
