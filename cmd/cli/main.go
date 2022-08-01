package main

import (
	"flag"
	"github.com/hazcod/go-intigriti/cmd/cli/company"
	"github.com/hazcod/go-intigriti/pkg/config"
	v2 "github.com/hazcod/go-intigriti/v2"
	"github.com/sirupsen/logrus"
	"strings"
)

func main() {
	logger := logrus.New()

	configPath := flag.String("config", "inti.yml", "Path to your config file.")
	logLevelStr := flag.String("log", "info", "Log level.")
	flag.Parse()

	logLevel, err := logrus.ParseLevel(*logLevelStr)
	if err != nil {
		logger.WithError(err).Fatal("could not parse log level")
	}

	logger.SetLevel(logLevel)
	logger.WithField("level", logLevel.String()).Debugf("log level set")

	cfg, err := config.Load(logger, *configPath)
	if err != nil {
		logger.Fatalf("could not load configuration: %s", err)
	}

	if err := cfg.Validate(); err != nil {
		logger.WithError(err).Fatal("invalid configuration")
	}

	inti, err := v2.New(cfg.Auth.ClientID, cfg.Auth.ClientSecret, &cfg.Cache, logger)
	if err != nil {
		logger.WithError(err).Fatal("could not initialize client")
	}

	intiToken, err := inti.GetToken()
	if err != nil {
		logger.WithError(err).Warn("could not extract token, skipping token cache")
	} else {
		if err := cfg.CacheAuth(logger, *configPath, intiToken); err != nil {
			logger.WithError(err).Warn("could not cache token")
		}
	}
	logger.WithField("valid", intiToken.Valid()).Debug("retrieved auth token")

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
