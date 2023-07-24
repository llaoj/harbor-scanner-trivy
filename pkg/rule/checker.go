package rule

import (
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/etc"
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/harbor"
	log "github.com/sirupsen/logrus"
)

type Checker struct {
	Config etc.RuleChecker
}

func NewChecker(config etc.RuleChecker) *Checker {
	return &Checker{
		Config: config,
	}
}

func (c *Checker) Check(req harbor.ScanRequest) error {
	log.WithFields(log.Fields{
		"baseUrl":    req.Registry.URL,
		"auth":       req.Registry.Authorization,
		"repository": req.Artifact.Repository,
		"digest":     req.Artifact.Digest,
	}).Debug("Scan Request")
	return CheckBuildHistory(req, c.Config)
}
