package config

import (
	"errors"
	"strings"

	"github.com/SpiderOak/errstack"
	"github.com/spf13/viper"
)

const (
	DefaultConfigFileName string = "config"

	MsgErrNoAppName string = "no AppName provided"
)

var (
	replacer = strings.NewReplacer("-", "_", " ", "_")

	ErrNoAppName error = errors.New(MsgErrNoAppName)
)

type Config struct {
	Logger    Logger
	Server    Server
	FreeClimb FreeClimb
}

type Logger struct {
	Production bool
	Level      string
}

type Server struct {
	Port     string
	Accept   string
	Secure   bool
	CertFile string
	KeyFile  string
}
type FreeClimb struct {
	AccountId string
	AuthToken string
	From      string
	To        string
}

func New() *Config {
	return &Config{}
}

func InitDefaults(vpr *viper.Viper) *viper.Viper {
	vpr.SetDefault("logger.production", true)
	vpr.SetDefault("logger.level", "debug")

	vpr.SetDefault("server.accept", "0.0.0.0")
	vpr.SetDefault("server.port", "8080")
	vpr.SetDefault("server.secure", false)
	vpr.SetDefault("server.certFile", "")
	vpr.SetDefault("server.keyFile", "")

	return vpr
}

func InitConfig(appName string) (*Config, error) {
	vpr, err := buildViper(appName)
	if err != nil {
		return nil, errstack.Push(err, "failed to setup configuration management")
	}
	vpr = InitDefaults(vpr)

	cfg := New()

	cfg.Logger.Production = vpr.GetBool("logger.production")
	cfg.Logger.Level = vpr.GetString("logger.level")

	cfg.Server.Accept = vpr.GetString("server.accept")
	cfg.Server.Port = vpr.GetString("server.port")
	cfg.Server.Secure = vpr.GetBool("server.secure")
	cfg.Server.CertFile = vpr.GetString("server.certFile")
	cfg.Server.KeyFile = vpr.GetString("server.keyFile")

	cfg.FreeClimb.AccountId = vpr.GetString("freeclimb.accountId")
	cfg.FreeClimb.AuthToken = vpr.GetString("freeclimb.authToken")
	cfg.FreeClimb.From = vpr.GetString("freeclimb.from")
	cfg.FreeClimb.To = vpr.GetString("freeclimb.to")

	err = validateConfig(cfg)
	if err != nil {
		return nil, errstack.Push(err, "config validation")
	}
	return cfg, nil
}

func buildViper(appName string) (*viper.Viper, error) {
	vpr := viper.New()

	if appName == "" {
		return nil, ErrNoAppName
	}

	vpr.SetEnvPrefix(replacer.Replace(strings.ToUpper(appName)))
	vpr.SetConfigName(DefaultConfigFileName)
	vpr.AddConfigPath(".")
	vpr.AddConfigPath("/etc/" + appName + "/")
	groupedPath := "$HOME/.config"
	vpr.AddConfigPath(groupedPath + "/." + appName + "/")
	vpr.AddConfigPath("$HOME/." + appName + "/")

	if err := vpr.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errstack.Pushf(err, "failed reading config file")
		}
	}

	vpr.AutomaticEnv()

	return vpr, nil
}

func validateConfig(cfg *Config) error {
	return nil
}
