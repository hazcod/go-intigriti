package company

import (
	intigriti "github.com/hazcod/go-intigriti/pkg/api"
	"github.com/sirupsen/logrus"
)

func ListPrograms(l *logrus.Logger, inti intigriti.Endpoint) {
	l.Info("Listing company programs")

	programs, err := inti.GetPrograms()
	if err != nil {
		l.WithError(err).Fatal("could not list programs")
	}

	for _, program := range programs {
		l.Infof("- %s (type %s, status %s, handle %s)", program.Name, program.Type.Value, program.Status.Value, program.Handle)
	}
}
