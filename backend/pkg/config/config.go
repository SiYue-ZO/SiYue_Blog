package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	HTTP     HTTP     `mapstructure:"http"`
	Postgres Postgres `mapstructure:"postgres"`
	JWT      JWT      `mapstructure:"jwt"`
	Auth     Auth     `mapstructure:"auth"`
	Admin    Admin    `mapstructure:"admin"`
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
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	MaxBodyBytes int64         `mapstructure:"max_body_bytes"`
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
	connectionURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Host:   net.JoinHostPort(p.Host, fmt.Sprintf("%d", p.Port)),
		Path:   p.DBName,
	}
	query := connectionURL.Query()
	query.Set("sslmode", p.SSLMode)
	connectionURL.RawQuery = query.Encode()
	return connectionURL.String()
}

type JWT struct {
	AccessSecret  string        `mapstructure:"access_secret"`
	RefreshSecret string        `mapstructure:"refresh_secret"`
	AccessTTL     time.Duration `mapstructure:"access_ttl"`
	RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
}

type Auth struct {
	AllowRegistration bool   `mapstructure:"allow_registration"`
	CookieDomain      string `mapstructure:"cookie_domain"`
	CookieSecure      bool   `mapstructure:"cookie_secure"`
}

type Admin struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
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
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.env", "development")
	v.SetDefault("http.addr", ":8080")
	v.SetDefault("http.read_timeout", "10s")
	v.SetDefault("http.write_timeout", "10s")
	v.SetDefault("http.idle_timeout", "60s")
	v.SetDefault("http.max_body_bytes", 2<<20)
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "siyue")
	v.SetDefault("postgres.password", "siyue_dev")
	v.SetDefault("postgres.dbname", "siyue_blog")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.max_open_conns", 25)
	v.SetDefault("postgres.max_idle_conns", 10)
	v.SetDefault("postgres.max_lifetime", "300s")
	v.SetDefault("jwt.access_secret", "")
	v.SetDefault("jwt.refresh_secret", "")
	v.SetDefault("jwt.access_ttl", "15m")
	v.SetDefault("jwt.refresh_ttl", "168h")
	v.SetDefault("auth.allow_registration", false)
	v.SetDefault("auth.cookie_secure", false)
	v.SetDefault("admin.username", "admin")
	v.SetDefault("admin.email", "")
	v.SetDefault("admin.password", "")
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
	if origins := strings.TrimSpace(os.Getenv("SIYUE_CORS_ALLOW_ORIGINS")); origins != "" {
		cfg.CORS.AllowOrigins = strings.Split(origins, ",")
		for i := range cfg.CORS.AllowOrigins {
			cfg.CORS.AllowOrigins[i] = strings.TrimSpace(cfg.CORS.AllowOrigins[i])
		}
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.App.Env != "development" && c.App.Env != "staging" && c.App.Env != "production" {
		return fmt.Errorf("app.env must be development, staging, or production")
	}
	if c.HTTP.ReadTimeout <= 0 || c.HTTP.WriteTimeout <= 0 || c.HTTP.IdleTimeout <= 0 {
		return fmt.Errorf("http timeouts must be positive")
	}
	if c.HTTP.MaxBodyBytes < 1024 {
		return fmt.Errorf("http.max_body_bytes must be at least 1024")
	}
	if c.Postgres.Host == "" || c.Postgres.User == "" || c.Postgres.DBName == "" {
		return fmt.Errorf("postgres host, user, and dbname are required")
	}
	if c.Postgres.Port < 1 || c.Postgres.Port > 65535 {
		return fmt.Errorf("postgres port must be between 1 and 65535")
	}
	if c.Postgres.MaxOpenConns < 1 || c.Postgres.MaxIdleConns < 0 || c.Postgres.MaxIdleConns > c.Postgres.MaxOpenConns || c.Postgres.MaxLifetime <= 0 {
		return fmt.Errorf("invalid postgres connection pool settings")
	}
	switch c.Postgres.SSLMode {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
	default:
		return fmt.Errorf("invalid postgres sslmode")
	}
	if len(c.JWT.AccessSecret) < 32 || len(c.JWT.RefreshSecret) < 32 {
		return fmt.Errorf("JWT access and refresh secrets must each contain at least 32 characters")
	}
	if c.JWT.AccessSecret == c.JWT.RefreshSecret {
		return fmt.Errorf("JWT access and refresh secrets must be different")
	}
	if c.JWT.AccessTTL <= 0 || c.JWT.RefreshTTL <= c.JWT.AccessTTL {
		return fmt.Errorf("JWT refresh TTL must be greater than the positive access TTL")
	}
	for _, origin := range c.CORS.AllowOrigins {
		if origin == "*" {
			return fmt.Errorf("wildcard CORS origins are not allowed with credentialed requests")
		}
		parsed, err := url.ParseRequestURI(origin)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("invalid CORS origin %q", origin)
		}
	}
	if c.Admin.Password != "" {
		if len(c.Admin.Password) < 12 {
			return fmt.Errorf("admin bootstrap password must contain at least 12 characters")
		}
		if len([]byte(c.Admin.Password)) > 72 {
			return fmt.Errorf("admin bootstrap password must not exceed 72 bytes")
		}
		if c.Admin.Email == "" || !strings.Contains(c.Admin.Email, "@") {
			return fmt.Errorf("a valid admin bootstrap email is required when a password is configured")
		}
	}
	if len(c.Admin.Username) > 50 || len(c.Admin.Email) > 255 {
		return fmt.Errorf("admin username or email exceeds the database limit")
	}
	if c.App.Env == "production" {
		if !c.Auth.CookieSecure {
			return fmt.Errorf("auth.cookie_secure must be true in production")
		}
		if c.Admin.Email == "" || c.Admin.Password == "" {
			return fmt.Errorf("admin bootstrap email and password are required in production")
		}
		if len(c.Postgres.Password) < 16 {
			return fmt.Errorf("postgres password must contain at least 16 characters in production")
		}
		for label, value := range map[string]string{
			"JWT access secret":  c.JWT.AccessSecret,
			"JWT refresh secret": c.JWT.RefreshSecret,
			"admin password":     c.Admin.Password,
			"postgres password":  c.Postgres.Password,
		} {
			lower := strings.ToLower(value)
			if strings.Contains(lower, "replace-with") || strings.Contains(lower, "change-me") || strings.Contains(lower, "not-used") {
				return fmt.Errorf("%s still contains a placeholder value", label)
			}
		}
		adminEmail := strings.ToLower(c.Admin.Email)
		if strings.HasSuffix(adminEmail, "@example.com") || strings.HasSuffix(adminEmail, ".local") {
			return fmt.Errorf("admin email must not use a reserved placeholder domain in production")
		}
	}
	return nil
}
