package service

import v1 "github.com/klovercloud-ci/ctl/v1"

type Git interface {
	Apply(git v1.Git, companyId, appId, apiServerUrl, token string, skipSsl bool) error
}
