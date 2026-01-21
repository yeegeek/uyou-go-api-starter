// Package config 提供应用程序配置管理，支持多环境配置和环境变量覆盖
package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig        `mapstructure:"app" yaml:"app"`
	Database   DatabaseConfig   `mapstructure:"database" yaml:"database"`
	MongoDB    MongoDBConfig    `mapstructure:"mongodb" yaml:"mongodb"`
	Redis      RedisConfig      `mapstructure:"redis" yaml:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt" yaml:"jwt"`
	Server     ServerConfig     `mapstructure:"server" yaml:"server"`
	Logging    LoggingConfig    `mapstructure:"logging" yaml:"logging"`
	Ratelimit  RateLimitConfig  `mapstructure:"ratelimit" yaml:"ratelimit"`
	Migrations MigrationsConfig `mapstructure:"migrations" yaml:"migrations"`
	Health     HealthConfig     `mapstructure:"health" yaml:"health"`
	RabbitMQ   RabbitMQConfig   `mapstructure:"rabbitmq" yaml:"rabbitmq"`
	GRPC       GRPCConfig       `mapstructure:"grpc" yaml:"grpc"`
	Metrics    MetricsConfig    `mapstructure:"metrics" yaml:"metrics"`
	Scheduler  SchedulerConfig  `mapstructure:"scheduler" yaml:"scheduler"`
	Security   SecurityConfig   `mapstructure:"security" yaml:"security"`
}

// SchedulerConfig 定时任务配置
type SchedulerConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	Timezone string `mapstructure:"timezone" yaml:"timezone"`
}

type AppConfig struct {
	Name        string `mapstructure:"name" yaml:"name"`
	Version     string `mapstructure:"version" yaml:"version"`
	Environment string `mapstructure:"environment" yaml:"environment"`
	Debug       bool   `mapstructure:"debug" yaml:"debug"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host" yaml:"host"`
	Port            int    `mapstructure:"port" yaml:"port"`
	User            string `mapstructure:"user" yaml:"user"`
	Password        string `mapstructure:"password" yaml:"password"`
	Name            string `mapstructure:"name" yaml:"name"`
	SSLMode         string `mapstructure:"sslmode" yaml:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"` // 秒
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"` // 秒
}

type JWTConfig struct {
	Secret          string        `mapstructure:"secret" yaml:"secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl" yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl" yaml:"refresh_token_ttl"`
	TTLHours        int           `mapstructure:"ttlhours" yaml:"ttlhours"` // Deprecated: kept for backward compatibility
}

type ServerConfig struct {
	Port            string `mapstructure:"port" yaml:"port"`
	ReadTimeout     int    `mapstructure:"readtimeout" yaml:"readtimeout"`
	WriteTimeout    int    `mapstructure:"writetimeout" yaml:"writetimeout"`
	IdleTimeout     int    `mapstructure:"idletimeout" yaml:"idletimeout"`
	ShutdownTimeout int    `mapstructure:"shutdowntimeout" yaml:"shutdowntimeout"`
	MaxHeaderBytes  int    `mapstructure:"maxheaderbytes" yaml:"maxheaderbytes"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level" yaml:"level"`   // debug, info, warn, error
	Format string `mapstructure:"format" yaml:"format"` // json, text
	Output string `mapstructure:"output" yaml:"output"` // stdout, file
	File   string `mapstructure:"file" yaml:"file"`     // 日志文件路径
}

type RateLimitConfig struct {
	Enabled  bool          `mapstructure:"enabled" yaml:"enabled"`
	Requests int           `mapstructure:"requests" yaml:"requests"`
	Window   time.Duration `mapstructure:"window" yaml:"window"`
}

type MigrationsConfig struct {
	Directory   string `mapstructure:"directory" yaml:"directory"`
	Timeout     int    `mapstructure:"timeout" yaml:"timeout"`
	LockTimeout int    `mapstructure:"locktimeout" yaml:"locktimeout"`
}

type HealthConfig struct {
	Timeout              int  `mapstructure:"timeout" yaml:"timeout"`
	DatabaseCheckEnabled bool `mapstructure:"database_check_enabled" yaml:"database_check_enabled"`
}

// MongoDBConfig MongoDB 数据库配置
type MongoDBConfig struct {
	Enabled        bool   `mapstructure:"enabled" yaml:"enabled"`
	URI            string `mapstructure:"uri" yaml:"uri"`
	Database       string `mapstructure:"database" yaml:"database"`
	MaxPoolSize    int    `mapstructure:"max_pool_size" yaml:"max_pool_size"`
	MinPoolSize    int    `mapstructure:"min_pool_size" yaml:"min_pool_size"`
	ConnectTimeout int    `mapstructure:"connect_timeout" yaml:"connect_timeout"`
}

