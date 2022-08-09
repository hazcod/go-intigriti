package company

import (
	"flag"
	"strings"

	intigriti "github.com/intigriti/sdk-go/pkg/api"
	"github.com/intigriti/sdk-go/pkg/config"
	"github.com/sirupsen/logrus"
)

func Command(l *logrus.Logger, cfg *config.Config, inti intigriti.Endpoint) {
	if len(flag.Args()) < 2 {
		l.Fatal("Missing subcommand. See: company <list,submissions>")
	}

	subCommand := strings.ToLower(flag.Arg(1))

	switch subCommand {
	case "ls", "list", "list-programs":
		ListPrograms(l, inti)
		return

	case "sub", "submissions", "list-submissions":
		ListSubmissions(l, inti)
		return

	case "check-ip", "ip":
		CheckIP(l, inti)
		return

	default:
		l.Fatalf("Unknown subcommand '%s'. See: company <list,submissions>", subCommand)
	}
}
