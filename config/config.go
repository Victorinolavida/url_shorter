package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/flashlabs/rootpath/location"
	"github.com/knadh/koanf/parsers/yaml"
	env2 "github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func getConfigFile(env []string) (*koanf.Koanf, error) {
	k = koanf.New(".")
	confile := "config/config.yaml"
	if len(env) > 0 && env[0] != "" {
		confile = fmt.Sprintf("config/config.%s.yaml", env[0])
		if err := location.Chdir(); err != nil {
			return k, nil
		}
	}

	if err := k.Load(file.Provider(confile), yaml.Parser()); err != nil {
		return nil, err
	}

	return k, nil
}

func NewConfig[T any](env ...string) (*T, error) {
	k, err := getConfigFile(env)
	if err != nil {
		return nil, err
	}

	err = k.Load(env2.Provider("", ".", func(s string) string {
		val := strings.ToLower(strings.ReplaceAll(s, "_", "."))
		return val
	}), nil)
	if err != nil {
		return nil, err
	}
	conf := new(T)
	err = k.Unmarshal("", conf)
	if err != nil {
		return nil, err
	}

	if validator, ok := any(conf).(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return nil, err
		}
	}

	return conf, nil
}

type Server struct {
	Host         string `koanf:"host"`
	Port         int    `koanf:"port"`
	PublicDomain string `koanf:"public_domain"`
}

const (
	EnvLocal      = "local"
	EnvProduction = "production"

	defaultHost = "localhost"
	defaultPort = 8080
)

type Config struct {
	Env            string        `koanf:"env"`
	Server         Server        `koanf:"server"`
	Database       Database      `koanf:"database"`
	LogLevel       string        `koanf:"log_level"`
	WorkerInterval time.Duration `koanf:"worker_duration"`
}

func (c *Config) Validate() error {
	switch c.Env {
	case "":
		c.Env = EnvLocal
	case EnvLocal, EnvProduction:
	default:
		return fmt.Errorf("invalid env %q: allowed values are %q, %q or empty", c.Env, EnvLocal, EnvProduction)
	}

	if c.Env == EnvLocal {
		c.Server.PublicDomain = localPublicDomain(c.Server.Host, c.Server.Port)
		return nil
	}

	return validateProductionDomain(c.Server.PublicDomain)
}

func localPublicDomain(host string, port int) string {
	if host == "" {
		host = defaultHost
	}
	if port == 0 {
		port = defaultPort
	}
	return "http://" + net.JoinHostPort(host, strconv.Itoa(port))
}

func validateProductionDomain(domain string) error {
	if domain == "" {
		return nil
	}

	lower := strings.ToLower(domain)
	if strings.Contains(lower, "://") || strings.HasPrefix(lower, "http:") || strings.HasPrefix(lower, "https:") {
		return fmt.Errorf("invalid server.public_domain %q: url scheme is not allowed in production", domain)
	}

	u, err := url.Parse("//" + domain)
	if err != nil {
		return fmt.Errorf("invalid server.public_domain %q: %w", domain, err)
	}

	host := strings.ToLower(u.Hostname())
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return fmt.Errorf("invalid server.public_domain %q: localhost is not allowed in production", domain)
	}

	return nil
}

type Database struct {
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Name     string `koanf:"name"`
	MaxConns int    `koanf:"max_conns"`
	MinConns int    `koanf:"min_conns"`
}