// RedisConfig Redis 缓存配置
type RedisConfig struct {
	Enabled      bool   `mapstructure:"enabled" yaml:"enabled"`
	Host         string `mapstructure:"host" yaml:"host"`
	Port         int    `mapstructure:"port" yaml:"port"`
	Password     string `mapstructure:"password" yaml:"password"`
	DB           int    `mapstructure:"db" yaml:"db"`
	DialTimeout  int    `mapstructure:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
	PoolSize     int    `mapstructure:"pool_size" yaml:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`
}

// RabbitMQConfig RabbitMQ 消息队列配置
// 注意：所有微服务连接到同一个 RabbitMQ 实例
// 发布和订阅行为由代码逻辑控制，不需要配置限制
type RabbitMQConfig struct {
	Enabled       bool   `mapstructure:"enabled" yaml:"enabled"`
	URL           string `mapstructure:"url" yaml:"url"`
	Exchange      string `mapstructure:"exchange" yaml:"exchange"`
	ExchangeType  string `mapstructure:"exchange_type" yaml:"exchange_type"`
	Queue         string `mapstructure:"queue" yaml:"queue"`
	RoutingKey    string `mapstructure:"routing_key" yaml:"routing_key"`
	PrefetchCount int    `mapstructure:"prefetch_count" yaml:"prefetch_count"`
}

// GRPCConfig gRPC 服务配置
type GRPCConfig struct {
	Enabled         bool   `mapstructure:"enabled" yaml:"enabled"`
	Port            string `mapstructure:"port" yaml:"port"`
	MaxRecvMsgSize  int    `mapstructure:"max_recv_msg_size" yaml:"max_recv_msg_size"`
	MaxSendMsgSize  int    `mapstructure:"max_send_msg_size" yaml:"max_send_msg_size"`
	ConnectionTimeout int  `mapstructure:"connection_timeout" yaml:"connection_timeout"`
}

// MetricsConfig Prometheus 监控配置
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
	Port    string `mapstructure:"port" yaml:"port"`
	Path    string `mapstructure:"path" yaml:"path"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	// Bcrypt 成本因子（推荐 10-14）
	BcryptCost int `mapstructure:"bcrypt_cost" yaml:"bcrypt_cost"`
	// 密码最小长度
	PasswordMinLength int `mapstructure:"password_min_length" yaml:"password_min_length"`
	// 密码需要大写字母
	PasswordRequireUppercase bool `mapstructure:"password_require_uppercase" yaml:"password_require_uppercase"`
	// 密码需要小写字母
	PasswordRequireLowercase bool `mapstructure:"password_require_lowercase" yaml:"password_require_lowercase"`
	// 密码需要数字
	PasswordRequireNumber bool `mapstructure:"password_require_number" yaml:"password_require_number"`
	// 密码需要特殊字符
	PasswordRequireSpecial bool `mapstructure:"password_require_special" yaml:"password_require_special"`
	// 最大登录尝试次数
	MaxLoginAttempts int `mapstructure:"max_login_attempts" yaml:"max_login_attempts"`
	// 账户锁定时间（分钟）
	LockoutDuration int `mapstructure:"lockout_duration" yaml:"lockout_duration"`
	// 启用安全响应头
	EnableSecurityHeaders bool `mapstructure:"enable_security_headers" yaml:"enable_security_headers"`
}

