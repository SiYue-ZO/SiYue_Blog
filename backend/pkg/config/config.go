package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	HTTP     HTTP     `mapstructure:"http"`
	Postgres Postgres `mapstructure:"postgres"`
	JWT      JWT      `mapstructure:"jwt"`
	CORS     CORS     `mapstructure:"cors"`
	Log      Log      `mapstructure:"log"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

type HTTP struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type Postgres struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DBName       string        `mapstructure:"dbname"`
	SSLMode      string        `mapstructure:"sslmode"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
}

func (p Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode,
	)
}

type JWT struct {
	AccessSecret  string        `mapstructure:"access_secret"`
	RefreshSecret string        `mapstructure:"refresh_secret"`
	AccessTTL     time.Duration `mapstructure:"access_ttl"`
	RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
}

type CORS struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("/etc/siyue-blog")

	v.SetEnvPrefix("SIYUE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Defaults
	v.SetDefault("app.name", "siyue-blog")
	v.SetDefault("app.version", "0.1.0")
	v.SetDefault("app.env", "development")
	v.SetDefault("http.addr", ":8080")
	v.SetDefault("http.read_timeout", "10s")
	v.SetDefault("http.write_timeout", "10s")
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "siyue")
	v.SetDefault("postgres.password", "siyue_dev")
	v.SetDefault("postgres.dbname", "siyue_blog")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.max_open_conns", 25)
	v.SetDefault("postgres.max_idle_conns", 10)
	v.SetDefault("postgres.max_lifetime", "300s")
	v.SetDefault("jwt.access_ttl", "15m")
	v.SetDefault("jwt.refresh_ttl", "168h")
	v.SetDefault("cors.allow_origins", []string{"http://localhost:3000"})
	v.SetDefault("log.level", "debug")
	v.SetDefault("log.format", "console")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
