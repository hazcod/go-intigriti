package company

import (
	intigriti "github.com/hazcod/go-intigriti/pkg/api"
	"github.com/sirupsen/logrus"
)

func DoAuth(l *logrus.Logger, inti intigriti.Endpoint) {
	l.Info("checking authentication status")

	if !inti.IsAuthenticated() {
		l.Fatal("client is not authenticated")
	}

	l.Info("client is authenticated successfully and cached in your configuration file")
}