// LoadConfig loads configuration using Viper. If configPath is non-empty it
// will be used as the exact config file path, otherwise Viper searches common locations.
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	bindEnvVariables(v)

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		}
	} else {
		env := v.GetString("APP_ENVIRONMENT")
		if env == "" {
			env = "development"
		}

		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("configs")
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read base config file: %w", err)
			}
		}

		v.SetConfigName(fmt.Sprintf("config.%s", env))
		if err := v.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to merge environment config: %w", err)
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.App.Environment == "" {
		if e := v.GetString("app.environment"); e != "" {
			cfg.App.Environment = e
		} else if e := v.GetString("ENV"); e != "" {
			cfg.App.Environment = e
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func bindEnvVariables(v *viper.Viper) {
	envBindings := map[string]string{
		"app.name":                      "APP_NAME",
		"app.version":                   "APP_VERSION",
		"app.environment":               "APP_ENVIRONMENT",
		"app.debug":                     "APP_DEBUG",
		"database.host":                 "DATABASE_HOST",
		"database.port":                 "DATABASE_PORT",
		"database.user":                 "DATABASE_USER",
		"database.password":             "DATABASE_PASSWORD",
		"database.name":                 "DATABASE_NAME",
		"database.sslmode":              "DATABASE_SSLMODE",
		"jwt.secret":                    "JWT_SECRET",
		"jwt.access_token_ttl":          "JWT_ACCESS_TOKEN_TTL",
		"jwt.refresh_token_ttl":         "JWT_REFRESH_TOKEN_TTL",
		"jwt.ttlhours":                  "JWT_TTLHOURS",
		"server.port":                   "SERVER_PORT",
		"server.readtimeout":            "SERVER_READTIMEOUT",
		"server.writetimeout":           "SERVER_WRITETIMEOUT",
		"server.idletimeout":            "SERVER_IDLETIMEOUT",
		"server.shutdowntimeout":        "SERVER_SHUTDOWNTIMEOUT",
		"server.maxheaderbytes":         "SERVER_MAXHEADERBYTES",
		"logging.level":                 "LOGGING_LEVEL",
		"ratelimit.enabled":             "RATELIMIT_ENABLED",
		"ratelimit.requests":            "RATELIMIT_REQUESTS",
		"ratelimit.window":              "RATELIMIT_WINDOW",
		"migrations.directory":          "MIGRATIONS_DIRECTORY",
		"migrations.timeout":            "MIGRATIONS_TIMEOUT",
		"migrations.locktimeout":        "MIGRATIONS_LOCKTIMEOUT",
		"health.timeout":                "HEALTH_TIMEOUT",
		"health.database_check_enabled": "HEALTH_DATABASE_CHECK_ENABLED",

		// MongoDB
		"mongodb.uri":             "MONGODB_URI",
		"mongodb.database":        "MONGODB_DATABASE",

		// Redis
		"redis.host":           "REDIS_HOST",
		"redis.port":           "REDIS_PORT",
		"redis.password":       "REDIS_PASSWORD",
	
		// RabbitMQ
		"rabbitmq.url":            "RABBITMQ_URL",

		// gRPC
		"grpc.port":               "GRPC_PORT",

		// Metrics
		"metrics.port":    "METRICS_PORT",

	
	}
	for key, env := range envBindings {
		_ = v.BindEnv(key, env)
	}
}

func (l *LoggingConfig) GetLogLevel() slog.Level {
	switch strings.ToLower(l.Level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default to info level
	}
}

func GetSkipPaths(env string) []string {
	switch env {
	case "production":
		return []string{"/health", "/health/live", "/health/ready", "/metrics", "/debug", "/pprof"}
	case "development":
		return []string{"/health", "/health/live", "/health/ready"}
	case "test":
		return []string{"/health", "/health/live", "/health/ready"}
	default:
		return []string{"/health", "/health/live", "/health/ready"}
	}
}

func GetConfigPath() string {
	paths := []string{
		"configs/config.yaml",
		"./configs/config.yaml",
		"../configs/config.yaml",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	return "configs/config.yaml"
}

func (c *Config) LogSafeConfig(logger *slog.Logger) {
	logger.Info("Loaded Configuration:")
	logger.Info("App", "Name", c.App.Name, "Environment", c.App.Environment, "Debug", c.App.Debug)
	logger.Info("Database", "Host", c.Database.Host, "Port", c.Database.Port, "User", c.Database.User, "Password", "<redacted>", "Name", c.Database.Name, "SSLMode", c.Database.SSLMode)
	logger.Info("JWT", "Secret", "<redacted>", "AccessTokenTTL", c.JWT.AccessTokenTTL, "RefreshTokenTTL", c.JWT.RefreshTokenTTL)
	logger.Info("Server", "Port", c.Server.Port, "ReadTimeout", c.Server.ReadTimeout, "WriteTimeout", c.Server.WriteTimeout, "IdleTimeout", c.Server.IdleTimeout, "ShutdownTimeout", c.Server.ShutdownTimeout, "MaxHeaderBytes", c.Server.MaxHeaderBytes)
	logger.Info("Logging", "Level", c.Logging.Level)
	logger.Info("RateLimit", "Enabled", c.Ratelimit.Enabled, "Requests", c.Ratelimit.Requests, "Window", c.Ratelimit.Window)
	logger.Info("Migrations", "Directory", c.Migrations.Directory, "Timeout", c.Migrations.Timeout, "LockTimeout", c.Migrations.LockTimeout)
}
