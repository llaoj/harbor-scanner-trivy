package rule

import (
	"strings"

	"github.com/aquasecurity/harbor-scanner-trivy/pkg/etc"
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/harbor"
	log "github.com/sirupsen/logrus"
)

func CheckBuildHistory(req harbor.ScanRequest, config etc.RuleChecker) error {
	harbor := NewHarborClient(req.Registry.URL, req.Registry.Authorization)
	project, repo, ok := parseRepository(req.Artifact.Repository)
	log.WithFields(log.Fields{
		"project": project,
		"repo":    repo,
	}).Trace("Parse Repository")
	if !ok {
		return nil
	}
	history, err := harbor.GetBuildHistory(project, repo, req.Artifact.Digest)
	log.WithFields(log.Fields{
		"history": history,
	}).Trace("Get Build History")
	if err != nil {
		return err
	}

	// is validated
	labelExist := false
	baseImageExist := false
	for _, record := range history {
		for _, label := range strings.Split(config.ImageLabels, ",") {
			if strings.Contains(record.CreatedBy, label) {
				labelExist = true
				break
			}
		}
	}
	// first record for base image
	for _, digest := range strings.Split(config.BaseImageDigests, ",") {
		if strings.Contains(history[0].CreatedBy, digest) {
			baseImageExist = true
			break
		}
	}

	if labelExist && baseImageExist {
		return harbor.AddLegalLabel(project, repo, req.Artifact.Digest)
	}
	return harbor.AddIlegalLabel(project, repo, req.Artifact.Digest)
}
