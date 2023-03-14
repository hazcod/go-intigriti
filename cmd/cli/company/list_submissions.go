package company

import (
	"flag"
	"strings"

	v2 "github.com/finn-no/go-intigriti/v2"
	"github.com/sirupsen/logrus"
)

func ListSubmissions(l *logrus.Logger, inti v2.Endpoint) {
	programID := "*"
	if len(flag.Args()) >= 3 {
		programID = flag.Arg(2)
		l.WithField("program_id", programID).Debug("filtering submissions for program")
	}

	l.WithField("program_id", programID).Info("Listing company submissions")

	programIDs := make([]string, 0)

	if strings.TrimSpace(programID) == "" || programID == "*" {
		programs, err := inti.GetPrograms()
		if err != nil {
			l.WithError(err).Fatal("could not list programs")
		}
		for _, program := range programs {
			programIDs = append(programIDs, program.ID)
		}
	} else {
		programIDs = []string{programID}
	}

	// for _, programID := range programIDs {
	// submissions, err := inti.GetSubmissions(programID)
	submissions, err := inti.GetSubmissions()
	if err != nil {
		l.WithError(err).WithField("program_id", submissions).Warn("could not list submissions")
		return
		// continue
	}

	for _, subm := range submissions {
		l.Infof(
			"- %s (state %s, severity %s, code %s)",
			subm.Title, subm.State.Status.Value, subm.Severity.Value, subm.Code)
	}
}

// }
