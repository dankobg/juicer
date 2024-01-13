package config

import (
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// App contains common app settings
type App struct {
	ENV             string `koanf:"env"`
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	BaseURL         string `koanf:"base_url"`
	WebsiteURL      string `koanf:"website_url"`
	KratosPublicURL string `koanf:"kratos_public_url"`
	KratosAdminURL  string `koanf:"kratos_admin_url"`
}

// ServerConfig contains the http server settings
type ServerConfig struct {
	ReadHeaderTimeout time.Duration `koanf:"read_header_timeout"`
	ReadTimeout       time.Duration `koanf:"read_timeout"`
	WriteTimeout      time.Duration `koanf:"write_timeout"`
	IdleTimeout       time.Duration `koanf:"idle_timeout"`
	UseTLS            bool          `koanf:"use_tls"`
	CERT_FILE         string        `koanf:"cert_file"`
	KEY_FILE          string        `koanf:"key_file"`
}

// CorsConfig contains the CORS settings
type CorsConfig struct {
	AllowOrigins     []string `koanf:"allow_origins"`
	AllowMethods     []string `koanf:"allow_methods"`
	AllowHeaders     []string `koanf:"allow_headers"`
	ExposeHeaders    []string `koanf:"expose_headers"`
	AllowCredentials bool     `koanf:"allow_credentials"`
	MaxAge           int      `koanf:"max_age"`
	Debug            bool     `koanf:"debug"`
}

// DatabaseConfig contains DB settings
type DatabaseConfig struct {
	Host         string        `koanf:"host"`
	Port         int           `koanf:"port"`
	DB           string        `koanf:"db"`
	User         string        `koanf:"user"`
	Password     string        `koanf:"password"`
	Scheme       string        `koanf:"scheme"`
	SSLMode      string        `koanf:"ssl_mode"`
	RetriesNum   int           `koanf:"retries_num"`
	RetriesDelay time.Duration `koanf:"retries_delay"`
}

// RedisConfig contains Redis db settings
type RedisConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	DB       int    `koanf:"db"`
	Password string `koanf:"password"`
}

// EmailConfig contains email settings
type EmailConfig struct {
	Enabled         bool   `koanf:"enabled"`
	TLS             bool   `koanf:"tls"`
	FromName        string `koanf:"from_name"`
	FromAddress     string `koanf:"from_address"`
	SMTPHost        string `koanf:"smtp_host"`
	SMTPPort        int    `koanf:"smtp_port"`
	SMTPUsername    string `koanf:"smtp_username"`
	SMTPPassword    string `koanf:"smtp_password"`
	DevSMTPHost     string `koanf:"dev_smtp_host"`
	DevSMTPPort     int    `koanf:"dev_smtp_port"`
	DevSMTPUsername string `koanf:"dev_smtp_username"`
	DevSMTPPassword string `koanf:"dev_smtp_password"`
}

// LoggerConfig contains the logger settings
type LoggerConfig struct {
	PrettyConsole bool   `koanf:"pretty_console"`
	Level         string `koanf:"level"`
}

// Config represents the app config
type Config struct {
	App      `koanf:",squash"`
	Server   ServerConfig   `koanf:"server"`
	Cors     CorsConfig     `koanf:"cors"`
	Database DatabaseConfig `koanf:"postgres"`
	Redis    RedisConfig    `koanf:"redis"`
	Email    EmailConfig    `koanf:"email"`
	Logger   LoggerConfig   `koanf:"logger"`
}

// loadEnv loads env files by convention: https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use
func loadEnv(k *koanf.Koanf) error {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if err := godotenv.Load(".env." + env + ".local"); err != nil {
		log.Printf("file: %q not present, skipping", ".env."+env+".local")
	}

	if env != "test" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Printf("file: %q not present, skipping", ".env.local")
		}
	}

	if err := godotenv.Load(".env." + env); err != nil {
		log.Printf("file: %q not present, skipping", ".env."+env)
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("file: %q not present, skipping", ".env")
	}

	return nil
}

func loadFromEnv(k *koanf.Koanf) error {
	p := env.ProviderWithValue("", ".", func(s string, v string) (string, interface{}) {
		prefixes := []string{"server", "cors", "postgres", "redis", "email", "logger"}

		str := strings.ToLower(strings.TrimPrefix(s, "JUICER_"))

		key := str
		if slices.Contains(prefixes, str) {
			key = strings.Replace(str, "_", ".", 1)
		}

		if strings.Contains(v, ",") {
			return key, strings.Split(v, ",")
		}

		return key, v
	})

	return k.Load(p, nil)
}

func loadDefaults(k *koanf.Koanf) error {
	return k.Load(file.Provider("config/defaults.yaml"), yaml.Parser())
}

func getConfig(k *koanf.Koanf) (*Config, error) {
	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// New creates the new config
func New() (*Config, *koanf.Koanf, error) {
	k := koanf.New(".")

	if err := loadEnv(k); err != nil {
		return nil, nil, err
	}

	if err := loadDefaults(k); err != nil {
		return nil, nil, err
	}

	if err := loadFromEnv(k); err != nil {
		return nil, nil, err
	}

	cfg, err := getConfig(k)
	if err != nil {
		return nil, nil, err
	}

	return cfg, k, nil
}
