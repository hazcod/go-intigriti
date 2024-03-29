package main

import (
	"flag"
	"github.com/hazcod/go-intigriti/cmd/cli/company"
	"github.com/hazcod/go-intigriti/cmd/cli/ui"
	"github.com/hazcod/go-intigriti/cmd/config"
	intigriti "github.com/hazcod/go-intigriti/pkg/api"
	apiConfig "github.com/hazcod/go-intigriti/pkg/config"
	"github.com/sirupsen/logrus"
	"strings"
)

func main() {
	logger := logrus.New()

	configPath := flag.String("config", "inti.yml", "Path to your config file.")
	logLevelStr := flag.String("log", "", "Log level.")
	flag.Parse()

	if *logLevelStr != "" {
		logLevel, err := logrus.ParseLevel(*logLevelStr)
		if err != nil {
			logger.WithError(err).Fatal("could not parse log level")
		}

		logger.SetLevel(logLevel)
		logger.WithField("level", logLevel.String()).Debugf("log level set")
	}

	cfg, err := config.Load(logger, *configPath)
	if err != nil {
		logger.Fatalf("could not load configuration: %s", err)
	}

	if err := cfg.Validate(); err != nil {
		logger.WithError(err).Fatal("invalid configuration")
	}

	if cfg.Log.Level != "" && *logLevelStr == "" {
		logLevel, err := logrus.ParseLevel(cfg.Log.Level)
		if err != nil {
			logger.WithError(err).Fatal("could not parse log level")
		}

		logger.SetLevel(logLevel)
		logger.WithField("level", logLevel.String()).Debugf("log level set")
	}

	browser := ui.SystemBrowser{}

	inti, err := intigriti.New(apiConfig.Config{
		// our Intigriti API credentials
		Credentials: struct {
			ClientID     string
			ClientSecret string
		}{ClientID: cfg.Auth.ClientID, ClientSecret: cfg.Auth.ClientSecret},

		// pop up a browser when necessary to authenticate
		OpenBrowser:   true,
		Authenticator: browser,

		// cache tokens as much as possible to reduce times we have to authenticate
		TokenCache: &apiConfig.CachedToken{
			RefreshToken: cfg.Cache.RefreshToken,
			AccessToken:  cfg.Cache.AccessToken,
			ExpiryDate:   cfg.Cache.ExpiryDate,
			Type:         cfg.Cache.Type,
		},

		// use our logger and our logging levels
		Logger: logger,
	})
	if err != nil {
		logger.WithError(err).Fatal("could not initialize client")
	}

	logger.WithField("authenticated", inti.IsAuthenticated()).Debug("initialized client")

	if len(flag.Args()) == 0 {
		logger.Fatalf("no command provided. See: company")
	}

	command := strings.ToLower(flag.Args()[0])

	switch strings.ToLower(command) {
	case "company", "c", "com":
		company.Command(logger, cfg, inti)
		return
	default:
		logger.Fatalf("unknown command '%s'. See: company", command)
	}
}
