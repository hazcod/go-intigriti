package company

import (
	"flag"
	intigriti "github.com/hazcod/go-intigriti/pkg/api"
	"github.com/sirupsen/logrus"
	"net"
)

func isFaultyIP(ip net.IP) bool {
	return ip == nil || ip.IsPrivate() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalMulticast() || ip.IsLoopback()
}

func CheckIP(l *logrus.Logger, inti intigriti.Endpoint) {
	if len(flag.Args()) != 3 {
		l.Fatal("usage: inti company ip <ip-address>")
	}

	ipAddress := flag.Arg(2)
	logger := l.WithField("ip_address", ipAddress)

	ip := net.ParseIP(ipAddress)

	if isFaultyIP(ip) {
		logger.Fatal("invalid ip address provided")
	}

	isResearcherIP, err := inti.IsKnownIP(ip)
	if err != nil {
		logger.WithError(err).Fatal("could not verify IP address")
	}

	if isResearcherIP {
		logger.Info("this is an ip address linked to an account on the Intigriti platform.")
	} else {
		logger.Info("this ip address is not known to the Intigriti platform.")
	}
}
