package config

import (
	"github.com/juju/fslock"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

const (
	appEnvPrefix = "INTI"
)

type Config struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Auth struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"auth"`

	Cache TokenCache `yaml:"cache"`
}

type TokenCache struct {
	RefreshToken string    `yaml:"refresh_token"`
	AccessToken  string    `yanl:"access_token"`
	ExpiryDate   time.Time `yaml:"expiry"`
	Type         string    `yaml:"type"`
}

func Load(logger *logrus.Logger, path string) (*Config, error) {
	var config Config

	if path != "" {
		configBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrap(err, "could not load configuration file")
		}

		if err := yaml.Unmarshal(configBytes, &config); err != nil {
			return nil, errors.Wrap(err, "could not parse configuration file")
		}

		logger.WithField("config", path).Debug("loaded configuration")
	}

	if err := envconfig.Process(appEnvPrefix, &config); err != nil {
		return nil, errors.Wrap(err, "could not load environment variables")
	}

	return &config, nil
}

func (c *Config) Save(logger *logrus.Logger, path string) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "could not serialize config")
	}

	if err := ioutil.WriteFile(path, b, 0600); err != nil {
		return errors.Wrap(err, "could not write config")
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Auth.ClientID == "" {
		return errors.New("no clientid provided")
	}

	if c.Auth.ClientSecret == "" {
		return errors.New("no client secret provided")
	}

	return nil
}

func (c *Config) CacheAuth(l *logrus.Logger, path string, token *oauth2.Token) error {
	if !token.Valid() {
		return errors.New("token is not valid")
	}

	lock := fslock.New(path)
	if err := lock.TryLock(); err != nil {
		return errors.Wrap(err, "could not lock config file")
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			l.WithError(err).Warn("could not release config lock")
		}
	}()

	c.Cache.AccessToken = token.AccessToken
	c.Cache.RefreshToken = token.RefreshToken
	c.Cache.ExpiryDate = token.Expiry
	c.Cache.Type = token.TokenType

	if err := c.Save(l, path); err != nil {
		return errors.Wrap(err, "failed to save config")
	}

	return nil
}
