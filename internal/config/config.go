package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	envLocal = "LOCAL"
	envProd  = "PROD"

	configFileLocal = "config/config_local.yaml"
	configFileProd  = "config/config_prod.yaml"
)

type Config struct {
	Token         string `yaml:"token"`
	FixerAPIToken string `yaml:"fixer-api-token"`

	DBUri  string `yaml:"dbUri"`
	DBName string `yaml:"dbName"`
	DBUser string `yaml:"dbUser"`
	DBPass string `yaml:"dbPass"`

	RedisUri  string `yaml:"redisUri"`
	RedisPass string `yaml:"redisPass"`
}

type Service struct {
	config Config
	env    string
}

func New(env string) (*Service, error) {
	s := &Service{env: env}

	var configFile string
	if env == envLocal {
		configFile = configFileLocal
	} else if env == envProd {
		configFile = configFileProd
	}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

func (s *Service) Token() string {
	return s.config.Token
}

func (s *Service) FixerAPIToken() string {
	return s.config.FixerAPIToken
}

func (s *Service) DBUri() string {
	return s.config.DBUri
}

func (s *Service) RedisUri() string {
	return s.config.RedisUri
}

func (s *Service) RedisPass() string {
	return s.config.RedisPass
}
