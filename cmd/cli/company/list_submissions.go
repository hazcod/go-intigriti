package company

import (
	intigriti "github.com/hazcod/go-intigriti/pkg/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Filter struct {
	Program    string
	Title      string
	Code       string
	Status     string
	Researcher string
	Assignee   string
}

func CreateFilter(l *logrus.Logger, args []string) Filter {
	filter := Filter{}

	for _, arg := range args {
		l.WithField("arg", arg).Debug("adding to filter")

		if err := filter.Add(arg); err != nil {
			l.WithError(err).WithField("filter", arg).Warn("failed to apply filter")
		}
	}

	return filter
}

func (f *Filter) Add(filter string) error {
	parts := strings.SplitN(filter, "=", 2)
	if len(parts) != 2 {
		return errors.New("invalid filter: " + filter + " , use filter=value notation")
	}

	filterName := strings.ToLower(parts[0])
	filterVal := parts[1]

	switch filterName {
	case "program":
		f.Program = filterVal
		return nil
	case "title":
		f.Title = filterVal
		return nil
	case "code":
		f.Code = filterVal
		return nil
	case "status":
		f.Status = filterVal
		return nil
	case "researcher":
		f.Researcher = filterVal
		return nil
	case "assignee":
		f.Assignee = filterVal
		return nil
	default:
		return errors.New("unknown filter: " + filterName)
	}
}

func (f *Filter) Filter(submissions *[]intigriti.Submission) {
	for _, sub := range *submissions {
		if f.Code != "" && sub.Code != f.Code {
			continue
		}

		if f.Researcher != "" && sub.Submitter.UserName != f.Researcher {
			continue
		}

		if f.Assignee != "" && sub.Assignee.Username != f.Assignee {
			continue
		}

		if f.Title != "" && sub.Title != f.Title {
			continue
		}

		if f.Program != "" && sub.ProgramID != f.Program {
			continue
		}

		if f.Status != "" && sub.State.Status.Value != f.Status {
			continue
		}
	}
}

func ListSubmissions(l *logrus.Logger, inti intigriti.Endpoint) {
	filter := CreateFilter(l, os.Args[4:])

	l.Info("Listing company submissions")

	programIDs := make([]string, 0)

	if strings.TrimSpace(filter.Program) == "" || filter.Program == "*" {
		programs, err := inti.GetPrograms()
		if err != nil {
			l.WithError(err).Fatal("could not list programs")
		}

		for _, program := range programs {
			programIDs = append(programIDs, program.ID)
		}
	} else {
		l.WithField("program_id", filter.Program).Debug("filtering for program")
		programIDs = []string{filter.Program}
	}

	submissions := make([]intigriti.Submission, 0)

	l.WithField("programs", len(programIDs)).Debug("retrieving submissions")

	for _, programID := range programIDs {
		pSubmissions, err := inti.GetProgramSubmissions(programID)
		if err != nil {
			l.WithError(err).WithField("program_id", pSubmissions).Error("could not list submissions")
			continue
		}
		submissions = append(submissions, pSubmissions...)
	}

	l.WithField("submissions", len(submissions)).Debug("retrieved submissions, filtering...")

	filter.Filter(&submissions)

	l.WithField("submissions", len(submissions)).Debug("filtered submissions")

	for _, subm := range submissions {
		l.WithFields(logrus.Fields{
			"state":      subm.State.Status.Value,
			"severity":   subm.Severity.Value,
			"assignee":   subm.Assignee.Username,
			"researcher": subm.Submitter.UserName,
			"code":       subm.Code,
		}).Info(subm.Title)
	}
}
