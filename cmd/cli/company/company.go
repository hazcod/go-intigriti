package company

import (
	"flag"
	"strings"

	"github.com/finn-no/go-intigriti/pkg/config"
	v2 "github.com/finn-no/go-intigriti/v2"
	"github.com/sirupsen/logrus"
)

func Command(l *logrus.Logger, cfg *config.Config, inti v2.Endpoint) {
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

	default:
		l.Fatalf("Unknown subcommand '%s'. See: company <list,submissions>", subCommand)
	}
}
