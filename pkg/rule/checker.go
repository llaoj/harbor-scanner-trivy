package rule

import (
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/etc"
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/harbor"
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
	return CheckBuildHistory(req, c.Config)
}
